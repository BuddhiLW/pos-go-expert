# Third-Party Infrastructure Setup

This directory contains configuration and initialization files for the external dependencies of the Order System.

## 🗄️ MySQL Database

### Files:
- `mysql/init.sql` - Database initialization script that creates tables and sample data
- `mysql/my.cnf` - MySQL configuration for performance optimization

### Schema:
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

### Sample Data:
The initialization script includes 3 sample orders for testing purposes.

## 🐰 RabbitMQ Message Broker

RabbitMQ is configured with automatic queue and exchange setup for event-driven communication.

### Files:
- `rabbitmq/definitions.json` - RabbitMQ definitions for exchanges, queues, and bindings
- `rabbitmq/rabbitmq.conf` - RabbitMQ server configuration

### Messaging Architecture:

#### **Exchanges:**
- `amq.direct` - Built-in direct exchange (used by OrderCreatedHandler)
- `orders.direct` - Custom direct exchange for order-related routing
- `orders.topic` - Topic exchange for pattern-based routing

#### **Queues:**
- `orders.created` - Receives all order creation events from `amq.direct`
- `orders.events` - General order events queue for business logic processing

#### **Bindings:**
- `amq.direct` → `orders.created` (routing key: "")
- `orders.direct` → `orders.events` (routing key: "order.created")
- `orders.topic` → `orders.events` (routing key: "order.*")

### Message Flow:
```
Order Creation
     ↓
OrderCreatedHandler → amq.direct → orders.created queue
     ↓
OrderConsumer processes messages
     ↓
Publishes to orders.direct → orders.events queue
```

### Queue Properties:
- **Durability**: All queues are durable (survive server restarts)
- **TTL**: Messages expire after 24 hours (86400000ms)
- **Max Length**: 10,000 messages per queue
- **Auto-delete**: Disabled (queues persist when no consumers)

### Access:
- **Management UI**: http://localhost:15672
- **Default credentials**: guest/guest
- **AMQP Port**: 5672

## 🚀 Quick Start

1. **Start the infrastructure:**
   ```bash
   docker compose up -d
   ```

2. **Verify MySQL is ready:**
   ```bash
   docker compose logs mysql
   ```

3. **Verify RabbitMQ is ready:**
   ```bash
   docker compose logs rabbitmq
   ```

4. **Connect to MySQL:**
   ```bash
   docker exec -it mysql mysql -u root -p orders
   # Password: root
   ```

5. **Verify tables:**
   ```sql
   SHOW TABLES;
   SELECT * FROM orders;
   ```

6. **Check RabbitMQ queues:**
   - Open http://localhost:15672 in your browser
   - Login with guest/guest
   - Navigate to "Queues" tab to see pre-configured queues

## 🔧 Configuration

### Environment Variables:
```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
```

### Ports:
- **MySQL**: 3306
- **RabbitMQ AMQP**: 5672
- **RabbitMQ Management**: 15672

## 🏥 Health Checks

Both services include health checks:
- **MySQL**: `mysqladmin ping`
- **RabbitMQ**: `rabbitmq-diagnostics ping`

## 📦 Data Persistence

Data is persisted in Docker volumes:
- MySQL data: `.docker/mysql/`
- RabbitMQ data: `.docker/rabbitmq/`

## 🔄 Reset Database

To reset the database and start fresh:
```bash
docker compose down -v
docker compose up -d
```

This will remove all data and reinitialize with the sample data.

## 📨 Testing Message Flow

### 1. **Create an Order (triggers message)**
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"test-order","price":100.0,"tax":10.0}'
```

### 2. **Check Messages in RabbitMQ UI**
- Go to http://localhost:15672
- Navigate to Queues → `orders.created`
- You should see message count increase

### 3. **View Message Content**
- Click on the `orders.created` queue
- Click "Get Messages" to see message payloads

### 4. **Monitor Logs**
```bash
docker compose logs -f
```

Look for consumer processing logs:
```
📨 [OrderConsumer] Received message: {"ID":"test-order",...}
✅ [OrderConsumer] Processing order: test-order
```

## 🐛 Troubleshooting

### MySQL Connection Issues:
```bash
# Check if MySQL is running
docker compose ps mysql

# View MySQL logs
docker compose logs mysql

# Connect to MySQL container
docker exec -it mysql bash
```

### RabbitMQ Issues:
```bash
# Check RabbitMQ status
docker compose ps rabbitmq

# View RabbitMQ logs
docker compose logs rabbitmq

# Access RabbitMQ management UI
open http://localhost:15672

# Check queue definitions loaded
docker exec rabbitmq rabbitmqctl list_queues
docker exec rabbitmq rabbitmqctl list_exchanges
```

### Message Flow Issues:
```bash
# Check if queues are declared
docker exec rabbitmq rabbitmqctl list_queues name messages

# Check bindings
docker exec rabbitmq rabbitmqctl list_bindings

# Monitor real-time messages
# (Use RabbitMQ Management UI for real-time monitoring)
``` 