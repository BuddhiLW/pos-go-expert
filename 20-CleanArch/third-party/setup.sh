#!/bin/bash

# Order System Infrastructure Setup Script

set -e

echo "🚀 Setting up Order System Infrastructure..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker compose is available (try both old and new syntax)
COMPOSE_CMD=""
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    print_error "docker compose is not available."
    exit 1
fi

print_status "Using compose command: $COMPOSE_CMD"

print_status "Stopping any existing containers..."
$COMPOSE_CMD down -v

print_status "Starting infrastructure services..."
$COMPOSE_CMD up -d

print_status "Waiting for services to be ready..."

# Wait for MySQL
print_status "Waiting for MySQL to be ready..."
until docker exec mysql mysqladmin ping -h localhost --silent 2>/dev/null; do
    echo -n "."
    sleep 2
done
echo ""

# Wait for RabbitMQ
print_status "Waiting for RabbitMQ to be ready..."
until docker exec rabbitmq rabbitmq-diagnostics ping > /dev/null 2>&1; do
    echo -n "."
    sleep 2
done
echo ""

# Wait a bit more for RabbitMQ management to load definitions
print_status "Waiting for RabbitMQ management and definitions to load..."
sleep 5

# Verify database setup
print_status "Verifying database setup..."
ORDERS_COUNT=$(docker exec mysql mysql -u root -proot -e "SELECT COUNT(*) FROM orders.orders;" -s -N 2>/dev/null || echo "0")

if [ "$ORDERS_COUNT" -gt 0 ]; then
    print_status "Database initialized successfully with $ORDERS_COUNT sample orders."
else
    print_warning "Database might not be properly initialized. Check the logs."
fi

# Verify RabbitMQ queue setup
print_status "Verifying RabbitMQ configuration..."

# Check if queues were created
QUEUES_OUTPUT=$(docker exec rabbitmq rabbitmqctl list_queues name 2>/dev/null || echo "")
if echo "$QUEUES_OUTPUT" | grep -q "orders.created"; then
    print_status "✅ Queue 'orders.created' is available"
else
    print_warning "⚠️ Queue 'orders.created' not found"
fi

if echo "$QUEUES_OUTPUT" | grep -q "orders.events"; then
    print_status "✅ Queue 'orders.events' is available"
else
    print_warning "⚠️ Queue 'orders.events' not found"
fi

# Check if exchanges were created
EXCHANGES_OUTPUT=$(docker exec rabbitmq rabbitmqctl list_exchanges name 2>/dev/null || echo "")
if echo "$EXCHANGES_OUTPUT" | grep -q "orders.direct"; then
    print_status "✅ Exchange 'orders.direct' is available"
else
    print_warning "⚠️ Exchange 'orders.direct' not found"
fi

if echo "$EXCHANGES_OUTPUT" | grep -q "orders.topic"; then
    print_status "✅ Exchange 'orders.topic' is available"
else
    print_warning "⚠️ Exchange 'orders.topic' not found"
fi

# Check bindings
BINDINGS_COUNT=$(docker exec rabbitmq rabbitmqctl list_bindings 2>/dev/null | wc -l || echo "0")
if [ "$BINDINGS_COUNT" -gt 3 ]; then
    print_status "✅ Queue bindings are configured"
else
    print_warning "⚠️ Some queue bindings might be missing"
fi

# Show service status
echo ""
print_info "🎯 Service Status:"
echo "  📁 MySQL:    localhost:3306 (user: root, password: root)"
echo "  🐰 RabbitMQ: localhost:5672 (AMQP)"
echo "  🌐 RabbitMQ: http://localhost:15672 (Management UI)"
echo ""

print_info "📊 Database Information:"
echo "  • Database: orders"
echo "  • Sample orders: $ORDERS_COUNT"
echo ""

print_info "📨 RabbitMQ Information:"
echo "  • Exchanges: amq.direct, orders.direct, orders.topic"
echo "  • Queues: orders.created, orders.events" 
echo "  • Login: guest/guest"
echo ""

print_status "Infrastructure setup complete! 🎉"
print_status "You can now start the application with: go run cmd/ordersystem/main.go"
echo ""

print_info "🧪 Test the message flow:"
echo "1. Start the application: go run cmd/ordersystem/main.go"
echo "2. Create an order: curl -X POST http://localhost:8000/order -H 'Content-Type: application/json' -d '{\"id\":\"test-order\",\"price\":100.0,\"tax\":10.0}'"
echo "3. Check RabbitMQ UI: http://localhost:15672"
echo ""

# Optional: Show logs
read -p "Do you want to see the service logs? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    $COMPOSE_CMD logs -f
fi 