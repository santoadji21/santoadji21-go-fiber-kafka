package main

import (
	"fmt"
	"strconv"

	kafkaCfg "go-kafka/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func produceKafkaMessages(p *kafka.Producer, topic string, messages []string) {
	// Produce messages to topic (asynchronously)
	for _, message := range messages {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, nil)
	}

	// Wait for message deliveries
	p.Flush(15 * 1000)
}

func main() {
	app := fiber.New()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaCfg.BROKER_HOST})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

    // Fiber endpoint to produce messages with UUIDs
    app.Post("/produce", func(c *fiber.Ctx) error {
        var messages []string
        for i := 0; i < 7; i++ { // Assuming you want to send 7 messages
            id := uuid.New()
            message := "Message " + strconv.Itoa(i+1) + " with UUID: " + id.String()
            messages = append(messages, message)
        }

        go produceKafkaMessages(p, kafkaCfg.TOPIC, messages)
        return c.SendString("Messages with UUIDs are being produced...")
    })

	app.Listen(":3000")
}
