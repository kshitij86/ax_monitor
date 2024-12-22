package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

var USER_ID string

// iota starts at 0 and autoincrements for each value that follows
// closest thing to an enum in go
// const (
// 	loginUser = iota
// 	createUser
// )

type Project struct {
	Project_id   int
	Project_name string
	User_id      string
}

type UserConfig struct {
	Axmonitor_username string   `json:"AXMONITOR_USERNAME"`
	Axmonitor_API_Port int      `json:"AXMONITOR_API_PORT"`
	Axmonitor_API_Host string   `json:"AXMONITOR_API_HOST"`
	Axmonitor_API_List []string `json:"AXMONITOR_API_LIST"`
}

type APIConfig struct {
	Axmonitor_API_Port int      `json:"AXMONITOR_API_PORT"`
	Axmonitor_API_Host string   `json:"AXMONITOR_API_HOST"`
	Axmonitor_API_List []string `json:"AXMONITOR_API_LIST"`
}

// func getUserProjects(db *sql.DB, user_id string) ([]Project, error) {

// 	var foundProjects []Project

// 	projects, err := db.Query("SELECT * FROM ax_monitor_projects WHERE user_id = ?", user_id)
// 	if err != nil {
// 		return foundProjects, fmt.Errorf("error loading projects")
// 	}

// 	for projects.Next() {
// 		var proj Project
// 		if err := projects.Scan(&proj.Project_id, &proj.Project_name, &proj.User_id); err != nil {
// 			return foundProjects, fmt.Errorf("error deserializing projects")
// 		}
// 		foundProjects = append(foundProjects, proj)
// 	}
// 	return foundProjects, nil
// }

// func setUserConfigAndLogin(db *sql.DB) (string, bool, error) {

// 	file, err := os.Open("user_config.json")
// 	if err != nil {
// 		return "", false, fmt.Errorf("cannot find user_config file")
// 	}
// 	defer file.Close()

// 	var usrConfig UserConfig

// 	decoder := json.NewDecoder(file)
// 	err = decoder.Decode(&usrConfig)
// 	if err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return "", false, fmt.Errorf("error reading user_config file")
// 	} else {
// 		fmt.Println("read user config successfully...")
// 	}

// 	user_id, isUserReal, err := LoginUser(db, usrConfig.Axmonitor_username)
// 	if err != nil {
// 		panic(err)
// 	}

// 	//handle user
// 	if isUserReal {
// 		fmt.Println("logged in as: ", usrConfig.Axmonitor_username)
// 		fmt.Printf("%s", fmt.Sprintf("monitoring APIs on http://%s:%d", usrConfig.Axmonitor_API_Host, usrConfig.Axmonitor_API_Port))
// 		return user_id, true, nil
// 	}

// 	return "", false, fmt.Errorf("error setting user configs")
// }

func getAPIRoot(apiConfig APIConfig) string {
	return fmt.Sprintf("http://%s:%d", apiConfig.Axmonitor_API_Host, apiConfig.Axmonitor_API_Port)
}

func setAPIConfig() (APIConfig, bool, error) {

	var apiConfig APIConfig
	file, err := os.Open("config/api_config.json")
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

type DatabaseServerConfig struct {
	Mysql_username string `json:"MYSQL_USERNAME"`
	Mysql_password string `json:"MYSQL_PASSWORD"`
	Mysql_port     string `json:"MYSQL_PORT"`
	Mysql_host     string `json:"MYSQL_HOST"`
	Mysql_db       string `json:"MYSQL_DB"`
}

func main() {
	fmt.Println(ax_monitor_logo)

	fmt.Println("Starting AX_MONITOR...")

	file, err := os.Open("config/server_config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var dbServerConfig DatabaseServerConfig

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dbServerConfig)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("read database server config successfully...")
	}

	dbPath := fmt.Sprintf("%s:%s@tcp(%s:%s)/ax_monitor_db", dbServerConfig.Mysql_username, dbServerConfig.Mysql_password, dbServerConfig.Mysql_host, dbServerConfig.Mysql_port)
	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Conncted to database...")
	}

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

	startMainServer(apiConfig)

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

	defer db.Close()
}
