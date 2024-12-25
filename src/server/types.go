package main

type ApiStatus struct {
	Api_endpoint string  `json:"api_endpoint"`
	Api_status   string  `json:"api_status"`
	Api_uptime   float64 `json:"api_uptime"`
}

type DatabaseServerConfig struct {
	Mysql_username string `json:"MYSQL_USERNAME"`
	Mysql_password string `json:"MYSQL_PASSWORD"`
	Mysql_port     string `json:"MYSQL_PORT"`
	Mysql_host     string `json:"MYSQL_HOST"`
	Mysql_db       string `json:"MYSQL_DB"`
}

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

type User struct {
	user_id  string
	username string
	password string
}
