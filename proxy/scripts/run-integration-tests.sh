#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}PostgresMaster Integration Tests${NC}"
echo "=================================="
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    echo "Please start Docker and try again"
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null 2>&1; then
    echo -e "${RED}Error: Docker Compose is not installed${NC}"
    exit 1
fi

# Use docker-compose or docker compose based on availability
DOCKER_COMPOSE="docker compose"
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
fi

cd "$(dirname "$0")/.."

echo -e "${YELLOW}Starting test Postgres database...${NC}"
$DOCKER_COMPOSE -f docker-compose.test.yml up -d

echo -e "${YELLOW}Waiting for database to be ready...${NC}"
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if docker exec postgres-master-test pg_isready -U testuser -d testdb > /dev/null 2>&1; then
        echo -e "${GREEN}Database is ready!${NC}"
        break
    fi
    attempt=$((attempt + 1))
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${RED}Database failed to start within 30 seconds${NC}"
        $DOCKER_COMPOSE -f docker-compose.test.yml logs
        $DOCKER_COMPOSE -f docker-compose.test.yml down
        exit 1
    fi
    sleep 1
done

echo ""
echo -e "${YELLOW}Running integration tests...${NC}"
export TEST_POSTGRES_URL="postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"

# Run the integration tests
if go test -tags=integration -v ./pkg/postgres/...; then
    echo ""
    echo -e "${GREEN}All integration tests passed!${NC}"
    exit_code=0
else
    echo ""
    echo -e "${RED}Some integration tests failed${NC}"
    exit_code=1
fi

echo ""
echo -e "${YELLOW}Cleaning up test database...${NC}"
$DOCKER_COMPOSE -f docker-compose.test.yml down -v

exit $exit_code
