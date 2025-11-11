package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/protocol"
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

// Tests for ExecuteQuery

func TestClient_Integration_ExecuteQuery_SimpleSelect(t *testing.T) {
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

	// Execute a simple SELECT query
	result, err := client.ExecuteQuery(ctx, "SELECT 1 as num, 'test' as text", nil)
	if err != nil {
		t.Fatalf("ExecuteQuery() failed: %v", err)
	}

	// Verify results
	if result.RowCount != 1 {
		t.Errorf("Expected 1 row, got %d", result.RowCount)
	}

	if len(result.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(result.Columns))
	}

	if len(result.Rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(result.Rows))
	}

	row := result.Rows[0]
	if row["num"] == nil {
		t.Error("Expected 'num' column to have value")
	}
	if row["text"] != "test" {
		t.Errorf("Expected 'text' to be 'test', got %v", row["text"])
	}

	// Verify execution time is measured
	if result.ExecutionTime == 0 {
		t.Error("ExecutionTime should be non-zero")
	}

	t.Logf("Query executed in %v", result.ExecutionTime)
}

func TestClient_Integration_ExecuteQuery_WithParameters(t *testing.T) {
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

	// Execute query with parameters
	params := []interface{}{42, "hello"}
	result, err := client.ExecuteQuery(ctx, "SELECT $1::int as num, $2::text as text", params)
	if err != nil {
		t.Fatalf("ExecuteQuery() with parameters failed: %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("Expected 1 row, got %d", result.RowCount)
	}

	row := result.Rows[0]
	// Note: pgx returns numeric types as int32, int64, etc.
	if numVal, ok := row["num"].(int32); !ok || numVal != 42 {
		t.Errorf("Expected 'num' to be 42, got %v (%T)", row["num"], row["num"])
	}
	if row["text"] != "hello" {
		t.Errorf("Expected 'text' to be 'hello', got %v", row["text"])
	}
}

func TestClient_Integration_ExecuteQuery_NullValues(t *testing.T) {
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

	// Execute query with NULL values
	result, err := client.ExecuteQuery(ctx, "SELECT NULL as null_col, 1 as int_col", nil)
	if err != nil {
		t.Fatalf("ExecuteQuery() failed: %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("Expected 1 row, got %d", result.RowCount)
	}

	row := result.Rows[0]
	if row["null_col"] != nil {
		t.Errorf("Expected 'null_col' to be nil, got %v", row["null_col"])
	}
	if row["int_col"] == nil {
		t.Error("Expected 'int_col' to have value")
	}
}

func TestClient_Integration_ExecuteQuery_MultipleRows(t *testing.T) {
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

	// Execute query that returns multiple rows
	result, err := client.ExecuteQuery(ctx, "SELECT generate_series(1, 10) as n", nil)
	if err != nil {
		t.Fatalf("ExecuteQuery() failed: %v", err)
	}

	if result.RowCount != 10 {
		t.Errorf("Expected 10 rows, got %d", result.RowCount)
	}

	if len(result.Rows) != 10 {
		t.Errorf("Expected 10 rows in result, got %d", len(result.Rows))
	}

	// Verify first and last rows
	if firstVal, ok := result.Rows[0]["n"].(int32); !ok || firstVal != 1 {
		t.Errorf("Expected first row to have n=1, got %v", result.Rows[0]["n"])
	}
	if lastVal, ok := result.Rows[9]["n"].(int32); !ok || lastVal != 10 {
		t.Errorf("Expected last row to have n=10, got %v", result.Rows[9]["n"])
	}
}

func TestClient_Integration_ExecuteQuery_DataTypes(t *testing.T) {
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

	// Execute query with various data types
	query := `
		SELECT
			true as bool_col,
			42 as int_col,
			3.14::float as float_col,
			'text value' as text_col,
			NOW() as timestamp_col,
			'2024-01-01'::date as date_col
	`
	result, err := client.ExecuteQuery(ctx, query, nil)
	if err != nil {
		t.Fatalf("ExecuteQuery() failed: %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("Expected 1 row, got %d", result.RowCount)
	}

	row := result.Rows[0]

	// Check boolean
	if boolVal, ok := row["bool_col"].(bool); !ok || !boolVal {
		t.Errorf("Expected 'bool_col' to be true, got %v (%T)", row["bool_col"], row["bool_col"])
	}

	// Check integer
	if row["int_col"] == nil {
		t.Error("Expected 'int_col' to have value")
	}

	// Check float
	if row["float_col"] == nil {
		t.Error("Expected 'float_col' to have value")
	}

	// Check text
	if row["text_col"] != "text value" {
		t.Errorf("Expected 'text_col' to be 'text value', got %v", row["text_col"])
	}

	// Check timestamp (should be converted to string in RFC3339 format)
	if timestampStr, ok := row["timestamp_col"].(string); !ok {
		t.Errorf("Expected 'timestamp_col' to be string, got %T", row["timestamp_col"])
	} else {
		// Verify it's a valid RFC3339 timestamp
		_, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			t.Errorf("Expected 'timestamp_col' to be RFC3339 format, got %v: %v", timestampStr, err)
		}
	}

	// Check date (should also be converted to string)
	if dateStr, ok := row["date_col"].(string); !ok {
		t.Errorf("Expected 'date_col' to be string, got %T", row["date_col"])
	} else {
		// Verify it's a valid timestamp
		_, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			t.Errorf("Expected 'date_col' to be RFC3339 format, got %v: %v", dateStr, err)
		}
	}

	// Verify column metadata
	if len(result.Columns) != 6 {
		t.Errorf("Expected 6 columns, got %d", len(result.Columns))
	}

	// Check that column data types are set
	for _, col := range result.Columns {
		if col.Name == "" {
			t.Error("Column name should not be empty")
		}
		if col.DataType == "" {
			t.Errorf("Column %s should have data type", col.Name)
		}
		t.Logf("Column: %s, Type: %s, OID: %d", col.Name, col.DataType, col.TypeOID)
	}
}

func TestClient_Integration_ExecuteQuery_Timeout(t *testing.T) {
	url, ok := getTestDatabaseURL()
	if !ok {
		t.Skip("Skipping integration test: TEST_POSTGRES_URL not set")
	}

	client, err := NewClient(context.Background(), url)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}
	defer client.Close()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Execute a query that takes longer than the timeout
	// pg_sleep(1) sleeps for 1 second
	_, err = client.ExecuteQuery(ctx, "SELECT pg_sleep(1)", nil)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Verify it's a timeout error
	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "canceled") {
		t.Errorf("Expected timeout or canceled error, got: %v", err)
	}

	t.Logf("Query timeout error: %v", err)
}

func TestClient_Integration_ExecuteQuery_SyntaxError(t *testing.T) {
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

	// Execute query with syntax error
	_, err = client.ExecuteQuery(ctx, "SELECT * FROM WHERE", nil)

	if err == nil {
		t.Fatal("Expected syntax error, got nil")
	}

	// Verify error message mentions syntax
	if !strings.Contains(err.Error(), "syntax") {
		t.Errorf("Expected 'syntax' in error message, got: %v", err)
	}

	t.Logf("Syntax error: %v", err)
}

func TestClient_Integration_ExecuteQuery_TableNotFound(t *testing.T) {
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

	// Execute query on non-existent table
	_, err = client.ExecuteQuery(ctx, "SELECT * FROM nonexistent_table_12345", nil)

	if err == nil {
		t.Fatal("Expected table not found error, got nil")
	}

	// Verify error message mentions table
	if !strings.Contains(err.Error(), "table") && !strings.Contains(err.Error(), "exist") {
		t.Errorf("Expected 'table' or 'exist' in error message, got: %v", err)
	}

	t.Logf("Table not found error: %v", err)
}

func TestClient_Integration_ExecuteQuery_EmptyResult(t *testing.T) {
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

	// Execute query that returns no rows
	result, err := client.ExecuteQuery(ctx, "SELECT 1 as n WHERE false", nil)
	if err != nil {
		t.Fatalf("ExecuteQuery() failed: %v", err)
	}

	if result.RowCount != 0 {
		t.Errorf("Expected 0 rows, got %d", result.RowCount)
	}

	if len(result.Rows) != 0 {
		t.Errorf("Expected empty rows array, got %d rows", len(result.Rows))
	}

	// Should still have column metadata
	if len(result.Columns) != 1 {
		t.Errorf("Expected 1 column in metadata, got %d", len(result.Columns))
	}
}

func TestClient_Integration_ExecuteQuery_InsertUpdate(t *testing.T) {
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

	// Create a temporary table
	_, err = client.ExecuteQuery(ctx, "CREATE TEMP TABLE test_table (id serial PRIMARY KEY, name text)", nil)
	if err != nil {
		t.Fatalf("Failed to create temp table: %v", err)
	}

	// Insert data
	result, err := client.ExecuteQuery(ctx, "INSERT INTO test_table (name) VALUES ('test1'), ('test2') RETURNING id, name", nil)
	if err != nil {
		t.Fatalf("INSERT query failed: %v", err)
	}

	if result.RowCount != 2 {
		t.Errorf("Expected 2 rows inserted, got %d", result.RowCount)
	}

	// Update data
	_, err = client.ExecuteQuery(ctx, "UPDATE test_table SET name = 'updated' WHERE name = 'test1'", nil)
	if err != nil {
		t.Fatalf("UPDATE query failed: %v", err)
	}

	// Verify update
	result, err = client.ExecuteQuery(ctx, "SELECT name FROM test_table WHERE name = 'updated'", nil)
	if err != nil {
		t.Fatalf("SELECT after UPDATE failed: %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("Expected 1 updated row, got %d", result.RowCount)
	}

	// Delete data
	_, err = client.ExecuteQuery(ctx, "DELETE FROM test_table", nil)
	if err != nil {
		t.Fatalf("DELETE query failed: %v", err)
	}
}

// Benchmark for ExecuteQuery

func BenchmarkClient_ExecuteQuery(b *testing.B) {
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
		_, err := client.ExecuteQuery(ctx, "SELECT 1", nil)
		if err != nil {
			b.Fatalf("ExecuteQuery() failed: %v", err)
		}
	}
}

func ExampleClient_ExecuteQuery() {
	ctx := context.Background()
	client, err := NewClient(ctx, "postgres://user:pass@localhost:5432/mydb")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	result, err := client.ExecuteQuery(ctx, "SELECT * FROM users WHERE id = $1", []interface{}{42})
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return
	}

	fmt.Printf("Query returned %d rows in %v\n", result.RowCount, result.ExecutionTime)
	for _, row := range result.Rows {
		fmt.Printf("Row: %v\n", row)
	}
}

// Tests for IntrospectSchema

func TestClient_Integration_IntrospectSchema_EmptyDatabase(t *testing.T) {
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

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Verify schema is not nil
	if schema == nil {
		t.Fatal("Expected schema to be non-nil")
	}

	// In a fresh database, there might be no user tables
	// but the query should still succeed
	t.Logf("Found %d tables and %d functions", len(schema.Tables), len(schema.Functions))
}

func TestClient_Integration_IntrospectSchema_WithTables(t *testing.T) {
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

	// Create a test table
	_, err = client.ExecuteQuery(ctx, `
		CREATE TEMP TABLE test_users (
			id serial PRIMARY KEY,
			username text NOT NULL,
			email text,
			created_at timestamp DEFAULT NOW()
		)
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Find the test_users table
	var testUsersTable *protocol.TableInfo
	for i := range schema.Tables {
		if schema.Tables[i].Name == "test_users" {
			testUsersTable = &schema.Tables[i]
			break
		}
	}

	if testUsersTable == nil {
		t.Fatal("Expected to find test_users table in schema")
	}

	// Verify table properties
	if testUsersTable.Type != "table" {
		t.Errorf("Expected table type to be 'table', got %s", testUsersTable.Type)
	}

	// Verify columns
	expectedColumns := map[string]bool{
		"id":         false,
		"username":   false,
		"email":      false,
		"created_at": false,
	}

	if len(testUsersTable.Columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(testUsersTable.Columns))
	}

	for _, col := range testUsersTable.Columns {
		if _, exists := expectedColumns[col.Name]; !exists {
			t.Errorf("Unexpected column: %s", col.Name)
		}
		expectedColumns[col.Name] = true

		// Verify column has data type
		if col.DataType == "" {
			t.Errorf("Column %s should have data type", col.Name)
		}

		// Verify specific column properties
		if col.Name == "username" {
			if col.Nullable {
				t.Error("Column 'username' should not be nullable")
			}
		}
		if col.Name == "email" {
			if !col.Nullable {
				t.Error("Column 'email' should be nullable")
			}
		}

		t.Logf("Column: %s, Type: %s, Nullable: %v", col.Name, col.DataType, col.Nullable)
	}

	// Verify all expected columns were found
	for colName, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column %s not found", colName)
		}
	}
}

func TestClient_Integration_IntrospectSchema_WithViews(t *testing.T) {
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

	// Create a test table
	_, err = client.ExecuteQuery(ctx, `
		CREATE TEMP TABLE test_products (
			id serial PRIMARY KEY,
			name text NOT NULL,
			price numeric(10,2)
		)
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Create a test view
	_, err = client.ExecuteQuery(ctx, `
		CREATE TEMP VIEW test_expensive_products AS
		SELECT id, name, price FROM test_products WHERE price > 100
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test view: %v", err)
	}

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Find the test view
	var testView *protocol.TableInfo
	for i := range schema.Tables {
		if schema.Tables[i].Name == "test_expensive_products" {
			testView = &schema.Tables[i]
			break
		}
	}

	if testView == nil {
		t.Fatal("Expected to find test_expensive_products view in schema")
	}

	// Verify it's identified as a view
	if testView.Type != "view" {
		t.Errorf("Expected type to be 'view', got %s", testView.Type)
	}

	// Verify view has columns
	if len(testView.Columns) != 3 {
		t.Errorf("Expected 3 columns in view, got %d", len(testView.Columns))
	}

	t.Logf("View: %s.%s (type: %s) with %d columns", testView.Schema, testView.Name, testView.Type, len(testView.Columns))
}

func TestClient_Integration_IntrospectSchema_WithFunctions(t *testing.T) {
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

	// Create a test function
	_, err = client.ExecuteQuery(ctx, `
		CREATE OR REPLACE FUNCTION test_add_numbers(a integer, b integer)
		RETURNS integer AS $$
		BEGIN
			RETURN a + b;
		END;
		$$ LANGUAGE plpgsql
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test function: %v", err)
	}

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Find the test function
	var testFunc *protocol.FunctionInfo
	for i := range schema.Functions {
		if schema.Functions[i].Name == "test_add_numbers" {
			testFunc = &schema.Functions[i]
			break
		}
	}

	if testFunc == nil {
		t.Fatal("Expected to find test_add_numbers function in schema")
	}

	// Verify function properties
	if testFunc.ReturnType == "" {
		t.Error("Expected function to have return type")
	}

	t.Logf("Function: %s.%s returns %s", testFunc.Schema, testFunc.Name, testFunc.ReturnType)

	// Clean up
	_, err = client.ExecuteQuery(ctx, "DROP FUNCTION IF EXISTS test_add_numbers(integer, integer)", nil)
	if err != nil {
		t.Logf("Warning: Failed to clean up test function: %v", err)
	}
}

func TestClient_Integration_IntrospectSchema_MultipleSchemas(t *testing.T) {
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

	// Create a custom schema
	_, err = client.ExecuteQuery(ctx, "CREATE SCHEMA IF NOT EXISTS test_schema", nil)
	if err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}
	defer func() {
		_, _ = client.ExecuteQuery(ctx, "DROP SCHEMA IF EXISTS test_schema CASCADE", nil)
	}()

	// Create a table in the custom schema
	_, err = client.ExecuteQuery(ctx, `
		CREATE TABLE test_schema.test_table (
			id serial PRIMARY KEY,
			data text
		)
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create table in test schema: %v", err)
	}

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Find the table in the custom schema
	var testTable *protocol.TableInfo
	for i := range schema.Tables {
		if schema.Tables[i].Schema == "test_schema" && schema.Tables[i].Name == "test_table" {
			testTable = &schema.Tables[i]
			break
		}
	}

	if testTable == nil {
		t.Fatal("Expected to find test_schema.test_table in schema")
	}

	// Verify schema name
	if testTable.Schema != "test_schema" {
		t.Errorf("Expected schema to be 'test_schema', got %s", testTable.Schema)
	}

	t.Logf("Found table: %s.%s", testTable.Schema, testTable.Name)
}

func TestClient_Integration_IntrospectSchema_DataTypes(t *testing.T) {
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

	// Create a table with various data types
	_, err = client.ExecuteQuery(ctx, `
		CREATE TEMP TABLE test_datatypes (
			int_col integer,
			text_col text,
			bool_col boolean,
			timestamp_col timestamp,
			json_col jsonb,
			uuid_col uuid,
			numeric_col numeric(10,2),
			array_col text[]
		)
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Find the test table
	var testTable *protocol.TableInfo
	for i := range schema.Tables {
		if schema.Tables[i].Name == "test_datatypes" {
			testTable = &schema.Tables[i]
			break
		}
	}

	if testTable == nil {
		t.Fatal("Expected to find test_datatypes table in schema")
	}

	// Verify all columns have data types
	expectedColumns := []string{"int_col", "text_col", "bool_col", "timestamp_col", "json_col", "uuid_col", "numeric_col", "array_col"}
	if len(testTable.Columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(testTable.Columns))
	}

	for _, col := range testTable.Columns {
		if col.DataType == "" {
			t.Errorf("Column %s should have data type", col.Name)
		}
		if col.TypeOID == 0 {
			t.Errorf("Column %s should have type OID", col.Name)
		}
		t.Logf("Column: %s, Type: %s, OID: %d, Nullable: %v", col.Name, col.DataType, col.TypeOID, col.Nullable)
	}
}

func TestClient_Integration_IntrospectSchema_EmptyTableColumns(t *testing.T) {
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

	// Create a minimal table
	_, err = client.ExecuteQuery(ctx, "CREATE TEMP TABLE test_minimal (id serial)", nil)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema() failed: %v", err)
	}

	// Find the test table
	var testTable *protocol.TableInfo
	for i := range schema.Tables {
		if schema.Tables[i].Name == "test_minimal" {
			testTable = &schema.Tables[i]
			break
		}
	}

	if testTable == nil {
		t.Fatal("Expected to find test_minimal table in schema")
	}

	// Verify it has at least the id column
	if len(testTable.Columns) == 0 {
		t.Error("Expected at least one column")
	}

	// Verify Columns array is never nil (should be empty array instead)
	if testTable.Columns == nil {
		t.Error("Columns array should not be nil")
	}
}

// Benchmark for IntrospectSchema

func BenchmarkClient_IntrospectSchema(b *testing.B) {
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

	// Create some test data
	_, _ = client.ExecuteQuery(ctx, "CREATE TEMP TABLE bench_table1 (id serial, data text)", nil)
	_, _ = client.ExecuteQuery(ctx, "CREATE TEMP TABLE bench_table2 (id serial, value integer)", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.IntrospectSchema(ctx)
		if err != nil {
			b.Fatalf("IntrospectSchema() failed: %v", err)
		}
	}
}

func ExampleClient_IntrospectSchema() {
	ctx := context.Background()
	client, err := NewClient(ctx, "postgres://user:pass@localhost:5432/mydb")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		fmt.Printf("Failed to introspect schema: %v\n", err)
		return
	}

	fmt.Printf("Found %d tables and %d functions\n", len(schema.Tables), len(schema.Functions))
	for _, table := range schema.Tables {
		fmt.Printf("Table: %s.%s (%s) with %d columns\n",
			table.Schema, table.Name, table.Type, len(table.Columns))
	}
}

// TestConvertValue tests the convertValue helper function
func TestConvertValue(t *testing.T) {
	client := &Client{} // Don't need a real connection for this test

	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "integer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: true,
		},
		{
			name:     "time value",
			input:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: "2024-01-01T12:00:00Z",
		},
		{
			name:     "byte array",
			input:    []byte("binary data"),
			expected: "binary data",
		},
		{
			name:     "empty byte array",
			input:    []byte{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := client.convertValue(tc.input)
			if result != tc.expected {
				t.Errorf("convertValue(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestGetDataTypeName tests the getDataTypeName helper function
func TestGetDataTypeName(t *testing.T) {
	client := &Client{} // Don't need a real connection for this test

	testCases := []struct {
		name     string
		oid      uint32
		expected string
	}{
		{
			name:     "bool type",
			oid:      16,
			expected: "bool",
		},
		{
			name:     "int2 type",
			oid:      21,
			expected: "int2",
		},
		{
			name:     "int4 type",
			oid:      23,
			expected: "int4",
		},
		{
			name:     "int8 type",
			oid:      20,
			expected: "int8",
		},
		{
			name:     "text type",
			oid:      25,
			expected: "text",
		},
		{
			name:     "varchar type",
			oid:      1043,
			expected: "varchar",
		},
		{
			name:     "timestamp type",
			oid:      1114,
			expected: "timestamp",
		},
		{
			name:     "timestamptz type",
			oid:      1184,
			expected: "timestamptz",
		},
		{
			name:     "date type",
			oid:      1082,
			expected: "date",
		},
		{
			name:     "time type",
			oid:      1083,
			expected: "time",
		},
		{
			name:     "json type",
			oid:      114,
			expected: "json",
		},
		{
			name:     "jsonb type",
			oid:      3802,
			expected: "jsonb",
		},
		{
			name:     "uuid type",
			oid:      2950,
			expected: "uuid",
		},
		{
			name:     "numeric type",
			oid:      1700,
			expected: "numeric",
		},
		{
			name:     "float4 type",
			oid:      700,
			expected: "float4",
		},
		{
			name:     "float8 type",
			oid:      701,
			expected: "float8",
		},
		{
			name:     "bytea type",
			oid:      17,
			expected: "bytea",
		},
		{
			name:     "unknown OID",
			oid:      99999,
			expected: "unknown(99999)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := client.getDataTypeName(tc.oid)
			if result != tc.expected {
				t.Errorf("getDataTypeName(%d) = %s, want %s", tc.oid, result, tc.expected)
			}
		})
	}
}

// TestHandleQueryError tests the handleQueryError helper function
func TestHandleQueryError(t *testing.T) {
	client := &Client{} // Don't need a real connection for this test

	testCases := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "context deadline exceeded",
			inputErr:    context.DeadlineExceeded,
			expectedMsg: "query timeout exceeded",
		},
		{
			name:        "generic error",
			inputErr:    fmt.Errorf("some generic error"),
			expectedMsg: "query failed: some generic error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := client.handleQueryError(tc.inputErr)
			if !strings.Contains(result.Error(), tc.expectedMsg) {
				t.Errorf("handleQueryError() = %v, want error containing %q", result, tc.expectedMsg)
			}
		})
	}
}
