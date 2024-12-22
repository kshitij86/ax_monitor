package main

// handles calling each API using a goroutine

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Api struct {
	Api_status string
}

var apiStatus = make(map[string]Api)
var mu sync.Mutex // Mutex to ensure thread-safe access to apiStatus

// Function to probe API and update the status
func probeAPI(apiEndpoint string) {
	var api Api

	status := "up"
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		status = "down"
	}
	defer resp.Body.Close()

	// Lock the mutex to update the shared map safely
	mu.Lock()
	api.Api_status = status
	apiStatus[apiEndpoint] = api
	mu.Unlock()
}

func startMainServer(apiConfig APIConfig) {

	/*
		dashboard calls main server and the main server
		keeps monitoring APIs
	*/

	/*
		TODO: make the server constantly monitor all connections
		and send info when the API is called
	*/

	// mute all warnings
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	apiRoot := getAPIRoot(apiConfig)

	for i := 0; i < len(apiConfig.Axmonitor_API_List); i++ {
		apiEndpoint := apiRoot + "" + apiConfig.Axmonitor_API_List[i]
		go probeAPI(apiEndpoint)
	}

	// Define a route
	r.GET("/get-status", func(c *gin.Context) {
		// returns a status list of all API endpoints
		for key, value := range apiStatus {
			c.JSON(200, gin.H{"apiEndpoint": key, "status": value.Api_status})
		}
	})

	// Start the server
	r.Run(":8080") // Run on port 8080
}
