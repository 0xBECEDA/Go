package route

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"messanger/db"

	"github.com/valyala/fasthttp"
)

var (
	ErrBanned         = errors.New("you are banned")
	ErrNoSuchUser     = errors.New("no user is found")
	ErrFriendNotFound = errors.New("your friend is not found")
)

type MessageFrom struct {
	FromID int    `json:"from_id"`
	ToID   int    `json:"to_id"`
	Data   string `json:"data"`
}

type Handler struct {
	logger *zap.Logger
	dbConn *db.DB
	client *fasthttp.Client
	// TODO добавить коннект к реббиту
}

func New(logger *zap.Logger, db *db.DB) *Handler {
	return &Handler{
		logger: logger,
		dbConn: db,
		client: &fasthttp.Client{},
	}
}

func (h *Handler) Send(reqCtx *fasthttp.RequestCtx) {
	var (
		msg     MessageFrom
		accFrom db.Account
		accTo   db.Account
	)

	data := reqCtx.Request.Body()
	if err := jsoniter.Unmarshal(data, &msg); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
		return
	}

	if err := h.dbConn.FindAccountByID(msg.FromID, &accFrom); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	if accFrom.ID <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrNoSuchUser.Error()))
		return
	}

	if accFrom.Banned {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrBanned.Error()))
		return
	}

	if err := h.dbConn.FindAccountByID(msg.ToID, &accTo); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	if accTo.ID <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusNoContent)
		_, _ = reqCtx.Write([]byte(ErrFriendNotFound.Error()))
		return
	}

	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://" + accTo.Host)
	req.SetBody(data)

	if err := h.client.Do(req, resp); err != nil {
		reqCtx.SetStatusCode(resp.StatusCode())
		_, _ = reqCtx.Write([]byte(err.Error()))
		return
	}
}
