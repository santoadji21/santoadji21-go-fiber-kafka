package config

const (
	// KafkaBroker is the Kafka broker address
	BROKER_HOST = "localhost:9092"
	// KafkaTopic is the Kafka topic to produce to
	TOPIC = "myTopic"
	// KafkaGroupID is the Kafka consumer group ID
	GROUP_ID = "myGroup"
	// KafkaAutoOffsetReset is the Kafka consumer auto offset reset
	AUTO_RESET_OFFSET = "earliest"
	// KafkaTimeout is the Kafka consumer timeout
	TIMEOUT = 15
)