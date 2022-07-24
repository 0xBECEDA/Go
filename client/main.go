package main

import (
	"context"
	"messanger/client/http"
	"messanger/client/send_msg"
	"messanger/client/write_msg"
	"messanger/internal"
	"time"

	log "go.uber.org/zap"
)

func main() {
	l := log.NewExample()
	ctx := context.Background()

	server := http.New(l, &http.Handler{Logger: l})
	server.Start()
	defer server.Stop(5 * time.Second)

	name := write_msg.EnterUserName()
	ch := make(chan internal.Message, 100)

	go write_msg.GetInput(ch, name)
	go send_msg.SendMessage(ch, "localhost:8080")
	<-ctx.Done()
}
