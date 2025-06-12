# Desafio Clean Architecture

## Como rodar o projeto

### Opção 1: Docker Compose (Recomendado)

Construa as imagens e rode todo o projeto com um único comando:

``` bash
rm -rf .docker                # Limpa o cache do docker
docker compose up --build     # constrói as imagens e sobe todos os serviços
```

Isso irá subir:
- **MySQL** na porta 3306
- **RabbitMQ** na porta 5672 (management na 15672)  
- **Aplicação Go** com:
  - HTTP REST API na porta 8000
  - gRPC na porta 50051
  - GraphQL na porta 9002

### Opção 2: Desenvolvimento Local

Se preferir rodar apenas as dependências no Docker e a aplicação localmente:

``` bash
# Sobe apenas as dependências
docker compose up mysql rabbitmq -d

# Roda a aplicação localmente
go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go
```

## Testando a aplicação

### REST API (porta 8000)

**Criar um pedido:**
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{
    "id":"123e4567-e89b-12d3-a456-426614174000",
    "price": 100.5,
    "tax": 0.5,
    "final_price": 101.0
  }'
```

**Listar pedidos:**
```bash
curl http://localhost:8000/orders
```

**Status da API:**
```bash
curl http://localhost:8000/
```

### GraphQL (porta 9002)

Acesse http://localhost:9002 no navegador para usar o GraphQL Playground.

**Exemplo de query:**
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

**Exemplo de mutation:**
```graphql
mutation {
  createOrder(input: {
    id: "new-order-123"
    Price: 199.99
    Tax: 20.0
    FinalPrice: 219.99
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

### gRPC (porta 50051)

Use um cliente gRPC como `grpcurl` ou `evans` para testar:

```bash
# Listar serviços disponíveis
grpcurl -plaintext localhost:50051 list

# Criar um pedido
grpcurl -plaintext -d '{
  "id": "grpc-order-001",
  "price": 150.0,
  "tax": 15.0,
  "final_price": 165.0
}' localhost:50051 pb.OrderService/CreateOrder

# Listar pedidos
grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
```

Respondendo:

```json
{
  "orders": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "price": 100.5,
      "tax": 0.5,
      "finalPrice": 101
    },
    {
      "id": "aewqeqewq3",
      "price": 100.5,
      "tax": 0.53,
      "finalPrice": 101.03
    },
    {
      "id": "grpc-order-001",
      "price": 150,
      "tax": 15,
      "finalPrice": 165
    },
    {
      "id": "order-001",
      "price": 100,
      "tax": 10,
      "finalPrice": 110
    },
    {
      "id": "order-002",
      "price": 250.5,
      "tax": 25.05,
      "finalPrice": 275.55
    },
    {
      "id": "order-003",
      "price": 75.25,
      "tax": 7.53,
      "finalPrice": 82.78
    }
  ]
}
```



### RabbitMQ Management

Acesse http://localhost:15672 no navegador:
- **Usuário:** guest
- **Senha:** guest

## Arquitetura

O projeto segue os princípios da Clean Architecture:

```
cmd/ordersystem/          # Main application
internal/
├── entity/              # Entidades de negócio
├── usecase/             # Casos de uso
├── infra/
│   ├── database/        # Implementações de banco de dados
│   ├── web/             # Handlers HTTP
│   ├── grpc/            # Serviços gRPC
│   └── graph/           # Resolvers GraphQL
└── event/               # Sistema de eventos
```

## Tecnologias utilizadas

- **Go 1.22**
- **MySQL 5.7** - Banco de dados principal
- **RabbitMQ** - Message broker para eventos
- **Chi Router** - HTTP router
- **gRPC** - Comunicação entre serviços
- **GraphQL** - API flexível para consultas
- **Wire** - Injeção de dependência
- **Docker & Docker Compose** - Containerização

## Funcionalidades

- ✅ **CRUD de Pedidos** via REST, gRPC e GraphQL
- ✅ **Sistema de Eventos** com RabbitMQ
- ✅ **Clean Architecture** com separação de responsabilidades
- ✅ **Injeção de Dependência** com Google Wire
- ✅ **Múltiplas interfaces** (HTTP, gRPC, GraphQL)
- ✅ **Containerização** completa com Docker
- ✅ **Health Checks** para todos os serviços

## Estrutura de dados

### Order (Pedido)
```json
{
  "id": "string",
  "price": "float64", 
  "tax": "float64",
  "final_price": "float64"
}
```

## Logs e Monitoramento

Para visualizar os logs da aplicação:

```bash
# Logs de todos os serviços
docker compose logs -f

# Logs apenas da aplicação
docker compose logs -f ordersystem

# Logs do MySQL
docker compose logs -f mysql

# Logs do RabbitMQ  
docker compose logs -f rabbitmq
```

## Parando os serviços

```bash
# Para todos os serviços
docker compose down

# Para e remove volumes (limpa dados)
docker compose down -v
```
