package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MPJHorner/PostgresMaster/proxy/pkg/protocol"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client represents a connection to a PostgreSQL database
type Client struct {
	pool *pgxpool.Pool
}

// NewClient creates a new Postgres client with connection pooling and retry logic
func NewClient(ctx context.Context, connString string) (*Client, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 5
	config.MinConns = 1

	// Retry logic with exponential backoff
	maxAttempts := 4
	backoffDurations := []time.Duration{0, 2 * time.Second, 4 * time.Second, 8 * time.Second}

	var pool *pgxpool.Pool
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Wait before retry (except first attempt)
		if attempt > 1 {
			waitDuration := backoffDurations[attempt-1]
			fmt.Printf("Retrying... (attempt %d/%d) after %v\n", attempt, maxAttempts, waitDuration)
			select {
			case <-time.After(waitDuration):
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			}
		} else {
			fmt.Printf("Connecting to database... (attempt %d/%d)\n", attempt, maxAttempts)
		}

		// Attempt connection
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			lastErr = fmt.Errorf("failed to create connection pool: %w", err)
			continue
		}

		// Test connection with ping
		if err := pool.Ping(ctx); err != nil {
			lastErr = fmt.Errorf("failed to ping database: %w", err)
			pool.Close()
			continue
		}

		// Success!
		fmt.Println("✓ Connected!")
		return &Client{pool: pool}, nil
	}

	// All attempts failed
	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxAttempts, lastErr)
}

// Close closes the database connection pool
func (c *Client) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

// Ping tests the database connection
func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// QueryResult contains the results of a query execution
type QueryResult struct {
	Rows          []map[string]interface{}
	Columns       []protocol.ColumnInfo
	RowCount      int
	ExecutionTime time.Duration
}

// ExecuteQuery executes a SQL query and returns the results
func (c *Client) ExecuteQuery(ctx context.Context, sql string, params []interface{}) (*QueryResult, error) {
	// Measure execution time
	startTime := time.Now()

	// Execute the query
	rows, err := c.pool.Query(ctx, sql, params...)
	if err != nil {
		return nil, c.handleQueryError(err)
	}
	defer rows.Close()

	// Parse field descriptions (column metadata)
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]protocol.ColumnInfo, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = protocol.ColumnInfo{
			Name:     string(fd.Name),
			DataType: c.getDataTypeName(fd.DataTypeOID),
			TypeOID:  fd.DataTypeOID,
		}
	}

	// Parse result rows
	resultRows := []map[string]interface{}{}
	for rows.Next() {
		// Get values for this row
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("failed to read row values: %w", err)
		}

		// Build row map
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col.Name] = c.convertValue(values[i])
		}
		resultRows = append(resultRows, rowMap)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	executionTime := time.Since(startTime)

	return &QueryResult{
		Rows:          resultRows,
		Columns:       columns,
		RowCount:      len(resultRows),
		ExecutionTime: executionTime,
	}, nil
}

// convertValue converts database values to JSON-friendly types
func (c *Client) convertValue(value interface{}) interface{} {
	// Handle NULL values
	if value == nil {
		return nil
	}

	// Handle time types - convert to RFC3339 string
	switch v := value.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	case []byte:
		// Convert byte arrays to strings for JSON compatibility
		return string(v)
	default:
		// All other types are already JSON-compatible
		return value
	}
}

// getDataTypeName returns a human-readable name for a Postgres OID
func (c *Client) getDataTypeName(oid uint32) string {
	// Common Postgres type OIDs
	// Full list: https://github.com/postgres/postgres/blob/master/src/include/catalog/pg_type.dat
	typeNames := map[uint32]string{
		16:   "bool",
		17:   "bytea",
		18:   "char",
		19:   "name",
		20:   "int8",
		21:   "int2",
		23:   "int4",
		25:   "text",
		114:  "json",
		142:  "xml",
		194:  "pg_node_tree",
		700:  "float4",
		701:  "float8",
		705:  "unknown",
		790:  "money",
		829:  "macaddr",
		869:  "inet",
		1000: "_bool",
		1001: "_bytea",
		1002: "_char",
		1003: "_name",
		1005: "_int2",
		1007: "_int4",
		1009: "_text",
		1014: "_bpchar",
		1015: "_varchar",
		1016: "_int8",
		1021: "_float4",
		1022: "_float8",
		1042: "bpchar",
		1043: "varchar",
		1082: "date",
		1083: "time",
		1114: "timestamp",
		1115: "_timestamp",
		1182: "_date",
		1183: "_time",
		1184: "timestamptz",
		1185: "_timestamptz",
		1186: "interval",
		1187: "_interval",
		1231: "_numeric",
		1266: "timetz",
		1270: "_timetz",
		1560: "bit",
		1562: "varbit",
		1700: "numeric",
		2950: "uuid",
		3802: "jsonb",
	}

	if name, ok := typeNames[oid]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", oid)
}

// handleQueryError categorizes and formats query errors
func (c *Client) handleQueryError(err error) error {
	// Check if it's a pgconn error with code
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "42601": // Syntax error
			return fmt.Errorf("syntax error: %s", pgErr.Message)
		case "42501": // Insufficient privilege
			return fmt.Errorf("permission denied: %s", pgErr.Message)
		case "42P01": // Undefined table
			return fmt.Errorf("table does not exist: %s", pgErr.Message)
		case "42703": // Undefined column
			return fmt.Errorf("column does not exist: %s", pgErr.Message)
		case "57014": // Query canceled
			return fmt.Errorf("query canceled: %s", pgErr.Message)
		default:
			// Return the full Postgres error
			return fmt.Errorf("database error [%s]: %s", pgErr.Code, pgErr.Message)
		}
	}

	// Check for context timeout
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("query timeout exceeded")
	}

	// Return generic error
	return fmt.Errorf("query failed: %w", err)
}

// IntrospectSchema queries the database schema and returns information about tables and functions
func (c *Client) IntrospectSchema(ctx context.Context) (*protocol.SchemaPayload, error) {
	// Query for tables (including views and materialized views)
	tables, err := c.queryTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}

	// Query for columns for each table
	for i := range tables {
		columns, err := c.queryColumns(ctx, tables[i].Schema, tables[i].Name)
		if err != nil {
			return nil, fmt.Errorf("failed to query columns for %s.%s: %w", tables[i].Schema, tables[i].Name, err)
		}
		tables[i].Columns = columns
	}

	// Query for functions
	functions, err := c.queryFunctions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query functions: %w", err)
	}

	return &protocol.SchemaPayload{
		Tables:    tables,
		Functions: functions,
	}, nil
}

// queryTables retrieves all user-defined tables, views, and materialized views
func (c *Client) queryTables(ctx context.Context) ([]protocol.TableInfo, error) {
	query := `
		SELECT n.nspname, c.relname, c.relkind
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind IN ('r', 'v', 'm')
		  AND n.nspname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY n.nspname, c.relname
	`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []protocol.TableInfo
	for rows.Next() {
		var schema, name, kind string
		if err := rows.Scan(&schema, &name, &kind); err != nil {
			return nil, fmt.Errorf("failed to scan table row: %w", err)
		}

		// Map relkind to readable type
		var tableType string
		switch kind {
		case "r":
			tableType = "table"
		case "v":
			tableType = "view"
		case "m":
			tableType = "materialized view"
		default:
			tableType = kind
		}

		tables = append(tables, protocol.TableInfo{
			Schema:  schema,
			Name:    name,
			Type:    tableType,
			Columns: []protocol.ColumnInfo{}, // Will be filled later
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table rows: %w", err)
	}

	return tables, nil
}

// queryColumns retrieves all columns for a specific table
func (c *Client) queryColumns(ctx context.Context, schema, table string) ([]protocol.ColumnInfo, error) {
	query := `
		SELECT
			a.attname,
			format_type(a.atttypid, a.atttypmod) as type_name,
			NOT a.attnotnull as nullable,
			a.atttypid as type_oid
		FROM pg_attribute a
		WHERE a.attrelid = ($1 || '.' || $2)::regclass
		  AND a.attnum > 0
		  AND NOT a.attisdropped
		ORDER BY a.attnum
	`

	rows, err := c.pool.Query(ctx, query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []protocol.ColumnInfo
	for rows.Next() {
		var name, dataType string
		var nullable bool
		var typeOID uint32

		if err := rows.Scan(&name, &dataType, &nullable, &typeOID); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}

		columns = append(columns, protocol.ColumnInfo{
			Name:     name,
			DataType: dataType,
			TypeOID:  typeOID,
			Nullable: nullable,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating column rows: %w", err)
	}

	return columns, nil
}

// queryFunctions retrieves all user-defined functions
func (c *Client) queryFunctions(ctx context.Context) ([]protocol.FunctionInfo, error) {
	query := `
		SELECT
			n.nspname,
			p.proname,
			pg_get_function_result(p.oid) as return_type
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY n.nspname, p.proname
	`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var functions []protocol.FunctionInfo
	for rows.Next() {
		var schema, name, returnType string
		if err := rows.Scan(&schema, &name, &returnType); err != nil {
			return nil, fmt.Errorf("failed to scan function row: %w", err)
		}

		functions = append(functions, protocol.FunctionInfo{
			Schema:     schema,
			Name:       name,
			ReturnType: returnType,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating function rows: %w", err)
	}

	return functions, nil
}
