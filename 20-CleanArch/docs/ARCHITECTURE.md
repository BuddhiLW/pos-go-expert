# 🏗️ Order System - Architecture Documentation

## 🎯 Overview

This document describes the architectural patterns, design decisions, and implementation details of the Order System, which follows **Clean Architecture** principles to ensure maintainability, testability, and separation of concerns.

---

## 🏛️ Clean Architecture Layers

### **1. Domain Layer (Entities)**
**Location**: `internal/entity/`

**Purpose**: Contains enterprise business rules and domain entities.

**Components**:
- `Order` entity with business logic
- `OrderRepositoryInterface` defining data access contracts
- Domain-specific validation rules

**Key Principles**:
- No dependencies on external frameworks
- Pure business logic and entities
- Stable and independent of external changes

```go
// Example: Order entity with business rules
type Order struct {
    ID         string
    Price      float64
    Tax        float64
    FinalPrice float64
}

func (o *Order) CalculateFinalPrice() error {
    o.FinalPrice = o.Price + o.Tax
    return o.IsValid()
}
```

### **2. Application Layer (Use Cases)**
**Location**: `internal/usecase/`

**Purpose**: Contains application-specific business rules and orchestrates domain entities.

**Components**:
- `CreateOrderUseCase` - Order creation workflow
- Input/Output DTOs for data transfer
- Event dispatching logic

**Key Principles**:
- Orchestrates domain entities
- Depends only on domain layer
- Contains application-specific business rules

```go
// Example: Use case with dependency injection
type CreateOrderUseCase struct {
    OrderRepository entity.OrderRepositoryInterface
    OrderCreated    events.EventInterface
    EventDispatcher events.EventDispatcherInterface
}
```

### **3. Interface Adapters Layer**
**Location**: `internal/infra/`

**Purpose**: Converts data between use cases and external interfaces.

**Components**:
- **Web** (`web/`): HTTP handlers and routing
- **gRPC** (`grpc/`): Protocol buffer services
- **GraphQL** (`graph/`): Schema and resolvers
- **Database** (`database/`): Repository implementations

**Key Principles**:
- Adapts external interfaces to internal use cases
- Handles protocol-specific concerns
- Converts between DTOs and domain entities

### **4. Frameworks & Drivers Layer**
**Location**: `cmd/`, `configs/`, `pkg/`

**Purpose**: Contains frameworks, databases, web servers, and external tools.

**Components**:
- Application entry point (`cmd/ordersystem/main.go`)
- Configuration management (`configs/`)
- External event system (`pkg/events/`)
- Docker and infrastructure setup

---

## 🔄 Data Flow Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          REQUEST FLOW                               │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────┐  1. HTTP/gRPC/GraphQL   ┌─────────────────────────────┐
│   Client    │─────────────────────────▶│     Interface Layer        │
│             │                         │  (Handlers/Services)        │
└─────────────┘                         └─────────────────────────────┘
                                                    │
                                          2. Convert to DTO
                                                    ▼
                                        ┌─────────────────────────────┐
                                        │      Use Case Layer         │
                                        │   (Business Logic)          │
                                        └─────────────────────────────┘
                                                    │
                                          3. Domain Operations
                                                    ▼
                                        ┌─────────────────────────────┐
                                        │       Domain Layer          │
                                        │   (Entities & Rules)        │
                                        └─────────────────────────────┘
                                                    │
                                          4. Repository Interface
                                                    ▼
                                        ┌─────────────────────────────┐
                                        │   Infrastructure Layer      │
                                        │  (Database/Events)          │
                                        └─────────────────────────────┘
```

---

## 🎭 Interface Implementations

### **HTTP REST API**
**Implementation**: `internal/infra/web/order_handler.go`

**Features**:
- JSON input/output
- HTTP status code handling
- Error response formatting

**Endpoints**:
- `POST /order` - Create new order

```go
func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
    var dto usecase.OrderInputDTO
    // 1. Parse HTTP request
    json.NewDecoder(r.Body).Decode(&dto)
    
    // 2. Execute use case
    createOrder := usecase.NewCreateOrderUseCase(/*deps*/)
    output, err := createOrder.Execute(dto)
    
    // 3. Return HTTP response
    json.NewEncoder(w).Encode(output)
}
```

### **gRPC Service**
**Implementation**: `internal/infra/grpc/service/order_service.go`

**Features**:
- Protocol buffer definitions
- Type-safe RPC calls
- Streaming support (future)

**Methods**:
- `CreateOrder(CreateOrderRequest) → CreateOrderResponse`

```go
func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    // 1. Convert protobuf to DTO
    dto := usecase.OrderInputDTO{
        ID:    in.Id,
        Price: float64(in.Price),
        Tax:   float64(in.Tax),
    }
    
    // 2. Execute use case
    output, err := s.CreateOrderUseCase.Execute(dto)
    
    // 3. Convert DTO to protobuf
    return &pb.CreateOrderResponse{/*...*/}, nil
}
```

### **GraphQL API**
**Implementation**: `internal/infra/graph/schema.resolvers.go`

**Features**:
- Schema-first approach
- Type-safe resolvers
- Introspection support

**Operations**:
- `Mutation.createOrder(input: OrderInput) → Order`

```go
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
    // 1. Convert GraphQL input to DTO
    dto := usecase.OrderInputDTO{
        ID:    input.ID,
        Price: float64(input.Price),
        Tax:   float64(input.Tax),
    }
    
    // 2. Execute use case
    output, err := r.CreateOrderUseCase.Execute(dto)
    
    // 3. Convert DTO to GraphQL model
    return &model.Order{/*...*/}, nil
}
```

---

## 🗄️ Data Persistence Strategy

### **Repository Pattern**
**Interface**: `internal/entity/interface.go`
**Implementation**: `internal/infra/database/order_repository.go`

**Benefits**:
- Abstraction over data access
- Testability through mocking
- Database-agnostic business logic

```go
// Domain interface (stable)
type OrderRepositoryInterface interface {
    Save(order *Order) error
    GetAll() ([]*Order, error)  // Future implementation
}

// Infrastructure implementation (replaceable)
type OrderRepository struct {
    Db *sql.DB
}

func (r *OrderRepository) Save(order *entity.Order) error {
    stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
    // ... implementation details
}
```

### **Database Schema**
```sql
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    price DECIMAL(10,2) NOT NULL,
    tax DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🔥 Event-Driven Architecture

### **Event System**
**Location**: `pkg/events/`

**Components**:
- `EventInterface` - Event contract
- `EventDispatcher` - Event routing
- `OrderCreatedHandler` - RabbitMQ publisher

**Flow**:
```go
// 1. Use case creates event
c.OrderCreated.SetPayload(dto)
c.EventDispatcher.Dispatch(c.OrderCreated)

// 2. Handler processes event
type OrderCreatedHandler struct {
    RabbitMQChannel *amqp.Channel
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface) {
    // Publish to RabbitMQ
}
```

### **Benefits**:
- Loose coupling between components
- Asynchronous processing capability
- Scalability through message queuing
- Audit trail and event sourcing potential

---

## 🏗️ Dependency Injection

### **Google Wire**
**Configuration**: `cmd/ordersystem/wire.go`

**Benefits**:
- Compile-time dependency resolution
- No runtime reflection
- Clear dependency graphs

```go
//go:build wireinject
// +build wireinject

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
    wire.Build(
        database.NewOrderRepository,
        entity.NewOrderCreated,
        usecase.NewCreateOrderUseCase,
        wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
    )
    return &usecase.CreateOrderUseCase{}
}
```

### **Dependency Graph**:
```
main.go
├── CreateOrderUseCase
│   ├── OrderRepository (interface)
│   │   └── MySQL Implementation
│   ├── OrderCreated (event)
│   └── EventDispatcher
│       └── RabbitMQ Handler
├── WebOrderHandler
├── gRPC OrderService
└── GraphQL Resolver
```

---

## 🧪 Testing Strategy

### **Layer-Specific Testing**

#### **Domain Layer Testing**
```go
func TestOrder_CalculateFinalPrice(t *testing.T) {
    tests := []struct {
        name  string
        price float64
        tax   float64
        want  float64
    }{
        {"valid calculation", 100.0, 10.0, 110.0},
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            order := &Order{Price: tt.price, Tax: tt.tax}
            order.CalculateFinalPrice()
            assert.Equal(t, tt.want, order.FinalPrice)
        })
    }
}
```

#### **Use Case Testing**
```go
func TestCreateOrderUseCase_Execute(t *testing.T) {
    // Arrange
    mockRepo := &MockOrderRepository{}
    mockDispatcher := &MockEventDispatcher{}
    useCase := NewCreateOrderUseCase(mockRepo, mockEvent, mockDispatcher)
    
    // Act
    result, err := useCase.Execute(validInput)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedOutput, result)
    mockRepo.AssertExpectations(t)
}
```

#### **Integration Testing**
```go
func TestOrderRepository_Save(t *testing.T) {
    // Use test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewOrderRepository(db)
    order := &entity.Order{/*...*/}
    
    err := repo.Save(order)
    assert.NoError(t, err)
    
    // Verify persistence
    var count int
    db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", order.ID).Scan(&count)
    assert.Equal(t, 1, count)
}
```

---

## 🚀 Deployment Architecture

### **Local Development**
```yaml
# docker-compose.yaml
services:
  mysql:
    image: mysql:5.7
    environment:
      MYSQL_DATABASE: orders
    ports:
      - "3306:3306"
  
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
```

### **Production Considerations**
- **Containerization**: Multi-stage Docker builds
- **Service Mesh**: Kubernetes with Istio
- **Observability**: Prometheus metrics, Jaeger tracing
- **Configuration**: Environment-based configuration
- **Security**: TLS termination, API authentication

---

## 🔧 Configuration Management

### **Environment-Based Configuration**
```go
type Config struct {
    DBDriver        string `mapstructure:"DB_DRIVER"`
    DBHost          string `mapstructure:"DB_HOST"`
    DBPort          string `mapstructure:"DB_PORT"`
    DBUser          string `mapstructure:"DB_USER"`
    DBPassword      string `mapstructure:"DB_PASSWORD"`
    DBName          string `mapstructure:"DB_NAME"`
    WebServerPort   string `mapstructure:"WEB_SERVER_PORT"`
    GRPCServerPort  string `mapstructure:"GRPC_SERVER_PORT"`
    GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}
```

### **Configuration Loading**
```go
func LoadConfig(path string) (*Config, error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")
    viper.AutomaticEnv()
    
    err := viper.ReadInConfig()
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = viper.Unmarshal(&config)
    return &config, err
}
```

---

## 📊 Performance Considerations

### **Database Optimization**
- Connection pooling for MySQL
- Prepared statements for security and performance
- Proper indexing strategy
- Query optimization

### **Concurrency**
- Go routines for parallel processing
- Context-based cancellation
- Resource pooling

### **Caching Strategy** (Future)
- Redis for frequently accessed data
- Application-level caching
- Cache invalidation strategies

---

## 🔒 Security Considerations

### **Data Validation**
- Input validation at interface boundaries
- Domain-level business rule validation
- SQL injection prevention through prepared statements

### **Error Handling**
- Proper error propagation
- No sensitive information in error messages
- Structured logging for security monitoring

### **Authentication & Authorization** (Future)
- JWT token validation
- Role-based access control
- Rate limiting

---

## 🎯 Next Steps for Implementation

### **GetAll Orders Feature**
Following the same architectural patterns:

1. **Domain**: Add `GetAll() ([]*Order, error)` to interface
2. **Use Case**: Create `ListOrdersUseCase` with pagination
3. **Infrastructure**: Implement in all adapters (HTTP, gRPC, GraphQL)
4. **Testing**: Add comprehensive tests at each layer
5. **Documentation**: Update API documentation

### **Future Enhancements**
- Order filtering and searching
- Pagination and sorting
- Real-time updates via WebSocket
- Distributed caching
- Microservice decomposition

---

**This architecture ensures maintainability, testability, and scalability while adhering to Clean Architecture principles and SOLID design patterns.** 