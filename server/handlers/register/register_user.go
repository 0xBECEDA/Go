package register

import (
	"errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"messanger/db"
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
	email := reqCtx.QueryArgs().Peek("email")
	if len(string(email)) <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrEmptyEmail.Error()))
		return
	}

	userName := reqCtx.QueryArgs().Peek("name")
	if len(string(userName)) <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrEmptyUserName.Error()))
		return
	}

	if err := h.dbConn.CreateAccount(string(userName), string(email)); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(err.Error()))
		return
	}

	reqCtx.SetStatusCode(fasthttp.StatusOK)
	_, _ = reqCtx.Write([]byte("account registered successfully"))
	return
}
