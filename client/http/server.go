package http

import (
	"context"
	"runtime"
	"strconv"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Server struct {
	server        *fasthttp.Server
	logger        *zap.Logger
	getMsgHandler *Handler
	host          string
}

func New(logger *zap.Logger, handler *Handler, host string) *Server {
	s := &Server{
		server: &fasthttp.Server{
			TCPKeepalivePeriod: 20 * time.Second,
			MaxRequestsPerConn: 24,
		},
		logger:        logger,
		getMsgHandler: handler,
		host:          host,
	}

	r := router.New()
	r.GET("/health", func(reqCtx *fasthttp.RequestCtx) {
		_, _ = reqCtx.WriteString(strconv.Itoa(runtime.NumGoroutine()))
		reqCtx.SetStatusCode(fasthttp.StatusOK)
	})

	r.POST("/get_msg", s.getMsgHandler.GetMessage)
	s.server.Handler = r.Handler
	return s
}

func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(s.host); err != nil {
			s.logger.Fatal("error during start http", zap.Error(err))
		}
	}()
}

func (s *Server) Stop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		if err := s.server.Shutdown(); err != nil {
			s.logger.Error("error stopping http", zap.Error(err))
		}
		cancel()
	}()
	<-ctx.Done()
}
