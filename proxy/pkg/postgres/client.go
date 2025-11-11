package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client represents a connection to a PostgreSQL database
type Client struct {
	pool *pgxpool.Pool
}

// NewClient creates a new Postgres client with connection pooling
func NewClient(ctx context.Context, connString string) (*Client, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 5
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{pool: pool}, nil
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
