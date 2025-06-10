# Order System API Documentation

## 🚀 Available Endpoints

### HTTP REST API

#### Create Order
```http
POST /order
Content-Type: application/json

{
  "id": "order-123",
  "price": 100.0,
  "tax": 10.0
}
```

**Response:**
```json
{
  "id": "order-123",
  "price": 100.0,
  "tax": 10.0,
  "final_price": 110.0
}
```

#### List All Orders
```http
GET /orders
```

**Response:**
```json
{
  "orders": [
    {
      "id": "order-123",
      "price": 100.0,
      "tax": 10.0,
      "final_price": 110.0
    },
    {
      "id": "order-456",
      "price": 200.0,
      "tax": 20.0,
      "final_price": 220.0
    }
  ]
}
```

### gRPC API

#### Create Order
```protobuf
rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);

message CreateOrderRequest {
  string id = 1;
  float price = 2;
  float tax = 3;
}

message CreateOrderResponse {
  string id = 1;
  float price = 2;
  float tax = 3;
  float final_price = 4;
}
```

#### List Orders
```protobuf
rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);

message ListOrdersRequest {
  // Empty for now, future: pagination parameters
}

message ListOrdersResponse {
  repeated CreateOrderResponse orders = 1;
}
```

### GraphQL API

#### Create Order (Mutation)
```graphql
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
```

#### List Orders (Query)
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

## 🏃‍♂️ Running the Application

1. **Start Dependencies:**
   ```bash
   docker-compose up -d  # MySQL + RabbitMQ
   ```

2. **Run the Application:**
   ```bash
   go run cmd/ordersystem/main.go
   ```

3. **Access Endpoints:**
   - **HTTP REST**: `http://localhost:8000`
   - **gRPC**: `localhost:50051`
   - **GraphQL Playground**: `http://localhost:8080`

## 🧪 Testing

### Run Tests
```bash
go test ./internal/usecase/... -v
```

### API Testing
Use the provided `.http` files in the `api/` directory:
- `create_order.http` - Test order creation
- `list_orders.http` - Test order listing

## 🏗️ Architecture

This API follows **Clean Architecture** principles:

- **Domain Layer**: Business entities and interfaces
- **Application Layer**: Use cases and DTOs
- **Infrastructure Layer**: Database, HTTP, gRPC, GraphQL implementations
- **Dependency Injection**: Google Wire for clean dependency management

## 📊 Features

- ✅ **Multi-Interface Support**: HTTP REST, gRPC, and GraphQL
- ✅ **Event-Driven Architecture**: RabbitMQ integration
- ✅ **Clean Architecture**: SOLID principles implementation
- ✅ **Comprehensive Testing**: Unit tests with mocks
- ✅ **Database Integration**: MySQL with prepared statements
- ✅ **Error Handling**: Proper error responses across all interfaces 