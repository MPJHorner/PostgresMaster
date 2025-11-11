package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/auth"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/postgres"
	"github.com/MPJHorner/PostgresMaster/proxy/pkg/server"
	"golang.org/x/term"
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
	var connString string
	var err error

	// Check if connection string provided as argument
	if len(os.Args) >= 2 {
		// Option A: Connection string provided
		connString = os.Args[1]

		// Validate connection string format
		if err := validateConnectionString(connString); err != nil {
			return fmt.Errorf("invalid connection string: %w", err)
		}

		log.Printf("Using connection string from arguments")
	} else {
		// Option B: Interactive mode
		connString, err = promptForConnection()
		if err != nil {
			return fmt.Errorf("failed to get connection details: %w", err)
		}
	}

	// Generate secret
	secret, err := auth.GenerateSecret()
	if err != nil {
		return fmt.Errorf("failed to generate secret: %w", err)
	}

	// Connect to Postgres (NewClient handles retry logic internally)
	ctx := context.Background()
	pgClient, err := postgres.NewClient(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer pgClient.Close()

	// Start WebSocket server
	wsServer := server.NewServer(secret, pgClient)
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

// validateConnectionString validates a Postgres connection string format
func validateConnectionString(connStr string) error {
	if connStr == "" {
		return fmt.Errorf("connection string is empty")
	}

	// Parse the connection string to validate format
	u, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check for postgres or postgresql scheme
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("connection string must start with postgres:// or postgresql://")
	}

	// Check for required components
	if u.Host == "" {
		return fmt.Errorf("host is required")
	}

	if u.Path == "" || u.Path == "/" {
		return fmt.Errorf("database name is required")
	}

	return nil
}

// promptForConnection prompts the user interactively for connection details
func promptForConnection() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== PostgreSQL Connection Setup ===")
	fmt.Println()

	// Host
	host, err := readInput(reader, "Host", "")
	if err != nil {
		return "", err
	}
	if host == "" {
		return "", fmt.Errorf("host is required")
	}

	// Port
	port, err := readInput(reader, "Port", "5432")
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "5432"
	}

	// Database
	database, err := readInput(reader, "Database", "")
	if err != nil {
		return "", err
	}
	if database == "" {
		return "", fmt.Errorf("database is required")
	}

	// Username
	username, err := readInput(reader, "Username", "")
	if err != nil {
		return "", err
	}
	if username == "" {
		return "", fmt.Errorf("username is required")
	}

	// Password (hidden input)
	password, err := readPassword("Password")
	if err != nil {
		return "", err
	}

	// SSL Mode
	sslMode, err := readInput(reader, "SSL Mode", "prefer")
	if err != nil {
		return "", err
	}
	if sslMode == "" {
		sslMode = "prefer"
	}

	// Build connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(username),
		url.QueryEscape(password),
		host,
		port,
		database,
		sslMode,
	)

	fmt.Println()
	fmt.Printf("Connecting to postgres://%s:%s/%s...\n", host, port, database)

	return connString, nil
}

// readInput reads a line of input from the user with a prompt and optional default value
func readInput(reader *bufio.Reader, prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	// Return default if input is empty and default is provided
	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return input, nil
}

// readPassword reads a password from the user with hidden input
func readPassword(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)

	// Read password with hidden input
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Println() // Print newline after password input

	password := string(passwordBytes)
	if password == "" {
		return "", fmt.Errorf("password is required")
	}

	return password, nil
}
