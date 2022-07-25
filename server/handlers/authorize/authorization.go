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
		reqCtx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	if err := h.dbConn.FindAccountByUserName(msg.Name, &acc); err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	if acc.ID <= 0 {
		reqCtx.Error(internal.ErrAuthorization.Error(), fasthttp.StatusBadRequest)
		return
	}

	if err := h.dbConn.UpdateAccountHost(acc.ID, msg.Host); err != nil {
		reqCtx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	reqCtx.SetStatusCode(fasthttp.StatusOK)
}
