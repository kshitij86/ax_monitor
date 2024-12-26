package axmonitor

import (
	"context"
	"fmt"
	"net/http"
	"sync"
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

// Function to probe API and update the status
func probeAPI(ctx context.Context, apiEndpoint string, wgProbe *sync.WaitGroup, kafkaProducer *kafka.Producer) {
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

			// probe the API
			resp, err := http.Get(apiEndpoint)
			if err != nil || resp.StatusCode == http.StatusNotFound {
				status = "down"
			}
			api.Api_status = status

			// set/reset the uptime
			if api.Api_status == "up" {
				api.Api_uptime += 1.0
			} else if api.Api_status == "down" {
				api.Api_uptime = 0.0
			}

			// set the last updated
			api.Api_lastupdated = time.Now().String()

			// Lock the mutex to update the shared map safely
			// mu.Lock()
			// apiMap[apiEndpoint] = api
			// mu.Unlock()

			// publish new status to kakfa using a new goroutine
			// send the update api struct to the kafka topic
			go func() {
				_, err := publishToKafkaTopic(api, kafkaProducer)
				if err != nil {
					panic(err)
				}
			}()
		}
	}
}

func startDaemon(ctx context.Context, apiConfig APIConfig) {
	/* create a wait group for all goroutines */
	var wgProbe sync.WaitGroup

	select {
	case <-ctx.Done():
		return
	default:
		/* set the kafka topic */

		/* initialize a kafka producer */
		kafkaProducer, err := createKafkaProducer()
		if err != nil {
			panic(err)
		}
		defer kafkaProducer.Close() // ensure producer closes as daemon finishes

		/* get the API root */
		apiRoot := getAPIRoot(apiConfig)

		// start the goroutines to monitor the APIs
		for i := 0; i < len(apiConfig.Axmonitor_API_List); i++ {
			apiEndpoint := apiRoot + "" + apiConfig.Axmonitor_API_List[i] // get API endpoint
			wgProbe.Add(1)                                                // add a new probe goroutine
			go probeAPI(ctx, apiEndpoint, &wgProbe, kafkaProducer)        // pass a cancelable context and pointer to wait group
		}
	}
	wgProbe.Wait() // wait for all goroutines to finish
}
