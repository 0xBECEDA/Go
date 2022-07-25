package register

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"messanger/db"
	"messanger/internal"
)

var (
	ErrEmptyEmail    = errors.New("empty email address")
	ErrEmptyUserName = errors.New("empty user name")
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

func (h *Handler) RegisterNewUser(reqCtx *fasthttp.RequestCtx) {
	var msg internal.AuthorizeMessage
	if err := jsoniter.Unmarshal(reqCtx.Request.Body(), &msg); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	if len(msg.Email) <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrEmptyEmail.Error()))
		return
	}

	if len(msg.Name) <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrEmptyUserName.Error()))
		return
	}

	if err := h.dbConn.CreateAccount(msg.Name, msg.Email, msg.Host); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
		return
	}

	reqCtx.SetStatusCode(fasthttp.StatusOK)
	_, _ = reqCtx.Write([]byte("account registered successfully"))
	return
}
