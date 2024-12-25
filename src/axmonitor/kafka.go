package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/* kafka related funcitonality */

func structToBytes(structVar ApiStatus) ([]byte, error) {
	/* serialize a struct to bytes to add to a topic */
	structBytes, err := json.Marshal(structVar)
	if err != nil {
		return structBytes, err
	}
	return structBytes, nil
}

func createKafkaProducer() (*kafka.Producer, error) {
	/* create a new kafka producer and return pointer to it */

	// create a new kafka producer
	createdProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return createdProducer, err
	}
	// return a pointer to the kafka producer
	return createdProducer, nil
}

func createKafkaConsumer() (*kafka.Consumer, error) {
	/* create a new kafka consumer and return pointer to it */

	createdConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		return createdConsumer, err
	}

	return createdConsumer, nil
}

func consumeLastMessageFromTopic(topics []string, kafkaConsumer *kafka.Consumer) ([]byte, error) {
	/* reads the latest message from the given topic */
	kafkaConsumer.SubscribeTopics(topics, nil) // subscribe to the given topics

	var consumedMessage []byte

	msg, err := kafkaConsumer.ReadMessage(-1) // read the last message
	if err != nil {
		return consumedMessage, err
	}

	consumedMessage = msg.Value

	return consumedMessage, nil
}

func publishToKafkaTopic(topic string, apiStatus ApiStatus, kafkaProducer *kafka.Producer) (bool, error) {

	/* convert the apiStatus to bytes */
	apiStatusBytes, err := structToBytes(apiStatus)
	if err != nil {
		return false, err
	}

	// nil here mentions the chan
	kafkaProducer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: apiStatusBytes}, nil)

	return true, nil
}
