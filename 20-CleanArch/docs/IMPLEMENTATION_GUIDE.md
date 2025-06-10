# 🛠️ Implementation Guide: Get All Orders Feature

## 🎯 Overview

This guide demonstrates how to implement the "Get All Orders" feature following Clean Architecture principles. We'll work from the **inside-out**, starting with the domain layer and progressing to the infrastructure layer.

---

## 📋 Implementation Steps

### **Step 1: Update Domain Interface**

**File**: `internal/entity/interface.go`

**Goal**: Add the new method to our repository contract.

```go
package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetAll() ([]*Order, error)  // 🆕 NEW METHOD
}
```

**Why This First?**
- The domain layer drives the architecture
- Interface changes propagate outward
- Ensures all implementations must support this operation

---

### **Step 2: Create List Orders Use Case**

**File**: `internal/usecase/list_orders.go`

**Goal**: Implement business logic for retrieving orders.

```go
package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrdersOutputDTO struct {
	Orders []OrderOutputDTO `json:"orders"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() (ListOrdersOutputDTO, error) {
	orders, err := c.OrderRepository.GetAll()
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}

	var orderDTOs []OrderOutputDTO
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		orderDTOs = append(orderDTOs, dto)
	}

	return ListOrdersOutputDTO{
		Orders: orderDTOs,
	}, nil
}
```

**Key Points**:
- Uses the same `OrderOutputDTO` for consistency
- Wraps results in a `ListOrdersOutputDTO` for future extensibility
- No business logic complexity (just data retrieval and conversion)
- Follows SRP: Single responsibility of listing orders

---

### **Step 3: Implement Database Repository**

**File**: `internal/infra/database/order_repository.go`

**Goal**: Add actual database implementation.

```go
// Add this method to the existing OrderRepository struct

func (r *OrderRepository) GetAll() ([]*entity.Order, error) {
	stmt, err := r.Db.Prepare("SELECT id, price, tax, final_price FROM orders ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
```

**Key Points**:
- Uses prepared statements for security
- Proper resource cleanup with `defer`
- Returns domain entities, not DTOs
- Orders by ID for consistent results

---

### **Step 4: Add HTTP Handler**

**File**: `internal/infra/web/order_handler.go`

**Goal**: Add REST endpoint for listing orders.

```go
// Add this method to the existing WebOrderHandler struct

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	listOrders := usecase.NewListOrdersUseCase(h.OrderRepository)
	output, err := listOrders.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
```

**Usage**: `GET /orders`

**Key Points**:
- Creates use case instance (could be dependency injected)
- Handles HTTP-specific concerns (headers, status codes)
- Separates protocol concerns from business logic

---

### **Step 5: Add gRPC Service Method**

**File**: `internal/infra/grpc/protofiles/order.proto`

**Goal**: Define protobuf contract.

```protobuf
syntax = "proto3";
package order;
option go_package = "github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb";

// ... existing messages ...

message ListOrdersRequest {
    // Empty for now, future: pagination parameters
}

message ListOrdersResponse {
    repeated CreateOrderResponse orders = 1;
}

service OrderService {
    rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
    rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);  // 🆕 NEW
}
```

**File**: `internal/infra/grpc/service/order_service.go`

**Goal**: Implement gRPC service method.

```go
// Add this method to the existing OrderService struct

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	listOrders := usecase.NewListOrdersUseCase(s.OrderRepository)
	output, err := listOrders.Execute()
	if err != nil {
		return nil, err
	}

	var grpcOrders []*pb.CreateOrderResponse
	for _, order := range output.Orders {
		grpcOrder := &pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		grpcOrders = append(grpcOrders, grpcOrder)
	}

	return &pb.ListOrdersResponse{
		Orders: grpcOrders,
	}, nil
}
```

**Key Points**:
- Reuses existing `CreateOrderResponse` message type
- Converts DTOs to protobuf messages
- Handles gRPC-specific type conversions (float64 → float32)

---

### **Step 6: Add GraphQL Query**

**File**: `internal/infra/graph/schema.graphqls`

**Goal**: Define GraphQL schema.

```graphql
# ... existing types ...

type Query {
    orders: [Order!]!  # 🆕 NEW QUERY
}

type Mutation {
    createOrder(input: OrderInput): Order
}
```

**File**: `internal/infra/graph/schema.resolvers.go`

**Goal**: Implement GraphQL resolver.

```go
// Add this method - gqlgen will generate the interface

func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	listOrders := usecase.NewListOrdersUseCase(r.OrderRepository)
	output, err := listOrders.Execute()
	if err != nil {
		return nil, err
	}

	var gqlOrders []*model.Order
	for _, order := range output.Orders {
		gqlOrder := &model.Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		gqlOrders = append(gqlOrders, gqlOrder)
	}

	return gqlOrders, nil
}

// Add this method to complete the resolver interface
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
```

**Key Points**:
- GraphQL schema-first approach
- Converts DTOs to GraphQL models
- Maintains type safety through generated interfaces

---

### **Step 7: Update Dependency Injection**

**File**: `cmd/ordersystem/wire.go`

**Goal**: Add dependencies for the new use case.

```go
// Add this function
func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		database.NewOrderRepository,
		usecase.NewListOrdersUseCase,
		wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
	)
	return &usecase.ListOrdersUseCase{}
}
```

**File**: `cmd/ordersystem/main.go`

**Goal**: Wire up the new endpoints.

```go
func main() {
	// ... existing setup ...

	// Add list orders use case
	listOrdersUseCase := NewListOrdersUseCase(db)

	// Update web server
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.List)  // 🆕 NEW

	// Update gRPC service
	createOrderService := service.NewOrderService(*createOrderUseCase)
	// Note: You'll need to update the service constructor to accept listOrdersUseCase

	// Update GraphQL resolver
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CreateOrderUseCase: *createOrderUseCase,
			// Add: ListOrdersUseCase: *listOrdersUseCase,
		},
	}))

	// ... rest of main function
}
```

---

### **Step 8: Add Comprehensive Tests**

#### **Domain Layer Test**

**File**: `internal/entity/order_test.go`

```go
// Add test for multiple orders if needed
func TestOrder_MultipleOrders(t *testing.T) {
	// Test any domain logic related to collections
	// Currently not needed as GetAll is just data retrieval
}
```

#### **Use Case Test**

**File**: `internal/usecase/list_orders_test.go`

```go
package usecase

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order *entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetAll() ([]*entity.Order, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func TestListOrdersUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		mockOrders     []*entity.Order
		mockError      error
		expectedOrders int
		expectedError  bool
	}{
		{
			name: "successful retrieval",
			mockOrders: []*entity.Order{
				{ID: "1", Price: 100, Tax: 10, FinalPrice: 110},
				{ID: "2", Price: 200, Tax: 20, FinalPrice: 220},
			},
			mockError:      nil,
			expectedOrders: 2,
			expectedError:  false,
		},
		{
			name:           "empty list",
			mockOrders:     []*entity.Order{},
			mockError:      nil,
			expectedOrders: 0,
			expectedError:  false,
		},
		{
			name:           "repository error",
			mockOrders:     nil,
			mockError:      errors.New("database error"),
			expectedOrders: 0,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockOrderRepository)
			mockRepo.On("GetAll").Return(tt.mockOrders, tt.mockError)
			
			useCase := NewListOrdersUseCase(mockRepo)

			// Act
			result, err := useCase.Execute()

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result.Orders, tt.expectedOrders)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
```

#### **Integration Test**

**File**: `internal/infra/database/order_repository_test.go`

```go
// Add this test to the existing file

func TestOrderRepository_GetAll(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Setup test table
	_, err = db.Exec(`CREATE TABLE orders (
		id VARCHAR(255) PRIMARY KEY,
		price DECIMAL(10,2),
		tax DECIMAL(10,2),
		final_price DECIMAL(10,2)
	)`)
	assert.NoError(t, err)

	// Insert test data
	_, err = db.Exec("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)", "1", 100.0, 10.0, 110.0)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)", "2", 200.0, 20.0, 220.0)
	assert.NoError(t, err)

	// Test
	repo := NewOrderRepository(db)
	orders, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, "1", orders[0].ID)
	assert.Equal(t, "2", orders[1].ID)
}
```

---

## 🔄 Testing the Implementation

### **HTTP REST API**
```bash
# Get all orders
curl -X GET http://localhost:8000/orders
```

**Expected Response**:
```json
{
  "orders": [
    {
      "id": "1",
      "price": 100.0,
      "tax": 10.0,
      "final_price": 110.0
    },
    {
      "id": "2", 
      "price": 200.0,
      "tax": 20.0,
      "final_price": 220.0
    }
  ]
}
```

### **gRPC**
```bash
# Using grpcurl
grpcurl -plaintext localhost:50051 order.OrderService/ListOrders
```

### **GraphQL**
```graphql
query {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

---

## 🎯 Key Architectural Benefits

### **1. Separation of Concerns**
- **Domain**: Business rules isolated from infrastructure
- **Use Case**: Application logic independent of delivery mechanism
- **Infrastructure**: Technology-specific implementations

### **2. Testability**
- Each layer can be tested independently
- Mock repositories for use case testing
- Integration tests for database layer

### **3. Flexibility**
- Easy to add new interfaces (WebSocket, CLI, etc.)
- Database can be changed without affecting business logic
- Use cases are reusable across interfaces

### **4. SOLID Compliance**
- **SRP**: Each class has single responsibility
- **OCP**: Open for extension (new interfaces)
- **LSP**: Repository implementations are substitutable
- **ISP**: Interface segregation through focused contracts
- **DIP**: Dependencies flow inward to domain

---

## 🔮 Future Enhancements

### **Pagination Support**
```go
type ListOrdersInputDTO struct {
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
}

type ListOrdersOutputDTO struct {
    Orders     []OrderOutputDTO `json:"orders"`
    TotalCount int              `json:"total_count"`
    Page       int              `json:"page"`
    PageSize   int              `json:"page_size"`
}
```

### **Filtering and Sorting**
```go
type ListOrdersInputDTO struct {
    MinPrice   *float64 `json:"min_price,omitempty"`
    MaxPrice   *float64 `json:"max_price,omitempty"`
    SortBy     string   `json:"sort_by"`     // "price", "tax", "final_price"
    SortOrder  string   `json:"sort_order"`  // "asc", "desc"
}
```

### **Caching Layer**
```go
type CachedOrderRepository struct {
    repo  entity.OrderRepositoryInterface
    cache Cache
}

func (c *CachedOrderRepository) GetAll() ([]*entity.Order, error) {
    // Check cache first
    if cached := c.cache.Get("all_orders"); cached != nil {
        return cached.([]*entity.Order), nil
    }
    
    // Fallback to repository
    orders, err := c.repo.GetAll()
    if err == nil {
        c.cache.Set("all_orders", orders, 5*time.Minute)
    }
    return orders, err
}
```

---

## ✅ Implementation Checklist

- [ ] **Domain**: Update `OrderRepositoryInterface`
- [ ] **Use Case**: Create `ListOrdersUseCase`  
- [ ] **Database**: Implement `GetAll()` in repository
- [ ] **HTTP**: Add `List` handler method
- [ ] **gRPC**: Update proto and service
- [ ] **GraphQL**: Add query to schema and resolver
- [ ] **Dependency Injection**: Wire new dependencies
- [ ] **Testing**: Add unit and integration tests
- [ ] **Documentation**: Update API docs
- [ ] **Validation**: Test all endpoints manually

---

**🎯 Result**: A fully functional "Get All Orders" feature that maintains Clean Architecture principles, follows SOLID design patterns, and provides multiple interface options while remaining highly testable and maintainable.** 