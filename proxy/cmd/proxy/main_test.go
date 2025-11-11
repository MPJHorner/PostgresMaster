package main

import (
	"testing"
)

func TestValidateConnectionString(t *testing.T) {
	tests := []struct {
		name      string
		connStr   string
		wantError bool
	}{
		{
			name:      "valid postgres connection string",
			connStr:   "postgres://user:pass@localhost:5432/testdb",
			wantError: false,
		},
		{
			name:      "valid postgresql connection string",
			connStr:   "postgresql://user:pass@localhost:5432/testdb",
			wantError: false,
		},
		{
			name:      "valid with query params",
			connStr:   "postgres://user:pass@localhost:5432/testdb?sslmode=disable",
			wantError: false,
		},
		{
			name:      "empty string",
			connStr:   "",
			wantError: true,
		},
		{
			name:      "invalid scheme",
			connStr:   "mysql://user:pass@localhost:5432/testdb",
			wantError: true,
		},
		{
			name:      "missing host",
			connStr:   "postgres://user:pass@/testdb",
			wantError: true,
		},
		{
			name:      "missing database",
			connStr:   "postgres://user:pass@localhost:5432",
			wantError: true,
		},
		{
			name:      "missing database with slash",
			connStr:   "postgres://user:pass@localhost:5432/",
			wantError: true,
		},
		{
			name:      "invalid URL format",
			connStr:   "not a valid url",
			wantError: true,
		},
		{
			name:      "valid with special characters in password",
			connStr:   "postgres://user:p%40ss%23word@localhost:5432/testdb",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConnectionString(tt.connStr)
			if (err != nil) != tt.wantError {
				t.Errorf("validateConnectionString() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
