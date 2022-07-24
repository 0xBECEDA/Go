package main

import (
	"context"
	"messanger/config"
	"messanger/db"
	"messanger/server/handlers/authorize"
	"messanger/server/handlers/register"
	"messanger/server/handlers/route"
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
	if err := config.Load("./server/config.yml", &serverCfg); err != nil {
		l.Fatal("error loading db config", zapcore.Field{String: err.Error()})
	}

	routeHandler := route.New(l, dbConn)
	regHandler := register.New(l, dbConn)
	authorizeHandeler := authorize.New(l, dbConn)

	serv := server.New(serverCfg, l, routeHandler, regHandler, authorizeHandeler)
	serv.Start()
	defer serv.Stop(5 * time.Second)

	<-ctx.Done()
	l.Info("stopping service")
	time.Sleep(2 * time.Second)
}
