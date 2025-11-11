package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/auth"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/postgres"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/server"
)

const (
	defaultPort = "8080"
	version     = "0.1.0"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// TODO: Parse CLI arguments or start interactive mode
	// For now, just require connection string as first argument
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <postgres-connection-string>", os.Args[0])
	}

	connString := os.Args[1]

	// Generate secret
	secret, err := auth.GenerateSecret()
	if err != nil {
		return fmt.Errorf("failed to generate secret: %w", err)
	}

	// Connect to Postgres with retry logic
	log.Println("Connecting to PostgreSQL...")
	ctx := context.Background()

	var pgClient *postgres.Client
	maxRetries := 4
	for attempt := 1; attempt <= maxRetries; attempt++ {
		pgClient, err = postgres.NewClient(ctx, connString)
		if err == nil {
			break
		}

		if attempt < maxRetries {
			waitTime := time.Duration(1<<uint(attempt-1)) * 2 * time.Second
			log.Printf("Connection failed (attempt %d/%d). Retrying in %v...", attempt, maxRetries, waitTime)
			time.Sleep(waitTime)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxRetries, err)
	}
	defer pgClient.Close()

	log.Println("✓ Connected to PostgreSQL!")

	// Start WebSocket server
	wsServer := server.NewServer(secret)
	http.HandleFunc("/", wsServer.HandleConnection)

	// Print connection URL
	fmt.Printf("\n→ Open in browser: http://localhost:%s?secret=%s\n\n", defaultPort, secret)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:         ":" + defaultPort,
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("WebSocket server listening on port %s", defaultPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("\nShutting down gracefully...")

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
