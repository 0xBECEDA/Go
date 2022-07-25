package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"messanger/client/http"
	"messanger/client/reader"
	"messanger/client/sender"
	"messanger/internal"
	"time"

	log "go.uber.org/zap"
)

const serverHost = "localhost:8080"

func authorizeLoop(myHost string) *internal.AuthorizeMessage {
	for {
		msg := reader.Authorize(myHost)
		respCode, err := sender.Authorize(msg, serverHost)
		if err != nil && errors.Is(err, internal.ErrAuthorization) {
			msg = reader.SignUp(myHost)
			respCode, err := sender.SignUp(msg, serverHost)
			if err != nil {
				fmt.Errorf("error during sign up: %v, please try 5 seconds later", err)
				continue
			}
			if respCode == fasthttp.StatusOK {
				return msg
			}
		}
		if err != nil {
			fmt.Errorf("error during authorisation: %v, please try 5 seconds later", err)
		}
		if respCode == fasthttp.StatusOK {
			return msg
		} else {
		}
	}
}

func main() {
	l := log.NewExample()
	ctx := context.Background()

	myHost := reader.EnterHost()

	server := http.New(l, &http.Handler{Logger: l}, myHost)
	server.Start()
	defer server.Stop(5 * time.Second)

	msg := authorizeLoop(myHost)

	ch := make(chan internal.Message, 100)
	go reader.GetInput(ch, msg.Name)
	go sender.SendMessage(ch, serverHost)
	<-ctx.Done()
}
