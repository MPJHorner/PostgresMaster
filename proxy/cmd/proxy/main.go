package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Define flags
	showHelp := flag.Bool("help", false, "Show help message")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showHelp, "h", false, "Show help message (shorthand)")
	flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")

	// Custom usage message
	flag.Usage = printUsage

	// Parse flags
	flag.Parse()

	// Handle version flag
	if *showVersion {
		printVersion()
		return nil
	}

	// Handle help flag
	if *showHelp {
		printUsage()
		return nil
	}

	var connString string
	var err error

	// Check if connection string provided as argument
	args := flag.Args()
	if len(args) >= 1 {
		// Option A: Connection string provided
		connString = args[0]

		// Validate connection string format
		if err := validateConnectionString(connString); err != nil {
			return err
		}

		fmt.Printf("✓ Using connection string from arguments\n")
	} else {
		// Option B: Interactive mode
		connString, err = promptForConnection()
		if err != nil {
			return err
		}
	}

	// Generate secret
	fmt.Println()
	fmt.Printf("🔐 Generating session secret...\n")
	secret, err := auth.GenerateSecret()
	if err != nil {
		return fmt.Errorf("failed to generate secret: %w", err)
	}
	fmt.Printf("✓ Session secret generated\n\n")

	// Connect to Postgres (NewClient handles retry logic internally)
	fmt.Printf("🔌 Connecting to PostgreSQL...\n")
	ctx := context.Background()
	pgClient, err := postgres.NewClient(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w\n\n"+
			"Troubleshooting tips:\n"+
			"  • Verify the database server is running\n"+
			"  • Check that the host and port are correct\n"+
			"  • Ensure your credentials are valid\n"+
			"  • Confirm the database exists\n"+
			"  • Check firewall settings if connecting remotely", err)
	}
	defer pgClient.Close()
	fmt.Printf("✓ Connected to PostgreSQL successfully\n\n")

	// Start WebSocket server
	wsServer := server.NewServer(secret, pgClient)
	http.HandleFunc("/", wsServer.HandleConnection)

	// Print connection URL with box
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  🚀 Proxy Server Running\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Printf("  📍 Local Address:  http://localhost:%s\n", defaultPort)
	fmt.Printf("  🔑 Session Secret: %s\n", secret)
	fmt.Println()
	fmt.Printf("  → Open in browser: http://localhost:%s?secret=%s\n", defaultPort, secret)
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  💡 Press Ctrl+C to stop the server")
	fmt.Println()

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
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "\n❌ Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-stop
	fmt.Println("\n🛑 Shutting down gracefully...")

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("✓ Server stopped successfully")
	return nil
}

// validateConnectionString validates a Postgres connection string format
func validateConnectionString(connStr string) error {
	if connStr == "" {
		return fmt.Errorf("connection string is empty\n\n" +
			"Expected format: postgres://[user:pass@]host[:port]/database[?params]\n" +
			"Example: postgres://localhost/mydb\n" +
			"Run 'postgres-proxy --help' for more information")
	}

	// Parse the connection string to validate format
	u, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("invalid connection string format: %w\n\n"+
			"Expected format: postgres://[user:pass@]host[:port]/database[?params]\n"+
			"Example: postgres://user:pass@localhost:5432/mydb", err)
	}

	// Check for postgres or postgresql scheme
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("invalid scheme '%s'\n\n"+
			"Connection string must start with 'postgres://' or 'postgresql://'\n"+
			"Example: postgres://localhost/mydb", u.Scheme)
	}

	// Check for required components
	if u.Host == "" {
		return fmt.Errorf("host is missing from connection string\n\n" +
			"The connection string must include a hostname or IP address.\n" +
			"Example: postgres://localhost/mydb\n" +
			"Example: postgres://192.168.1.100/mydb")
	}

	if u.Path == "" || u.Path == "/" {
		return fmt.Errorf("database name is missing from connection string\n\n" +
			"The connection string must include a database name after the host.\n" +
			"Example: postgres://localhost/mydb\n" +
			"Example: postgres://localhost:5432/production")
	}

	return nil
}

// promptForConnection prompts the user interactively for connection details
func promptForConnection() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  📋 PostgreSQL Connection Setup")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Host
	host, err := readInput(reader, "Host", "")
	if err != nil {
		return "", err
	}
	if host == "" {
		return "", fmt.Errorf("host is required\n\nExample: localhost or db.example.com")
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
		return "", fmt.Errorf("database is required\n\nExample: mydb or postgres")
	}

	// Username
	username, err := readInput(reader, "Username", "")
	if err != nil {
		return "", err
	}
	if username == "" {
		return "", fmt.Errorf("username is required\n\nExample: postgres or your database user")
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
	fmt.Printf("✓ Configuration complete\n")

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

// printVersion prints version information
func printVersion() {
	fmt.Printf("PostgreSQL Proxy v%s\n", version)
	fmt.Println("A lightweight WebSocket-to-PostgreSQL bridge for browser-based SQL clients")
	fmt.Println("\nProject: https://github.com/MPJHorner/PostgresMaster")
	fmt.Println("License: AGPL-3.0")
}

// printUsage prints usage information
func printUsage() {
	fmt.Println("PostgreSQL Proxy - WebSocket to PostgreSQL Bridge")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  postgres-proxy [OPTIONS] [CONNECTION_STRING]")
	fmt.Println("  postgres-proxy [OPTIONS]")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Starts a local WebSocket server that bridges browser connections to PostgreSQL.")
	fmt.Println("  The proxy enables secure database access from web applications without exposing")
	fmt.Println("  credentials to third-party services.")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println()
	fmt.Println("USAGE MODES:")
	fmt.Println()
	fmt.Println("  1. Connection String Mode:")
	fmt.Println("     Provide a PostgreSQL connection string as an argument")
	fmt.Println()
	fmt.Println("     Example:")
	fmt.Println("       postgres-proxy \"postgres://user:pass@host:5432/dbname?sslmode=require\"")
	fmt.Println()
	fmt.Println("  2. Interactive Mode:")
	fmt.Println("     Run without arguments to be prompted for connection details")
	fmt.Println()
	fmt.Println("     Example:")
	fmt.Println("       postgres-proxy")
	fmt.Println()
	fmt.Println("CONNECTION STRING FORMAT:")
	fmt.Println("  postgres://[user[:password]@][host][:port]/database[?param=value]")
	fmt.Println()
	fmt.Println("  Required components:")
	fmt.Println("    - host: Database server hostname or IP")
	fmt.Println("    - database: Database name")
	fmt.Println()
	fmt.Println("  Optional components:")
	fmt.Println("    - user: Username (default: postgres)")
	fmt.Println("    - password: User password")
	fmt.Println("    - port: Server port (default: 5432)")
	fmt.Println("    - sslmode: SSL mode (disable, allow, prefer, require, verify-ca, verify-full)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println()
	fmt.Println("  # Local database with default settings")
	fmt.Println("  postgres-proxy \"postgres://localhost/mydb\"")
	fmt.Println()
	fmt.Println("  # Remote database with SSL")
	fmt.Println("  postgres-proxy \"postgres://admin:secret@db.example.com:5432/prod?sslmode=require\"")
	fmt.Println()
	fmt.Println("  # Interactive mode")
	fmt.Println("  postgres-proxy")
	fmt.Println()
	fmt.Println("AFTER STARTING:")
	fmt.Println("  The proxy will display a URL like:")
	fmt.Println("    → Open in browser: http://localhost:8080?secret=abc123...")
	fmt.Println()
	fmt.Println("  Copy this URL and open it in your browser to connect to the web interface.")
	fmt.Println()
	fmt.Println("SECURITY:")
	fmt.Println("  - The proxy runs locally on your machine (localhost only)")
	fmt.Println("  - A unique secret is generated for each session")
	fmt.Println("  - Database credentials never leave your machine")
	fmt.Println("  - All connections are authenticated with the session secret")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/MPJHorner/PostgresMaster")
}
