package http

import (
	"fmt"
	"messanger/internal"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger
}

func (h *Handler) GetMessage(reqCtx *fasthttp.RequestCtx) {
	var msg internal.Message
	if err := jsoniter.Unmarshal(reqCtx.Request.Body(), &msg); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	reqCtx.SetStatusCode(fasthttp.StatusOK)

	fmt.Printf("You got message from user %s!\n", msg.FromName)
	fmt.Println(string(msg.Data))
	return
}
