# Integration Tests

This document describes how to run the integration tests for the PostgresMaster proxy.

## Overview

The integration tests verify the complete functionality of the proxy by connecting to a real Postgres database. These tests cover:

- Full query execution flow (SELECT, INSERT, UPDATE, DELETE)
- Schema introspection (tables, views, functions)
- Connection retry logic
- Error handling (syntax errors, timeouts, missing tables, etc.)
- Data type handling (integers, text, booleans, timestamps, JSON, UUID, arrays, etc.)

## Prerequisites

- **Docker**: Required to run the test Postgres database
- **Docker Compose**: Used to manage the test container
- **Go 1.21+**: Required to run the tests

## Quick Start

### Option 1: Using Make (Recommended)

The easiest way to run integration tests is using the provided Makefile target:

```bash
cd proxy
make test-integration
```

This will:
1. Start a Postgres container using Docker Compose
2. Wait for the database to be ready
3. Run all integration tests
4. Clean up the container

### Option 2: Using the Script Directly

You can also run the script directly:

```bash
cd proxy
./scripts/run-integration-tests.sh
```

### Option 3: Manual Setup

If you want more control, you can manually manage the test database:

```bash
# Start the test database
cd proxy
docker-compose -f docker-compose.test.yml up -d

# Wait for it to be ready (check with docker ps)
docker ps

# Set the connection URL
export TEST_POSTGRES_URL="postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"

# Run the integration tests
go test -tags=integration -v ./pkg/postgres/...

# When done, stop the database
docker-compose -f docker-compose.test.yml down -v
```

## Running All Tests

To run both unit tests and integration tests:

```bash
make test-all
```

This will run:
1. Unit tests (all packages)
2. Integration tests (with Docker)

## Test Coverage

With integration tests, the coverage should reach **80%+** as documented in the PRD.

To check coverage:

```bash
# With TEST_POSTGRES_URL set
export TEST_POSTGRES_URL="postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

## Test Database Configuration

The test database uses:
- **Image**: `postgres:16-alpine`
- **Port**: `5433` (to avoid conflicts with default 5432)
- **Database**: `testdb`
- **Username**: `testuser`
- **Password**: `testpass`
- **Storage**: tmpfs (in-memory, fast and ephemeral)

Configuration is in `docker-compose.test.yml`.

## Continuous Integration

In CI environments, you can run integration tests by:

1. Starting the Postgres service
2. Setting `TEST_POSTGRES_URL`
3. Running the tests

Example GitHub Actions workflow:

```yaml
- name: Start Postgres
  run: |
    docker-compose -f proxy/docker-compose.test.yml up -d
    sleep 5

- name: Run Integration Tests
  env:
    TEST_POSTGRES_URL: postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable
  run: |
    cd proxy
    go test -tags=integration -v ./pkg/postgres/...
```

## Test Files

- `pkg/postgres/client_test.go` - Unit tests and integration tests (skipped without TEST_POSTGRES_URL)
- `pkg/postgres/integration_test.go` - Dedicated integration tests (build tag: `integration`)
- `scripts/run-integration-tests.sh` - Automated test runner
- `docker-compose.test.yml` - Test database configuration

## Troubleshooting

### Docker not running

```
Error: Docker is not running
```

**Solution**: Start Docker Desktop or the Docker daemon.

### Port already in use

```
Error: port 5433 is already allocated
```

**Solution**: Either stop the service using port 5433, or modify `docker-compose.test.yml` to use a different port.

### Tests timeout

If tests timeout, the database might not be ready. The script waits up to 30 seconds, but you can adjust this in `run-integration-tests.sh`.

### Permission denied on script

```
Permission denied: ./scripts/run-integration-tests.sh
```

**Solution**: Make the script executable:
```bash
chmod +x scripts/run-integration-tests.sh
```

## Development Workflow

When developing new features:

1. Write unit tests first (no database required)
2. Run unit tests: `make test`
3. Write integration tests for database interactions
4. Run integration tests: `make test-integration`
5. Verify coverage: `make test-coverage`

## Notes

- Integration tests are automatically skipped if `TEST_POSTGRES_URL` is not set
- The build tag `integration` allows running only integration tests: `go test -tags=integration ./...`
- All integration tests use temporary tables/views/functions to avoid conflicts
- The test database is ephemeral and destroyed after each test run
