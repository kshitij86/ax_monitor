package axmonitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// iota starts at 0 and autoincrements for each value that follows
// closest thing to an enum in go
// const (
// 	loginUser = iota
// 	createUser
// )

func getAPIRoot(apiConfig APIConfig) string {
	return fmt.Sprintf("http://%s:%d", apiConfig.Axmonitor_API_Host, apiConfig.Axmonitor_API_Port)
}

func setAPIConfig() (APIConfig, bool, error) {

	absPath, _ := os.Getwd()
	// fmt.Println(filepath.Join(absPath, "../src/axmonitor/config/api_config.json"))

	var apiConfig APIConfig

	file, err := os.Open(filepath.Join(absPath, "../axmonitor/config/api_config.json"))
	if err != nil {
		return apiConfig, false, fmt.Errorf("cannot find api_config file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&apiConfig)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return apiConfig, false, fmt.Errorf("error reading api_config file")
	} else {
		fmt.Println("read api config successfully...")
	}

	//handle user
	fmt.Println("monitoring APIs on: ", getAPIRoot(apiConfig))

	return apiConfig, true, nil

}

func StartMonitor() {
	fmt.Println(ax_monitor_logo)

	fmt.Println("Starting AX_MONITOR...")

	apiConfig, apiConfigStatus, err := setAPIConfig()
	if err != nil {
		panic(err)
	}
	if apiConfigStatus {
		fmt.Println("API config set...")
	} else {
		fmt.Println("cannot set API config, aborting...")
		os.Exit(1)
	}

	// send this context to both functions
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // called when the function exits, and stopChan gets SIGTERM

	go startDaemon(ctx, apiConfig)
	go startWebSocketDaemon(ctx)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	fmt.Println("Press Ctrl+C to exit...")

	<-stopChan // Wait for the signal
	fmt.Println("Shutting down gracefully...")
	// Perform any cleanup here
	// cancel() is called here
}

/*


// perform test query
	// ucnt := CountUsers(db)
	// fmt.Println("AX_MONITOR USER COUNT: ", ucnt)

	// user_id, setConfigStatus, err := setUserConfig(db)
	// if err != nil {
	// 	panic(err)
	// }

	// if setConfigStatus {
	// 	USER_ID = user_id
	// } else {
	// 	fmt.Println("cannot find user, aborting...")
	// 	os.Exit(1)
	// }

	// foundProjects, err := getUserProjects(db, USER_ID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(foundProjects)

	// defer db.Close()



*/
