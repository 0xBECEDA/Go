package main

import (
	"log"
	"messanger/config"
	"messanger/db"
)

func main() {
	var cfg db.Config
	if err := config.Load("./db/config.yml", &cfg); err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Connect(cfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	dbConn.Conn.AutoMigrate(&db.Account{})
}
