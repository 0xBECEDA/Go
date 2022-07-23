package main

import (
	"context"
	"messanger/config"
	"messanger/db"
	"messanger/server/handlers"
	"messanger/server/server"
	"time"

	logger "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	l := logger.NewExample()
	ctx := context.Background()
	var dbCfg db.Config
	if err := config.Load("./db/config.yml", &dbCfg); err != nil {
		l.Fatal("error loading db config", zapcore.Field{String: err.Error()})
	}

	dbConn, err := db.Connect(dbCfg, l)
	if err != nil {
		l.Fatal("error connecting", zapcore.Field{String: err.Error()})
	}

	var serverCfg server.Config
	if err := config.Load("../config.yml", &serverCfg); err != nil {
		l.Fatal("error loading db config", zapcore.Field{String: err.Error()})
	}

	handler := handlers.New(l, dbConn)
	serv := server.New(serverCfg, l, handler)
	serv.Start()
	defer serv.Stop(5 * time.Second)

	<-ctx.Done()
	l.Info("stopping service")
	time.Sleep(2 * time.Second)
}
