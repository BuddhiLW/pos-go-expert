#!/bin/bash

# Order System Setup Verification Script

set -e

echo "🔍 Verifying Order System Setup..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_pass() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_fail() {
    echo -e "${RED}❌ $1${NC}"
}

print_warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check Go installation
echo "📋 Checking prerequisites..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_pass "Go is installed: $GO_VERSION"
else
    print_fail "Go is not installed"
    exit 1
fi

# Check Docker installation
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    print_pass "Docker is installed: $DOCKER_VERSION"
else
    print_fail "Docker is not installed"
    exit 1
fi

# Check if Docker Compose is available
if command -v docker-compose &> /dev/null || docker compose version &> /dev/null 2>&1; then
    print_pass "Docker Compose is available"
else
    print_fail "Docker Compose is not available"
    exit 1
fi

# Check Go module
echo "📦 Checking Go module..."
if [ -f "go.mod" ]; then
    print_pass "go.mod exists"
    MODULE_NAME=$(grep "^module " go.mod | awk '{print $2}')
    print_pass "Module name: $MODULE_NAME"
else
    print_fail "go.mod not found"
    exit 1
fi

# Check dependencies
echo "🔗 Checking dependencies..."
if go mod verify > /dev/null 2>&1; then
    print_pass "Go modules are valid"
else
    print_warn "Running go mod tidy..."
    go mod tidy
fi

# Check if project builds
echo "🔨 Testing build..."
if go build -o /tmp/ordersystem ./cmd/ordersystem > /dev/null 2>&1; then
    print_pass "Application builds successfully"
    rm -f /tmp/ordersystem
else
    print_fail "Application failed to build"
    exit 1
fi

# Check tests
echo "🧪 Running tests..."
if go test ./internal/usecase/... > /dev/null 2>&1; then
    print_pass "All tests pass"
else
    print_fail "Some tests failed"
    exit 1
fi

# Check infrastructure files
echo "🏗️ Checking infrastructure setup..."
if [ -f "docker-compose.yaml" ]; then
    print_pass "docker-compose.yaml exists"
else
    print_fail "docker-compose.yaml not found"
    exit 1
fi

if [ -f "third-party/mysql/init.sql" ]; then
    print_pass "MySQL initialization script exists"
else
    print_fail "MySQL initialization script not found"
    exit 1
fi

if [ -f "third-party/setup.sh" ] && [ -x "third-party/setup.sh" ]; then
    print_pass "Setup script exists and is executable"
else
    print_fail "Setup script not found or not executable"
    exit 1
fi

# Check configuration files
echo "⚙️ Checking configuration..."
if [ -f "configs/config.go" ]; then
    print_pass "Configuration file exists"
else
    print_warn "Configuration file not found (this may be okay)"
fi

# Validate Docker Compose configuration
echo "🐳 Validating Docker Compose..."
if docker compose config > /dev/null 2>&1; then
    print_pass "Docker Compose configuration is valid"
else
    print_fail "Docker Compose configuration has errors"
    exit 1
fi

echo ""
echo "🎉 Setup verification completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Create a .env file with your environment variables"
echo "2. Start infrastructure: ./third-party/setup.sh"
echo "3. Run the application: go run cmd/ordersystem/main.go"
echo ""
echo "📊 Access points:"
echo "   - REST API: http://localhost:8000"
echo "   - gRPC: localhost:50051"
echo "   - GraphQL: http://localhost:8080"
echo "   - RabbitMQ: http://localhost:15672" 