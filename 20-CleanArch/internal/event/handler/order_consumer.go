package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type OrderConsumer struct {
	RabbitMQChannel *amqp.Channel
	QueueName       string
}

func NewOrderConsumer(rabbitMQChannel *amqp.Channel, queueName string) *OrderConsumer {
	return &OrderConsumer{
		RabbitMQChannel: rabbitMQChannel,
		QueueName:       queueName,
	}
}

func (c *OrderConsumer) StartConsuming() error {
	// Declare the queue (in case it doesn't exist)
	queue, err := c.RabbitMQChannel.QueueDeclare(
		c.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		amqp.Table{
			"x-message-ttl": 86400000, // 24 hours
			"x-max-length":  10000,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind the queue to amq.direct exchange (for compatibility with existing handler)
	err = c.RabbitMQChannel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key (empty for amq.direct)
		"amq.direct", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Start consuming messages
	msgs, err := c.RabbitMQChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	// Process messages
	go func() {
		log.Printf("🐰 [OrderConsumer] Started consuming from queue: %s", c.QueueName)
		for msg := range msgs {
			c.processMessage(msg)
		}
	}()

	return nil
}

func (c *OrderConsumer) processMessage(msg amqp.Delivery) {
	log.Printf("📨 [OrderConsumer] Received message: %s", string(msg.Body))

	// Parse the order data
	var orderData map[string]interface{}
	if err := json.Unmarshal(msg.Body, &orderData); err != nil {
		log.Printf("❌ [OrderConsumer] Failed to parse message: %v", err)
		return
	}

	// Process the order (this is where you'd add your business logic)
	if orderID, ok := orderData["ID"].(string); ok {
		log.Printf("✅ [OrderConsumer] Processing order: %s", orderID)

		// Example processing logic
		c.processOrder(orderData)
	} else {
		log.Printf("⚠️ [OrderConsumer] Invalid order format: missing ID")
	}
}

func (c *OrderConsumer) processOrder(orderData map[string]interface{}) {
	// This is where you'd implement your order processing logic
	// For example:
	// - Send email notifications
	// - Update inventory
	// - Trigger fulfillment processes
	// - Log to audit systems

	log.Printf("🔄 [OrderConsumer] Order processing completed for: %v", orderData["ID"])

	// You could also publish to other exchanges/queues here
	// Example: Publish to a fulfillment queue
	c.publishToFulfillment(orderData)
}

func (c *OrderConsumer) publishToFulfillment(orderData map[string]interface{}) {
	// Example: Publish processed order to fulfillment queue
	fulfillmentData := map[string]interface{}{
		"order_id":     orderData["ID"],
		"status":       "ready_for_fulfillment",
		"processed_at": "2024-01-01T00:00:00Z", // You'd use actual timestamp
	}

	jsonData, _ := json.Marshal(fulfillmentData)

	err := c.RabbitMQChannel.Publish(
		"orders.direct",     // exchange
		"order.fulfillment", // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)

	if err != nil {
		log.Printf("❌ [OrderConsumer] Failed to publish to fulfillment: %v", err)
	} else {
		log.Printf("📤 [OrderConsumer] Published to fulfillment queue")
	}
}
