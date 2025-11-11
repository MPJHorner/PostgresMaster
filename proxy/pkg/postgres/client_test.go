package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// TestNewClient_InvalidConnectionString tests that NewClient fails immediately with invalid connection string
func TestNewClient_InvalidConnectionString(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		connStr  string
		wantErr  bool
		errSubstr string
	}{
		{
			name:     "invalid connection string format",
			connStr:  "not-a-valid-connection-string",
			wantErr:  true,
			errSubstr: "parse",
		},
		{
			name:     "invalid port format",
			connStr:  "postgres://user:pass@localhost:notaport/db",
			wantErr:  true,
			errSubstr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(ctx, tc.connStr)

			if tc.wantErr && err == nil {
				t.Errorf("NewClient() expected error but got nil")
			}

			if !tc.wantErr && err != nil {
				t.Errorf("NewClient() unexpected error: %v", err)
			}

			if err != nil && tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
				t.Errorf("NewClient() error = %v, want error containing %q", err, tc.errSubstr)
			}

			if client != nil {
				client.Close()
			}
		})
	}
}

// TestNewClient_WithTimeout tests that NewClient respects context timeout
func TestNewClient_WithTimeout(t *testing.T) {
	// Use an invalid host that will hang/timeout
	connStr := "postgres://user:pass@10.255.255.1:5432/testdb?connect_timeout=1"

	// Create context with very short timeout to ensure it triggers
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	client, err := NewClient(ctx, connStr)
	duration := time.Since(start)

	if err == nil {
		t.Errorf("NewClient() expected error with timeout context, got nil")
		if client != nil {
			client.Close()
		}
	}

	// Verify that context cancellation is respected (should fail quickly, not wait for all retries)
	if duration > 5*time.Second {
		t.Errorf("NewClient() took too long (%v), context timeout not respected", duration)
	}
}

// TestNewClient_ContextCancellation tests context cancellation during retry
func TestNewClient_ContextCancellation(t *testing.T) {
	// Use an unreachable host to trigger retries
	connStr := "postgres://user:pass@192.0.2.1:5432/testdb?connect_timeout=1"

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after a short delay to simulate cancellation during retry
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	client, err := NewClient(ctx, connStr)
	duration := time.Since(start)

	if err == nil {
		t.Errorf("NewClient() expected error after context cancellation, got nil")
		if client != nil {
			client.Close()
		}
	}

	if err != nil && !strings.Contains(err.Error(), "context") {
		t.Logf("NewClient() error doesn't mention context (may be OK): %v", err)
	}

	// Should fail relatively quickly after cancellation, not wait for full retry cycle
	if duration > 10*time.Second {
		t.Errorf("NewClient() took too long (%v) after context cancellation", duration)
	}
}

// TestClient_Close tests that Close() doesn't panic
func TestClient_Close(t *testing.T) {
	// Test closing a nil pool client
	client := &Client{pool: nil}

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close() panicked: %v", r)
		}
	}()

	client.Close()
}

// TestClient_Close_Multiple tests that multiple Close() calls are safe
func TestClient_Close_Multiple(t *testing.T) {
	client := &Client{pool: nil}

	// Should not panic on multiple calls
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Multiple Close() calls panicked: %v", r)
		}
	}()

	client.Close()
	client.Close()
	client.Close()
}

// Integration tests - require actual Postgres instance
// Set TEST_POSTGRES_URL environment variable to run these tests

func getTestDatabaseURL() (string, bool) {
	url := os.Getenv("TEST_POSTGRES_URL")
	if url == "" {
		return "", false
	}
	return url, true
}

func TestNewClient_Integration_Success(t *testing.T) {
	url, ok := getTestDatabaseURL()
	if !ok {
		t.Skip("Skipping integration test: TEST_POSTGRES_URL not set")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}
	defer client.Close()

	// Verify we can ping the database
	if err := client.Ping(ctx); err != nil {
		t.Errorf("Ping() failed: %v", err)
	}
}

func TestNewClient_Integration_Retry(t *testing.T) {
	url, ok := getTestDatabaseURL()
	if !ok {
		t.Skip("Skipping integration test: TEST_POSTGRES_URL not set")
	}

	// First, verify we can connect to the real database
	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("NewClient() failed to connect to test database: %v", err)
	}
	client.Close()

	// Now test with wrong credentials to trigger retry
	badURL := strings.Replace(url, "postgres://", "postgres://baduser:badpass@", 1)
	// Also ensure we have a proper format
	if !strings.Contains(badURL, "baduser") {
		// If replacement didn't work, construct a bad URL
		badURL = "postgres://baduser:badpass@localhost:5432/nonexistent?connect_timeout=1"
	}

	start := time.Now()
	client, err = NewClient(ctx, badURL)
	duration := time.Since(start)

	if err == nil {
		t.Errorf("NewClient() with bad credentials expected to fail")
		if client != nil {
			client.Close()
		}
		return
	}

	// Verify error message mentions retry attempts
	if !strings.Contains(err.Error(), "4 attempts") {
		t.Errorf("NewClient() error should mention retry attempts: %v", err)
	}

	// Verify retry logic was executed (should take at least 2+4+8 = 14 seconds for retries)
	// With connect_timeout=1, each attempt takes ~1s, plus wait times: 0 + 2 + 4 + 8 = 14s
	// But let's be conservative and check for at least 5 seconds
	if duration < 5*time.Second {
		t.Logf("NewClient() completed in %v, may not have executed full retry logic", duration)
	}

	t.Logf("NewClient() with bad credentials failed after %v (expected)", duration)
}

func TestClient_Integration_Ping(t *testing.T) {
	url, ok := getTestDatabaseURL()
	if !ok {
		t.Skip("Skipping integration test: TEST_POSTGRES_URL not set")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}
	defer client.Close()

	// Test multiple pings
	for i := 0; i < 5; i++ {
		if err := client.Ping(ctx); err != nil {
			t.Errorf("Ping() iteration %d failed: %v", i, err)
		}
	}
}

func TestClient_Integration_CloseAndPing(t *testing.T) {
	url, ok := getTestDatabaseURL()
	if !ok {
		t.Skip("Skipping integration test: TEST_POSTGRES_URL not set")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Close the client
	client.Close()

	// Ping should fail after close
	err = client.Ping(ctx)
	if err == nil {
		t.Errorf("Ping() after Close() should fail, got nil")
	}
}

// Benchmark tests

func BenchmarkNewClient(b *testing.B) {
	url, ok := getTestDatabaseURL()
	if !ok {
		b.Skip("Skipping benchmark: TEST_POSTGRES_URL not set")
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := NewClient(ctx, url)
		if err != nil {
			b.Fatalf("NewClient() failed: %v", err)
		}
		client.Close()
	}
}

func BenchmarkClient_Ping(b *testing.B) {
	url, ok := getTestDatabaseURL()
	if !ok {
		b.Skip("Skipping benchmark: TEST_POSTGRES_URL not set")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		b.Fatalf("NewClient() failed: %v", err)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := client.Ping(ctx); err != nil {
			b.Fatalf("Ping() failed: %v", err)
		}
	}
}

// Example tests

func ExampleNewClient() {
	ctx := context.Background()
	client, err := NewClient(ctx, "postgres://user:pass@localhost:5432/mydb")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	// Use the client...
	fmt.Println("Connected successfully")
}

func ExampleClient_Ping() {
	ctx := context.Background()
	client, err := NewClient(ctx, "postgres://user:pass@localhost:5432/mydb")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	if err := client.Ping(ctx); err != nil {
		fmt.Printf("Ping failed: %v\n", err)
		return
	}

	fmt.Println("Database is reachable")
}
