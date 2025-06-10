# 🏗️ Order System - Project Context

## 📋 Executive Summary

This project implements a **Clean Architecture** order management system in Go, demonstrating the separation of concerns through well-defined layers and interfaces. The system currently supports **order creation** through multiple interfaces (HTTP, gRPC, GraphQL) with event-driven architecture and message queuing.

---

## 🎯 Current State

**Status**: ✅ **MILESTONE COMPLETED** - "Get All Orders" feature fully implemented across all interfaces

### 🎯 Recently Completed Features

#### ✅ Order Creation (Complete)
- **Domain Layer**: Order entity with business logic validation
- **Application Layer**: CreateOrderUseCase with event dispatching
- **Infrastructure Layer**: 
  - MySQL repository with prepared statements
  - HTTP REST endpoint (`POST /order`)
  - gRPC service (`CreateOrder`)
  - GraphQL mutation (`createOrder`)
- **Event-Driven**: OrderCreated events published to RabbitMQ

#### ✅ Order Listing (Complete) - **NEW**
- **Domain Layer**: Extended OrderRepositoryInterface with `GetAll()` method
- **Application Layer**: ListOrdersUseCase with comprehensive DTO mapping
- **Infrastructure Layer**:
  - MySQL repository `GetAll()` implementation with proper error handling
  - HTTP REST endpoint (`GET /orders`)
  - gRPC service (`ListOrders`)
  - GraphQL query (`orders`)
- **Testing**: Comprehensive test coverage with table-driven tests and mocks

### ✅ **Implemented Features**
- **Multi-Interface Support**: HTTP REST, gRPC, and GraphQL endpoints
- **Event-Driven Architecture**: RabbitMQ integration with order creation events
- **Clean Architecture**: Proper separation between entities, use cases, and infrastructure
- **Database Persistence**: MySQL integration with prepared statements
- **Dependency Injection**: Wire-based dependency management

### 🚧 **Next Milestone**
- **Get All Orders**: Implement method to retrieve all orders from the system
  - Add `GetAll() ([]*Order, error)` to `OrderRepositoryInterface`
  - Implement corresponding use case `ListOrdersUseCase`
  - Add endpoints across all interfaces (HTTP, gRPC, GraphQL)
  - Maintain Clean Architecture principles

---

## 🏛️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        PRESENTATION LAYER                        │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   HTTP REST     │      gRPC       │         GraphQL             │
│                 │                 │                             │
│ POST /order     │ CreateOrder()   │ mutation createOrder        │
│ GET  /orders    │ ListOrders()    │ query orders               │
└─────────────────┴─────────────────┴─────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────┐
│                      APPLICATION LAYER                          │
├─────────────────────────────────────────────────────────────────┤
│ • CreateOrderUseCase (with event dispatching)                  │
│ • ListOrdersUseCase (with DTO conversion)                      │
│ • OrderInputDTO / OrderOutputDTO / ListOrdersOutputDTO         │
└─────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────┐
│                        DOMAIN LAYER                             │
├─────────────────────────────────────────────────────────────────┤
│ • Order Entity (ID, Price, Tax, FinalPrice + business logic)    │
│ • OrderRepositoryInterface (Save, GetAll)                      │
│ • Domain Events (OrderCreated)                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER                         │
├─────────────────────────────────────────────────────────────────┤
│ • MySQL OrderRepository (prepared statements, connection pool)  │
│ • RabbitMQ Event Publisher                                      │
│ • HTTP Handlers (Chi router)                                    │
│ • gRPC Services (Protocol Buffers)                             │
│ • GraphQL Resolvers (gqlgen)                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📂 Project Structure

```
20-CleanArch/
├── api/                      # API testing files
│   └── create_order.http
├── cmd/ordersystem/          # Application entry point
│   ├── main.go              # Main application setup
│   ├── wire.go              # Dependency injection definitions
│   └── wire_gen.go          # Generated DI code
├── configs/                  # Configuration management
│   └── config.go            # Environment-based configuration
├── docs/                     # 📝 Documentation (NEW)
├── internal/                 # Private application code
│   ├── entity/              # 🏛️ Domain layer
│   │   ├── interface.go     # Repository interfaces
│   │   ├── order.go         # Order domain entity
│   │   └── order_test.go    # Domain tests
│   ├── usecase/             # 🎯 Application layer
│   │   └── create_order.go  # Order creation business logic
│   ├── infra/               # 🔧 Infrastructure layer
│   │   ├── database/        # Data persistence
│   │   │   ├── order_repository.go
│   │   │   └── order_repository_test.go
│   │   ├── web/             # HTTP interface
│   │   │   ├── order_handler.go
│   │   │   └── webserver/
│   │   ├── grpc/            # gRPC interface
│   │   │   ├── pb/          # Generated protobuf code
│   │   │   ├── protofiles/  # Proto definitions
│   │   │   └── service/     # gRPC service implementation
│   │   └── graph/           # GraphQL interface
│   │       ├── generated.go # Generated GraphQL code
│   │       ├── resolver.go  # GraphQL resolvers
│   │       ├── schema.graphqls
│   │       └── model/       # GraphQL models
│   └── event/               # Event handling
│       └── handler/         # Event handlers
├── pkg/                     # Public packages
│   └── events/              # Event system implementation
├── docker-compose.yaml      # Local development environment
├── gqlgen.yml              # GraphQL generation config
└── tools.go               # Build tools dependencies
```

---

## 🔧 Technology Stack

### **Core Framework**
- **Language**: Go 1.22+
- **Architecture**: Clean Architecture pattern
- **Dependency Injection**: Google Wire

### **Interfaces**
- **HTTP**: Native `net/http` with custom web server
- **gRPC**: `google.golang.org/grpc` with Protocol Buffers
- **GraphQL**: `github.com/99designs/gqlgen`

### **Infrastructure**
- **Database**: MySQL 5.7 with `database/sql`
- **Message Queue**: RabbitMQ with `github.com/streadway/amqp`
- **Configuration**: `github.com/spf13/viper`

### **Development**
- **Testing**: Go's built-in testing + table-driven tests
- **Containerization**: Docker Compose for local development
- **Code Generation**: Wire (DI), gqlgen (GraphQL), protoc (gRPC)

---

## 🎭 SOLID Principles Implementation

### **Single Responsibility Principle (SRP)**
- Each layer has a single concern
- Handlers only handle protocol translation
- Use cases contain only business logic
- Repositories only handle data persistence

### **Open/Closed Principle (OCP)**
- New interfaces can be added without modifying existing code
- Use cases are extensible through dependency injection

### **Liskov Substitution Principle (LSP)**
- Repository implementations are substitutable
- Event handlers can be swapped without breaking the system

### **Interface Segregation Principle (ISP)**
- `OrderRepositoryInterface` contains only necessary methods
- Each interface serves specific client needs

### **Dependency Inversion Principle (DIP)**
- High-level modules (use cases) don't depend on low-level modules
- Dependencies flow inward toward the domain layer

---

## 🗄️ Data Model

### **Order Entity**
```go
type Order struct {
    ID         string   // Unique identifier
    Price      float64  // Base price
    Tax        float64  // Tax amount
    FinalPrice float64  // Calculated final price
}
```

### **Business Rules**
- Order ID must not be empty
- Price must be greater than 0
- Tax must be greater than 0
- Final price is automatically calculated as Price + Tax

---

## 🔄 Current Data Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│  Interface  │───▶│  Use Case   │───▶│ Repository  │
│             │    │  (HTTP/     │    │ (Business   │    │ (Database   │
│             │    │ gRPC/GQL)   │    │  Logic)     │    │ Persistence)│
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                            │                   │
                            ▼                   ▼
                   ┌─────────────┐    ┌─────────────┐
                   │   DTO       │    │   Event     │
                   │ Conversion  │    │ Dispatch    │
                   │             │    │ (RabbitMQ)  │
                   └─────────────┘    └─────────────┘
```

---

## 🚀 Development Workflow

### **Adding New Features**
1. **Domain First**: Define entities and interfaces in `internal/entity/`
2. **Use Case**: Implement business logic in `internal/usecase/`
3. **Infrastructure**: Add implementations in `internal/infra/`
4. **Integration**: Wire dependencies in `cmd/ordersystem/wire.go`
5. **Testing**: Add comprehensive tests at each layer

### **Running the System**
```bash
# Start infrastructure
docker-compose up -d

# Run the application
cd cmd/ordersystem
go run .
```

### **Available Endpoints**
- **HTTP**: `POST http://localhost:8000/order`
- **gRPC**: `localhost:50051` (OrderService.CreateOrder)
- **GraphQL**: `http://localhost:8080` (Playground available)

---

## 🎯 Next Implementation Steps

### **Immediate: Get All Orders**
1. **Update Interface** (`internal/entity/interface.go`):
   ```go
   type OrderRepositoryInterface interface {
       Save(order *Order) error
       GetAll() ([]*Order, error)  // 🆕 Add this method
   }
   ```

2. **Create Use Case** (`internal/usecase/list_orders.go`):
   - Input: None or pagination parameters
   - Output: List of OrderOutputDTO
   - Business logic: Retrieve and format orders

3. **Implement Infrastructure**:
   - **Database**: Add `GetAll()` method to `OrderRepository`
   - **HTTP**: Add `GET /orders` endpoint
   - **gRPC**: Add `ListOrders` RPC method
   - **GraphQL**: Add `orders` query

4. **Wire Dependencies**: Update dependency injection

### **Future Enhancements**
- Order searching and filtering
- Order status management
- User authentication and authorization
- Order validation rules
- Caching layer
- Metrics and monitoring

---

## 🔍 Quality Metrics

### **Test Coverage**
- Target: ≥80% coverage on business logic
- Current: Domain layer fully tested
- Focus: Use case and integration testing

### **Code Quality**
- Clean Architecture compliance
- SOLID principles adherence
- Comprehensive error handling
- Proper dependency management

---

## 📚 References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design)

---

**🎉 Current Status**: The "Get All Orders" feature has been successfully implemented across all interfaces (HTTP REST, gRPC, GraphQL) following Clean Architecture principles and SOLID design patterns. The system now supports both order creation and retrieval with comprehensive test coverage and proper error handling. 