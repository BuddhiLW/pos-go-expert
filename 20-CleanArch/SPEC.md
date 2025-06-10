# Order System - Clean Architecture

A demonstration of Clean Architecture principles in Go with multiple API interfaces (REST, gRPC, GraphQL) and event-driven messaging.

## 🏗️ Architecture

This project implements Clean Architecture with the following layers:

- **Domain Layer**: Business entities and interfaces (`internal/entity/`)
- **Application Layer**: Use cases and business logic (`internal/usecase/`)
- **Infrastructure Layer**: External concerns (`internal/infra/`)
  - Database (MySQL)
  - Web servers (Chi, gRPC, GraphQL)
  - Message brokers (RabbitMQ)
- **Event Layer**: Event-driven messaging (`internal/event/`)

## 🚀 Features

✅ **Multiple API Interfaces**:
- REST API (Chi router)
- gRPC API
- GraphQL API with Playground

✅ **Event-Driven Architecture**:
- RabbitMQ message broker integration
- Order creation events
- Automatic queue and exchange setup

✅ **Clean Architecture**:
- SOLID principles implementation
- Dependency injection with Google Wire
- Repository pattern

✅ **Robust Infrastructure**:
- MySQL database with automated schema setup
- Docker Compose for local development
- Health checks and monitoring

## 📋 Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- Make (optional)

## 🚀 Quick Start

### 1. **Start Infrastructure**
```bash
# Option A: Use the automated setup script
./third-party/setup.sh

# Option B: Manual setup
docker compose up -d
```

### 2. **Run the Application**
```bash
go run cmd/ordersystem/main.go
```

### 3. **Access APIs**
- **REST API**: http://localhost:8000
- **gRPC**: localhost:50051
- **GraphQL Playground**: http://localhost:8080

### 4. **Test the APIs**
See the [API Documentation](api/README.md) for detailed endpoint information.

## 🗄️ Database Schema

The application uses MySQL with the following schema:

```sql
CREATE TABLE orders (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    price DECIMAL(10,2) NOT NULL,
    tax DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

Sample data is automatically loaded during initialization.

## 📨 Message Broker Architecture

RabbitMQ is configured with automatic setup for event-driven messaging:

### **Exchanges:**
- `amq.direct` - Built-in direct exchange (used by OrderCreatedHandler)
- `orders.direct` - Custom direct exchange for order routing
- `orders.topic` - Topic exchange for pattern-based routing

### **Queues:**
- `orders.created` - Receives order creation events
- `orders.events` - General order events processing

### **Message Flow:**
```
Order Creation → OrderCreatedHandler → amq.direct → orders.created queue
                                    ↓
                            OrderConsumer processes messages
                                    ↓  
                            Publishes to orders.direct → orders.events queue
```

### **Testing Message Flow:**
```bash
# Test RabbitMQ configuration
./third-party/test-messaging.sh

# Access RabbitMQ Management UI
open http://localhost:15672  # guest/guest
```

## 🔧 Configuration

Create a `.env` file with the following variables:

```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=:8000
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/usecase/... -v

# Test RabbitMQ messaging
./third-party/test-messaging.sh
```

## 📁 Project Structure

```
├── api/                    # API documentation and test files
├── cmd/ordersystem/        # Application entry point
├── configs/                # Configuration management
├── docs/                   # Project documentation
├── internal/
│   ├── entity/            # Domain entities
│   ├── usecase/           # Business use cases
│   ├── infra/
│   │   ├── database/      # Database implementations
│   │   ├── grpc/          # gRPC server and protobuf
│   │   ├── graph/         # GraphQL implementation
│   │   └── web/           # HTTP REST implementation
│   └── event/             # Event handling
│       └── handler/       # Message producers and consumers
├── pkg/events/            # Shared event utilities
├── third-party/           # Infrastructure setup
│   ├── mysql/            # Database initialization
│   ├── rabbitmq/         # Message broker config
│   ├── setup.sh          # Automated setup script
│   └── test-messaging.sh # RabbitMQ testing script
└── docker-compose.yaml   # Container orchestration
```

## 🔍 API Examples

### REST API
```bash
# Create an order
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"order-123","price":100.0,"tax":10.0}'

# List orders
curl http://localhost:8000/orders
```

### GraphQL
```graphql
# Create order mutation
mutation {
  createOrder(input: {
    id: "order-123"
    Price: 100.0
    Tax: 10.0
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}

# Query orders
query {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

### gRPC

``` bash
evans internal/infra/grpc/protofiles/order.proto
pb.OrderService@127.0.0.1:50051> call ListOrders
```

### Message Flow Testing
```bash
# Create an order and monitor messages
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"test-order","price":100.0,"tax":10.0}'

# Check message in RabbitMQ UI
open http://localhost:15672
# Navigate to Queues → orders.created to see the message
```

## 🐛 Troubleshooting

### Database Issues
```bash
# Check MySQL logs
docker compose logs mysql

# Connect to database
docker exec -it mysql mysql -u root -p orders

# Reset database
docker compose down -v && docker compose up -d
```

### RabbitMQ Issues
```bash
# Check RabbitMQ status
docker compose logs rabbitmq

# Test message configuration
./third-party/test-messaging.sh

# Access management UI
open http://localhost:15672

# Check queues via CLI
docker exec rabbitmq rabbitmqctl list_queues
docker exec rabbitmq rabbitmqctl list_exchanges
docker exec rabbitmq rabbitmqctl list_bindings
```

### Application Issues
```bash
# Check Go version
go version

# Verify dependencies
go mod tidy
go mod verify

# Run with verbose logging
go run cmd/ordersystem/main.go -v
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## 📜 License

This project is part of a Go Clean Architecture tutorial and is provided for educational purposes.

## 🔗 Related Links

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
- [GraphQL with gqlgen](https://gqlgen.com/getting-started/)
- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html) 