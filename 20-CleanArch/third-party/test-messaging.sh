#!/bin/bash

# RabbitMQ Message Flow Test Script

set -e

echo "🧪 Testing RabbitMQ Message Flow..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if services are running
print_info "Checking if services are running..."

if ! docker ps | grep -q rabbitmq; then
    print_error "RabbitMQ container is not running"
    print_info "Run: ./third-party/setup.sh"
    exit 1
fi

if ! docker ps | grep -q mysql; then
    print_error "MySQL container is not running"
    print_info "Run: ./third-party/setup.sh"
    exit 1
fi

print_success "Services are running"

# Test 1: Check queues exist
print_info "Test 1: Checking if queues exist..."

QUEUES=$(docker exec rabbitmq rabbitmqctl list_queues name -q 2>/dev/null)

if echo "$QUEUES" | grep -q "orders.created"; then
    print_success "Queue 'orders.created' exists"
else
    print_error "Queue 'orders.created' not found"
    exit 1
fi

if echo "$QUEUES" | grep -q "orders.events"; then
    print_success "Queue 'orders.events' exists"
else
    print_error "Queue 'orders.events' not found"
    exit 1
fi

# Test 2: Check exchanges
print_info "Test 2: Checking exchanges..."

EXCHANGES=$(docker exec rabbitmq rabbitmqctl list_exchanges name -q 2>/dev/null)

if echo "$EXCHANGES" | grep -q "orders.direct"; then
    print_success "Exchange 'orders.direct' exists"
else
    print_error "Exchange 'orders.direct' not found"
    exit 1
fi

# Test 3: Check bindings
print_info "Test 3: Checking queue bindings..."

BINDINGS=$(docker exec rabbitmq rabbitmqctl list_bindings source_name destination_name -q 2>/dev/null)

if echo "$BINDINGS" | grep -q "amq.direct.*orders.created"; then
    print_success "Binding: amq.direct → orders.created"
else
    print_warning "Binding: amq.direct → orders.created not found"
fi

if echo "$BINDINGS" | grep -q "orders.direct.*orders.events"; then
    print_success "Binding: orders.direct → orders.events"
else
    print_warning "Binding: orders.direct → orders.events not found"
fi

# Test 4: Get current message counts
print_info "Test 4: Checking current message counts..."

ORDERS_CREATED_COUNT=$(docker exec rabbitmq rabbitmqctl list_queues name messages -q 2>/dev/null | grep orders.created | awk '{print $2}' || echo "0")
ORDERS_EVENTS_COUNT=$(docker exec rabbitmq rabbitmqctl list_queues name messages -q 2>/dev/null | grep orders.events | awk '{print $2}' || echo "0")

print_info "Current messages in orders.created: $ORDERS_CREATED_COUNT"
print_info "Current messages in orders.events: $ORDERS_EVENTS_COUNT"

# Test 5: Send a test message directly to RabbitMQ
print_info "Test 5: Publishing test message to amq.direct..."

# Create a temporary test message
TEST_MESSAGE='{"ID":"test-msg-001","Price":99.99,"Tax":9.99,"FinalPrice":109.98,"test":true}'

# Publish directly using rabbitmqctl
docker exec rabbitmq bash -c "
rabbitmqctl eval \"
rabbit_basic:publish(
    {resource, <<\\\"/\\\">>, exchange, <<\\\"amq.direct\\\">>},
    <<\\\"\\\">>,
    false,
    false,
    {amqp_msg, {P_basic, <<\\\"application/json\\\">>, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined}, <<\\\"$TEST_MESSAGE\\\">>}
).
\"" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    print_success "Test message published to amq.direct"
else
    print_warning "Failed to publish test message"
fi

# Wait a moment and check message count
sleep 2

NEW_ORDERS_CREATED_COUNT=$(docker exec rabbitmq rabbitmqctl list_queues name messages -q 2>/dev/null | grep orders.created | awk '{print $2}' || echo "0")

if [ "$NEW_ORDERS_CREATED_COUNT" -gt "$ORDERS_CREATED_COUNT" ]; then
    print_success "Message received in orders.created queue!"
    print_info "Messages increased from $ORDERS_CREATED_COUNT to $NEW_ORDERS_CREATED_COUNT"
else
    print_warning "No new messages detected in orders.created queue"
fi

# Test 6: Check RabbitMQ Management API
print_info "Test 6: Testing RabbitMQ Management API..."

# Test if management API is accessible
if curl -s -u guest:guest http://localhost:15672/api/overview > /dev/null 2>&1; then
    print_success "RabbitMQ Management API is accessible"
    
    # Get queue details via API
    QUEUE_INFO=$(curl -s -u guest:guest http://localhost:15672/api/queues/%2F/orders.created 2>/dev/null)
    if echo "$QUEUE_INFO" | grep -q "orders.created"; then
        print_success "orders.created queue details available via API"
    else
        print_warning "orders.created queue details not available via API"
    fi
else
    print_warning "RabbitMQ Management API not accessible"
fi

echo ""
print_success "RabbitMQ Message Flow Test Completed!"
echo ""
print_info "📋 Summary:"
echo "  • Queues: orders.created, orders.events"
echo "  • Exchanges: amq.direct, orders.direct, orders.topic"
echo "  • Management UI: http://localhost:15672 (guest/guest)"
echo ""
print_info "🚀 Ready to test with the application:"
echo "  1. Start app: go run cmd/ordersystem/main.go"
echo "  2. Create order: curl -X POST http://localhost:8000/order -H 'Content-Type: application/json' -d '{\"id\":\"test-order\",\"price\":100.0,\"tax\":10.0}'"
echo "  3. Monitor: http://localhost:15672" 