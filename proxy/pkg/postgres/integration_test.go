//go:build integration
// +build integration

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

// Integration tests for Postgres client
// These tests require Docker to be running
// Run with: go test -tags=integration -v ./pkg/postgres

// TestIntegration_FullQueryExecutionFlow tests the complete query execution flow
func TestIntegration_FullQueryExecutionFlow(t *testing.T) {
	url := getIntegrationTestURL(t)

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer client.Close()

	// Test 1: Simple SELECT query
	t.Run("SimpleSelect", func(t *testing.T) {
		result, err := client.ExecuteQuery(ctx, "SELECT 1 as num, 'test' as text, true as flag", nil)
		if err != nil {
			t.Fatalf("ExecuteQuery failed: %v", err)
		}

		if result.RowCount != 1 {
			t.Errorf("Expected 1 row, got %d", result.RowCount)
		}

		if len(result.Columns) != 3 {
			t.Errorf("Expected 3 columns, got %d", len(result.Columns))
		}

		if result.ExecutionTime == 0 {
			t.Error("ExecutionTime should be non-zero")
		}

		row := result.Rows[0]
		if row["text"] != "test" {
			t.Errorf("Expected 'text' to be 'test', got %v", row["text"])
		}
	})

	// Test 2: Query with parameters
	t.Run("ParameterizedQuery", func(t *testing.T) {
		params := []interface{}{100, "param_value"}
		result, err := client.ExecuteQuery(ctx, "SELECT $1::int as num, $2::text as str", params)
		if err != nil {
			t.Fatalf("Parameterized query failed: %v", err)
		}

		if result.RowCount != 1 {
			t.Errorf("Expected 1 row, got %d", result.RowCount)
		}

		row := result.Rows[0]
		if numVal, ok := row["num"].(int32); !ok || numVal != 100 {
			t.Errorf("Expected 'num' to be 100, got %v (%T)", row["num"], row["num"])
		}
		if row["str"] != "param_value" {
			t.Errorf("Expected 'str' to be 'param_value', got %v", row["str"])
		}
	})

	// Test 3: Multiple rows
	t.Run("MultipleRows", func(t *testing.T) {
		result, err := client.ExecuteQuery(ctx, "SELECT generate_series(1, 50) as n", nil)
		if err != nil {
			t.Fatalf("Multiple rows query failed: %v", err)
		}

		if result.RowCount != 50 {
			t.Errorf("Expected 50 rows, got %d", result.RowCount)
		}

		if len(result.Rows) != 50 {
			t.Errorf("Expected 50 rows in result, got %d", len(result.Rows))
		}
	})

	// Test 4: NULL values
	t.Run("NullValues", func(t *testing.T) {
		result, err := client.ExecuteQuery(ctx, "SELECT NULL as null_col, 1 as int_col", nil)
		if err != nil {
			t.Fatalf("NULL query failed: %v", err)
		}

		row := result.Rows[0]
		if row["null_col"] != nil {
			t.Errorf("Expected 'null_col' to be nil, got %v", row["null_col"])
		}
	})

	// Test 5: CREATE, INSERT, UPDATE, DELETE
	t.Run("DataModification", func(t *testing.T) {
		// Create table
		_, err := client.ExecuteQuery(ctx, `
			CREATE TEMP TABLE integration_test (
				id serial PRIMARY KEY,
				name text NOT NULL,
				value integer
			)
		`, nil)
		if err != nil {
			t.Fatalf("CREATE TABLE failed: %v", err)
		}

		// Insert data
		result, err := client.ExecuteQuery(ctx, `
			INSERT INTO integration_test (name, value)
			VALUES ('test1', 100), ('test2', 200)
			RETURNING id, name, value
		`, nil)
		if err != nil {
			t.Fatalf("INSERT failed: %v", err)
		}
		if result.RowCount != 2 {
			t.Errorf("Expected 2 rows inserted, got %d", result.RowCount)
		}

		// Update data
		_, err = client.ExecuteQuery(ctx, `
			UPDATE integration_test
			SET value = 150
			WHERE name = 'test1'
		`, nil)
		if err != nil {
			t.Fatalf("UPDATE failed: %v", err)
		}

		// Verify update
		result, err = client.ExecuteQuery(ctx, "SELECT value FROM integration_test WHERE name = 'test1'", nil)
		if err != nil {
			t.Fatalf("SELECT after UPDATE failed: %v", err)
		}
		if result.RowCount != 1 {
			t.Errorf("Expected 1 row, got %d", result.RowCount)
		}
		if val, ok := result.Rows[0]["value"].(int32); !ok || val != 150 {
			t.Errorf("Expected value to be 150, got %v", result.Rows[0]["value"])
		}

		// Delete data
		_, err = client.ExecuteQuery(ctx, "DELETE FROM integration_test WHERE value > 100", nil)
		if err != nil {
			t.Fatalf("DELETE failed: %v", err)
		}

		// Verify delete
		result, err = client.ExecuteQuery(ctx, "SELECT COUNT(*) as count FROM integration_test", nil)
		if err != nil {
			t.Fatalf("SELECT COUNT failed: %v", err)
		}
		if count, ok := result.Rows[0]["count"].(int64); !ok || count != 1 {
			t.Errorf("Expected count to be 1, got %v", result.Rows[0]["count"])
		}
	})
}

// TestIntegration_SchemaIntrospection tests schema introspection functionality
func TestIntegration_SchemaIntrospection(t *testing.T) {
	url := getIntegrationTestURL(t)

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer client.Close()

	// Create test schema with tables, views, and functions
	setupTestSchema(t, client, ctx)

	// Introspect schema
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		t.Fatalf("IntrospectSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Test 1: Verify table information
	t.Run("TableInfo", func(t *testing.T) {
		var testTable *protocol.TableInfo
		for i := range schema.Tables {
			if schema.Tables[i].Name == "test_users" {
				testTable = &schema.Tables[i]
				break
			}
		}

		if testTable == nil {
			t.Fatal("Expected to find test_users table")
		}

		if testTable.Type != "table" {
			t.Errorf("Expected type 'table', got %s", testTable.Type)
		}

		// Verify columns
		expectedCols := map[string]bool{
			"id":         false,
			"username":   false,
			"email":      false,
			"created_at": false,
		}

		for _, col := range testTable.Columns {
			if _, exists := expectedCols[col.Name]; exists {
				expectedCols[col.Name] = true

				// Verify column metadata
				if col.DataType == "" {
					t.Errorf("Column %s should have data type", col.Name)
				}
				if col.TypeOID == 0 {
					t.Errorf("Column %s should have type OID", col.Name)
				}

				// Verify NOT NULL constraint
				if col.Name == "username" && col.Nullable {
					t.Error("Column 'username' should not be nullable")
				}
				if col.Name == "email" && !col.Nullable {
					t.Error("Column 'email' should be nullable")
				}
			}
		}

		// Verify all expected columns found
		for colName, found := range expectedCols {
			if !found {
				t.Errorf("Expected column %s not found", colName)
			}
		}
	})

	// Test 2: Verify view information
	t.Run("ViewInfo", func(t *testing.T) {
		var testView *protocol.TableInfo
		for i := range schema.Tables {
			if schema.Tables[i].Name == "test_active_users" {
				testView = &schema.Tables[i]
				break
			}
		}

		if testView == nil {
			t.Fatal("Expected to find test_active_users view")
		}

		if testView.Type != "view" {
			t.Errorf("Expected type 'view', got %s", testView.Type)
		}

		if len(testView.Columns) == 0 {
			t.Error("View should have columns")
		}
	})

	// Test 3: Verify function information
	t.Run("FunctionInfo", func(t *testing.T) {
		var testFunc *protocol.FunctionInfo
		for i := range schema.Functions {
			if schema.Functions[i].Name == "test_get_user_count" {
				testFunc = &schema.Functions[i]
				break
			}
		}

		if testFunc == nil {
			t.Fatal("Expected to find test_get_user_count function")
		}

		if testFunc.ReturnType == "" {
			t.Error("Function should have return type")
		}

		if testFunc.Schema == "" {
			t.Error("Function should have schema name")
		}
	})

	// Test 4: Verify data types
	t.Run("DataTypes", func(t *testing.T) {
		var datatypesTable *protocol.TableInfo
		for i := range schema.Tables {
			if schema.Tables[i].Name == "test_datatypes" {
				datatypesTable = &schema.Tables[i]
				break
			}
		}

		if datatypesTable == nil {
			t.Fatal("Expected to find test_datatypes table")
		}

		expectedTypes := []string{
			"int_col", "text_col", "bool_col", "timestamp_col",
			"json_col", "uuid_col", "numeric_col", "array_col",
		}

		if len(datatypesTable.Columns) != len(expectedTypes) {
			t.Errorf("Expected %d columns, got %d", len(expectedTypes), len(datatypesTable.Columns))
		}

		for _, col := range datatypesTable.Columns {
			if col.DataType == "" {
				t.Errorf("Column %s should have data type", col.Name)
			}
			if col.TypeOID == 0 {
				t.Errorf("Column %s should have type OID", col.Name)
			}
		}
	})
}

// TestIntegration_ConnectionRetryLogic tests connection retry behavior
func TestIntegration_ConnectionRetryLogic(t *testing.T) {
	ctx := context.Background()

	t.Run("SuccessfulConnection", func(t *testing.T) {
		url := getIntegrationTestURL(t)
		client, err := NewClient(ctx, url)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		defer client.Close()

		// Verify we can ping
		if err := client.Ping(ctx); err != nil {
			t.Errorf("Ping failed: %v", err)
		}
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		badURL := "postgres://baduser:badpass@localhost:5432/testdb?connect_timeout=1"
		start := time.Now()
		client, err := NewClient(ctx, badURL)
		duration := time.Since(start)

		if err == nil {
			t.Error("Expected connection to fail with bad credentials")
			if client != nil {
				client.Close()
			}
			return
		}

		// Verify error mentions retry attempts
		if !strings.Contains(err.Error(), "4 attempts") {
			t.Errorf("Error should mention retry attempts: %v", err)
		}

		// Verify retry logic executed (should take some time)
		if duration < 2*time.Second {
			t.Logf("Connection failed quickly (%v), may not have executed full retry logic", duration)
		}

		t.Logf("Connection with bad credentials failed after %v (expected)", duration)
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		// Use unreachable IP to trigger timeout
		badURL := "postgres://user:pass@192.0.2.1:5432/db?connect_timeout=1"
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after short delay
		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		start := time.Now()
		client, err := NewClient(ctx, badURL)
		duration := time.Since(start)

		if err == nil {
			t.Error("Expected connection to fail after context cancellation")
			if client != nil {
				client.Close()
			}
			return
		}

		// Should fail relatively quickly after cancellation
		if duration > 5*time.Second {
			t.Errorf("Connection took too long (%v) after context cancellation", duration)
		}

		t.Logf("Connection cancelled after %v (expected)", duration)
	})
}

// TestIntegration_ErrorHandling tests various error scenarios
func TestIntegration_ErrorHandling(t *testing.T) {
	url := getIntegrationTestURL(t)

	ctx := context.Background()
	client, err := NewClient(ctx, url)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	t.Run("SyntaxError", func(t *testing.T) {
		_, err := client.ExecuteQuery(ctx, "SELECT * FROM WHERE", nil)
		if err == nil {
			t.Fatal("Expected syntax error")
		}
		if !strings.Contains(err.Error(), "syntax") {
			t.Errorf("Error should mention 'syntax', got: %v", err)
		}
	})

	t.Run("TableNotFound", func(t *testing.T) {
		_, err := client.ExecuteQuery(ctx, "SELECT * FROM nonexistent_table_xyz123", nil)
		if err == nil {
			t.Fatal("Expected table not found error")
		}
		errStr := err.Error()
		if !strings.Contains(errStr, "table") && !strings.Contains(errStr, "exist") && !strings.Contains(errStr, "relation") {
			t.Errorf("Error should mention table/exist/relation, got: %v", err)
		}
	})

	t.Run("QueryTimeout", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := client.ExecuteQuery(timeoutCtx, "SELECT pg_sleep(10)", nil)
		if err == nil {
			t.Fatal("Expected timeout error")
		}

		errStr := err.Error()
		if !strings.Contains(errStr, "timeout") && !strings.Contains(errStr, "canceled") && !strings.Contains(errStr, "context") {
			t.Errorf("Error should mention timeout/canceled/context, got: %v", err)
		}
	})

	t.Run("InvalidParameterCount", func(t *testing.T) {
		// Query expects 2 parameters, but we provide 1
		_, err := client.ExecuteQuery(ctx, "SELECT $1, $2", []interface{}{1})
		if err == nil {
			t.Fatal("Expected parameter count error")
		}
		// Error message varies, so just verify we got an error
		t.Logf("Got expected error: %v", err)
	})

	t.Run("EmptyQuery", func(t *testing.T) {
		_, err := client.ExecuteQuery(ctx, "", nil)
		if err == nil {
			t.Fatal("Expected error for empty query")
		}
	})

	t.Run("PingAfterClose", func(t *testing.T) {
		// Create a separate client for this test
		testClient, err := NewClient(ctx, url)
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		testClient.Close()

		// Ping should fail after close
		err = testClient.Ping(ctx)
		if err == nil {
			t.Error("Expected Ping to fail after Close")
		}
	})
}

// Helper functions

func getIntegrationTestURL(t *testing.T) string {
	url := os.Getenv("TEST_POSTGRES_URL")
	if url == "" {
		t.Skip("Skipping integration test: TEST_POSTGRES_URL not set")
	}
	return url
}

func setupTestSchema(t *testing.T, client *Client, ctx context.Context) {
	// Create test table
	_, err := client.ExecuteQuery(ctx, `
		CREATE TEMP TABLE test_users (
			id serial PRIMARY KEY,
			username text NOT NULL,
			email text,
			created_at timestamp DEFAULT NOW()
		)
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test_users table: %v", err)
	}

	// Create test view
	_, err = client.ExecuteQuery(ctx, `
		CREATE TEMP VIEW test_active_users AS
		SELECT id, username, email FROM test_users WHERE username IS NOT NULL
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test_active_users view: %v", err)
	}

	// Create test function
	_, err = client.ExecuteQuery(ctx, `
		CREATE OR REPLACE FUNCTION test_get_user_count()
		RETURNS integer AS $$
		BEGIN
			RETURN (SELECT COUNT(*) FROM test_users);
		END;
		$$ LANGUAGE plpgsql
	`, nil)
	if err != nil {
		t.Fatalf("Failed to create test_get_user_count function: %v", err)
	}

	// Create table with various data types
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
		t.Fatalf("Failed to create test_datatypes table: %v", err)
	}

	// Insert some test data
	_, err = client.ExecuteQuery(ctx, `
		INSERT INTO test_users (username, email)
		VALUES
			('alice', 'alice@example.com'),
			('bob', 'bob@example.com'),
			('charlie', NULL)
	`, nil)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}
