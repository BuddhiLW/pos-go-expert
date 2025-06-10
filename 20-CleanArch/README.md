# Desafio Clean Architecture

## Como rodar o projeto


Construa a imagens e rode-as, em fundo:

``` bash
rm -rf .docker                # Limpa o cache do docker
docker compose up --build -d  # constrói as imagens e sobe os containers em plano de fundo
```

Rode o projeto:

``` bash
go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go
```

## Chamadas


### HTTP

#### Criar ordem

``` bash
꧂ (λ) curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-order-1",
    "price": 100.0,
    "tax": 10.0
  }'
```

```
{"id":"test-order-1","price":100,"tax":10,"final_price":110}
```

#### Listar ordens

``` bash
curl -s http://localhost:8000/orders -H "Content-Type: application/json" | jq .
``` 

``` output
{
  "orders": [
    {
      "id": "eqwewq3",
      "price": 1999.99,
      "tax": 200,
      "final_price": 2199.99
    },
    {
      "id": "order-001",
      "price": 100,
      "tax": 10,
      "final_price": 110
    },
    {
      "id": "order-002",
      "price": 250.5,
      "tax": 25.05,
      "final_price": 275.55
    },
    {
      "id": "order-003",
      "price": 75.25,
      "tax": 7.53,
      "final_price": 82.78
    }
  ]
}
```


### GraphQL

#### Criar ordem

``` bash
curl -s http://localhost:9002/query -H "Content-Type: application/json" -d '{"query":"mutation { createOrder(input: {id: \"gql-working\", Price: 500.0, Tax: 50.0}) { id Price Tax FinalPrice } }"}' | jq .
```

``` output
Order created: {gql-working 500 50 550}{
  "data": {
    "createOrder": {
      "id": "gql-working",
      "Price": 500,
      "Tax": 50,
      "FinalPrice": 550
    }
  }
}
```

#### Listar ordens

``` bash
curl -s http://localhost:9002/query -H "Content-Type: application/json" -d '{"query": "query { orders { id Price Tax FinalPrice } }"}' | jq .
``` 

```
{
  "data": {
    "orders": [
      {
        "id": "manual-test",
        "Price": 100,
        "Tax": 10,
        "FinalPrice": 110
      },
      {
        "id": "order-001",
        "Price": 100,
        "Tax": 10,
        "FinalPrice": 110
      },
      {
        "id": "order-002",
        "Price": 250.5,
        "Tax": 25.05,
        "FinalPrice": 275.55
      },
      {
        "id": "order-003",
        "Price": 75.25,
        "Tax": 7.53,
        "FinalPrice": 82.78
      },
      {
        "id": "test-123",
        "Price": 100,
        "Tax": 10,
        "FinalPrice": 110
      }
    ]
  }
}
```

### gRPC


```bash
evans internal/infra/grpc/protofiles/order.proto
```

``` output
  ______
 |  ____|
 | |__    __   __   __ _   _ __    ___
 |  __|   \ \ / /  / _. | | '_ \  / __|
 | |____   \ V /  | (_| | | | | | \__ \
 |______|   \_/    \__,_| |_| |_| |___/

 more expressive universal gRPC client


pb.OrderService@127.0.0.1:50051> call CreateOrder
id (TYPE_STRING) => eqwewq3
price (TYPE_FLOAT) => 1999.99
tax (TYPE_FLOAT) => 200.00
{
  "finalPrice": 2199.99,
  "id": "eqwewq3",
  "price": 1999.99,
  "tax": 200
}

pb.OrderService@127.0.0.1:50051> call ListOrders
{
  "orders": [
    {
      "finalPrice": 2199.99,
      "id": "eqwewq3",
      "price": 1999.99,
      "tax": 200
    },
    {
      "finalPrice": 110,
      "id": "order-001",
      "price": 100,
      "tax": 10
    },
    {
      "finalPrice": 275.55,
      "id": "order-002",
      "price": 250.5,
      "tax": 25.05
    },
    {
      "finalPrice": 82.78,
      "id": "order-003",
      "price": 75.25,
      "tax": 7.53
    }
  ]
```
