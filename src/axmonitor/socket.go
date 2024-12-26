package axmonitor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

var kafkaConsumer *kafka.Consumer

func readAndPushToSocket(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from any origin (for development)
			return true
		}}

	/* upgrade the HTTP connection to a websocket one */
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close() // close the connection when the function exits

	/* read every second from the kafka topic */
	for {
		time.Sleep(time.Second)
		consumedMessage, err := consumeMessageFromTopic(kafkaConsumer)
		if err != nil {
			panic(err)
		}

		err = conn.WriteMessage(websocket.BinaryMessage, consumedMessage)
		if err != nil {
			log.Println("WriteMessage error:", err)
			panic(err)
		}
	}
}

func startWebSocketDaemon(ctx context.Context) {

	// Create a new Kafka consumer
	createdConsumer, err := createKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %s", err)
	}
	kafkaConsumer = createdConsumer
	defer kafkaConsumer.Close()

	// Create an HTTP server with graceful shutdown support
	server := &http.Server{Addr: ":3000"}

	http.HandleFunc("/ws", readAndPushToSocket)

	// Start the server in a goroutine
	fmt.Println("Started WebSocket server on localhost:3000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %s", err)
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown the server gracefully
	fmt.Println("Shutting down WebSocket server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("WebSocket server shutdown failed: %s", err)
	} else {
		fmt.Println("WebSocket server stopped cleanly.")
	}
}
