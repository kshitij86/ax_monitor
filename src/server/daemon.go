package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/*
	startDaemon()
		-> start one goroutine for each API
		-> get API status
		-> lock mutex, update status, unlock mutex
		-> push to kafka topic
		-> sigChan gets SIGTERM or SIGINT
		-> the cancel ctx is called and each goroutine is terminated
		-> each goroutine decrements the wait group when being terminated
		-> the wait group waits for all goroutines to exit and then ends the program

	the ctx ensures each goroutine gets the signal to exit
	when exiting, the wait group makes sure ALL goroutines exit BEFORE the main program exits

*/

var apiMap = make(map[string]ApiStatus)
var mu sync.Mutex // Mutex to ensure thread-safe access to apiStatus

var kafkaTopic string

// Function to probe API and update the status
func probeAPI(ctx context.Context, apiEndpoint string, wgProbe *sync.WaitGroup, wgKafka *sync.WaitGroup, kafkaProducer *kafka.Producer) {
	defer wgProbe.Done() // when the Probe goroutine exits, decrement wait group
	var api ApiStatus
	api.Api_endpoint = apiEndpoint

	for {
		select {
		case <-ctx.Done():
			// in case the ctx is cancelled
			fmt.Println("shutting down monitor for", apiEndpoint)
			return
		default:
			time.Sleep(time.Second)
			status := "up"

			resp, err := http.Get(apiEndpoint)
			if err != nil || resp.StatusCode == http.StatusNotFound {
				status = "down"
			}

			// Lock the mutex to update the shared map safely
			mu.Lock()
			api.Api_status = status
			if api.Api_status == "up" {
				api.Api_uptime += 1.0
			} else if api.Api_status == "down" {
				api.Api_uptime = 0.0
			}
			apiMap[apiEndpoint] = api
			mu.Unlock()

			// publish new status to kakfa using a new goroutine
			wgKafka.Add(1) // add a goroutine when publishing
			go func() {
				defer wgKafka.Done() // ensure it closes before it exits
				_, err := publishToKafkaTopic(kafkaTopic, api, kafkaProducer)
				if err != nil {
					panic(err)
				}
				// fmt.Println("record updated for endpoint: ", apiEndpoint, publishStatus, "on", time.Now())
			}()

		}
	}
}

func startDaemon(apiConfig APIConfig) {

	/* set the kafka topic */
	kafkaTopic = "test"

	/* initialize a kafka producer */
	kafkaProducer, err := createKafkaProducer()
	if err != nil {
		panic(err)
	}
	defer kafkaProducer.Close() // ensure producer closes as daemon finishes

	/* get the API root */
	apiRoot := getAPIRoot(apiConfig)

	/* create a wait group for all goroutines */
	var wgProbe sync.WaitGroup
	var wgKafka sync.WaitGroup

	/* the context will help to shutdown the goroutines */
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)                      // channel to handle OS signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // capture interrupt or termination

	// start the goroutines to monitor the APIs
	for i := 0; i < len(apiConfig.Axmonitor_API_List); i++ {
		apiEndpoint := apiRoot + "" + apiConfig.Axmonitor_API_List[i]    // get API endpoint
		wgProbe.Add(1)                                                   // add a new probe goroutine
		go probeAPI(ctx, apiEndpoint, &wgProbe, &wgKafka, kafkaProducer) // pass a cancelable context and pointer to wait group
	}

	<-sigChan // wait unitl sigChan gets one of the signals, hold the main thread

	cancel()       // cancel all contexts and return from all goroutines
	wgProbe.Wait() // wait for all goroutines to finish
}
