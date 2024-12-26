package axmonitor

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/* kafka related funcitonality */

func getKafkaTopic() string {
	return "test"
}

func structToBytes(structVar ApiStatus) ([]byte, error) {
	/* serialize a struct to bytes to add to a topic */
	structBytes, err := json.Marshal(structVar)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the struct")
	}
	return structBytes, nil
}

func createKafkaProducer() (*kafka.Producer, error) {
	/* create a new kafka producer and return pointer to it */

	// create a new kafka producer
	createdProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return nil, err
	}
	// return a pointer to the kafka producer
	return createdProducer, nil
}

func createKafkaConsumer() (*kafka.Consumer, error) {
	/* create a new kafka consumer and return pointer to it */

	createdConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "test-group",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		return nil, err
	}

	return createdConsumer, nil
}

func consumeMessageFromTopic(kafkaConsumer *kafka.Consumer) ([]byte, error) {
	/* reads the latest message from the given topic */

	kafkaConsumer.Subscribe(getKafkaTopic(), nil)

	var consumedMessage []byte

	msg, err := kafkaConsumer.ReadMessage(-1) // read a message
	if err != nil {
		return consumedMessage, err
	}
	var v ApiStatus
	json.Unmarshal(msg.Value, &v)
	fmt.Println("consuming: ", v)

	/*
		TODO: this consumed message is not able to get the lastupdated field
	*/

	consumedMessage = msg.Value

	return consumedMessage, nil
}

func publishToKafkaTopic(apiStatus ApiStatus, kafkaProducer *kafka.Producer) (bool, error) {

	/* convert the apiStatus to bytes */
	apiStatusBytes, err := structToBytes(apiStatus)
	if err != nil {
		return false, err
	}

	/* this runs fine */
	var temp ApiStatus
	json.Unmarshal(apiStatusBytes, &temp)
	fmt.Println("publishing: ", temp)

	// nil here mentions the chan
	topicName := getKafkaTopic()
	kafkaProducer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny}, Value: apiStatusBytes}, nil)

	return true, nil
}
