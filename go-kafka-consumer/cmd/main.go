package main

import (
	"fmt"
	"sync"
	"time"

	kafkaCfg "go-kafka/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
)

var (
	// Message store
	messages []string
	// Mutex to avoid race conditions
	mu sync.Mutex
)

func consumeKafkaMessages() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaCfg.BROKER_HOST,
		"group.id":          kafkaCfg.GROUP_ID,
		"auto.offset.reset": kafkaCfg.AUTO_RESET_OFFSET,
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{kafkaCfg.TOPIC}, nil)

	run := true
	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			mu.Lock()
			messages = append(messages, string(msg.Value))
			mu.Unlock()
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func main() {
	app := fiber.New()

	// Start Kafka consumer in a goroutine
	go consumeKafkaMessages()

	// Fiber endpoint to display messages
	app.Get("/messages", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		return c.JSON(messages)
	})

	app.Listen(":3002")
}
