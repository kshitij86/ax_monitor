package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

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
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return p, err
	}
	// return a pointer to the kafka producer
	return p, nil
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
