package authorize

import (
	"messanger/db"
	"messanger/internal"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	dbConn *db.DB
}

func New(logger *zap.Logger, dbConn *db.DB) *Handler {
	return &Handler{
		logger: logger,
		dbConn: dbConn,
	}
}

func (h *Handler) Authorize(reqCtx *fasthttp.RequestCtx) {
	var (
		msg internal.AuthorizeMessage
		acc db.Account
	)
	if err := jsoniter.Unmarshal(reqCtx.Request.Body(), &msg); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	if err := h.dbConn.FindAccountByUserName(msg.Name, &acc); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}
	if acc.ID <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte("please register!"))
	}

	if err := h.dbConn.UpdateAccountHost(acc.ID, msg.Host); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}
	reqCtx.SetStatusCode(fasthttp.StatusOK)
}
