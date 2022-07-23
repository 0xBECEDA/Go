package db

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBUser   string `json:"db_user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}
