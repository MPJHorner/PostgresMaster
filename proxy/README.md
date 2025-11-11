# PostgreSQL Proxy

A lightweight WebSocket-to-PostgreSQL proxy that enables browser-based database clients to connect to PostgreSQL servers.

## Features

- **WebSocket Server**: Accepts connections from web applications
- **PostgreSQL Client**: Connects to PostgreSQL databases using the native protocol
- **Secure**: Secret-based authentication for WebSocket connections
- **Connection Pooling**: Efficient connection management
- **Cross-Platform**: Supports Windows, macOS (Intel & Apple Silicon), and Linux

## Installation

### From Source

```bash
# Build for your current platform
make build

# Build for all platforms
make build-all
```

Binaries will be created in the `bin/` directory.

## Usage

### Basic Usage

```bash
# Using connection string
./postgres-proxy "postgres://user:password@host:5432/database"
```

### Interactive Mode (Coming Soon)

```bash
# Start interactive mode
./postgres-proxy

# Follow the prompts to enter connection details
```

## Development

### Prerequisites

- Go 1.21 or later
- PostgreSQL database (for testing)

### Running in Development Mode

```bash
# Run without building
make dev "postgres://user:password@localhost:5432/database"
```

### Running Tests

The proxy includes comprehensive unit and integration tests. Unit tests run without requiring a database connection, while integration tests require a real PostgreSQL instance.

#### Unit Tests

Unit tests run by default and cover:
- Secret generation and validation (87.5% coverage)
- Message protocol serialization (100% coverage)
- WebSocket server logic (84.5% coverage)
- Database client helper functions (type conversion, error handling)

```bash
# Run all tests (unit tests only, integration tests skipped)
go test ./... -v

# Run with coverage report
go test ./... -coverprofile=coverage.out

# View coverage summary
go tool cover -func=coverage.out | tail -1

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

#### Integration Tests

Integration tests validate the full database interaction flow including:
- Query execution with various data types
- Schema introspection (tables, columns, functions)
- Error handling (syntax errors, timeouts, missing tables)
- Connection retry logic

To run integration tests, set the `TEST_POSTGRES_URL` environment variable:

```bash
# Start a test PostgreSQL instance (using Docker)
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=test postgres:latest

# Run tests with integration tests enabled
TEST_POSTGRES_URL="postgres://postgres:test@localhost:5432/postgres" go test ./... -v

# Or export it for multiple test runs
export TEST_POSTGRES_URL="postgres://postgres:test@localhost:5432/postgres"
go test ./... -v -coverprofile=coverage.out
```

#### Coverage Targets

- **Current Coverage**: ~35% (unit tests only)
- **With Integration Tests**: 80%+ (when TEST_POSTGRES_URL is set)
- **Per-Package Coverage** (unit tests only):
  - `pkg/protocol`: 100%
  - `pkg/auth`: 87.5%
  - `pkg/server`: 84.5%
  - `pkg/postgres`: 27.9% (most code requires database connection)
  - `cmd/proxy`: 10.6% (main function and interactive prompts)

#### Makefile Shortcuts

```bash
# Run unit tests
make test

# Run tests with race detector
go test ./... -race

# Run specific package tests
go test ./pkg/postgres -v

# Run specific test function
go test ./pkg/postgres -v -run TestConvertValue
```

### Code Quality

```bash
# Run linter
make lint

# Format code
go fmt ./...
```

## Project Structure

```
proxy/
├── cmd/
│   └── proxy/           # Main entry point
│       └── main.go
├── pkg/
│   ├── auth/           # Secret generation and validation
│   │   └── secret.go
│   ├── postgres/       # PostgreSQL client
│   │   └── client.go
│   ├── protocol/       # Message protocol definitions
│   │   └── messages.go
│   └── server/         # WebSocket server
│       └── websocket.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## WebSocket Protocol

The proxy uses a JSON-based protocol for communication between the browser and the proxy.

### Client Messages

```json
{
  "id": "unique-request-id",
  "type": "query|introspect|ping",
  "payload": {
    "sql": "SELECT * FROM users",
    "params": [],
    "timeout": 30000
  }
}
```

### Server Messages

```json
{
  "id": "unique-request-id",
  "type": "result|error|schema|pong",
  "payload": {
    "rows": [...],
    "columns": [...],
    "rowCount": 10,
    "executionTime": 45
  }
}
```

## Security

- All WebSocket connections require a valid secret passed as a query parameter
- Secrets are 64-character hex-encoded strings (32 bytes of cryptographic randomness)
- CORS is restricted to localhost origins only
- The proxy never stores or logs sensitive connection information

## Contributing

Contributions are welcome! Please see the main repository's CONTRIBUTING.md for guidelines.

## License

This project is licensed under the AGPL-3.0 License - see the LICENSE file in the root directory for details.
