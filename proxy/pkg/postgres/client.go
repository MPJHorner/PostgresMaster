package postgres

import (
	"context"
	"fmt"
	"time"

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
