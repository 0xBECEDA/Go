package main

import (
	"context"
	"fmt"
	"messanger/client/http"
	"messanger/client/send_msg"
	"messanger/client/write_msg"
	"messanger/internal"
	"time"

	"github.com/valyala/fasthttp"
	log "go.uber.org/zap"
)

const addr = "localhost:8080"

func main() {
	l := log.NewExample()
	ctx := context.Background()

	host := write_msg.EnterHost()

	server := http.New(l, &http.Handler{Logger: l}, host)
	server.Start()
	defer server.Stop(5 * time.Second)

	msg := write_msg.Authorize(host)
	respCode, err := send_msg.Authorize(msg, host)
	if err != nil || respCode != fasthttp.StatusOK {
		fmt.Errorf("error during authorisation: response code %v, error %v", respCode, err)
		return
	}
	ch := make(chan internal.Message, 100)
	go write_msg.GetInput(ch, msg.Name)
	go send_msg.SendMessage(ch, addr)
	<-ctx.Done()
}
