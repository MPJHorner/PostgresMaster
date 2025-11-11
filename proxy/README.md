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

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
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
