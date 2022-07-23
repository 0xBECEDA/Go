package server

import (
	"context"
	"messanger/server/handlers"
	"runtime"
	"strconv"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Server struct {
	config       Config
	server       *fasthttp.Server
	routeHandler *handlers.RouteHandler
	logger       *zap.Logger
}

func New(cfg Config, logger *zap.Logger, handler *handlers.RouteHandler) *Server {
	s := &Server{
		config:       cfg,
		routeHandler: handler,
		server: &fasthttp.Server{
			TCPKeepalivePeriod: cfg.TCPAlivePeriod,
			MaxRequestsPerConn: cfg.MaxConn,
		},
		logger: logger,
	}

	r := router.New()
	r.GET("/health", func(reqCtx *fasthttp.RequestCtx) {
		_, _ = reqCtx.WriteString(strconv.Itoa(runtime.NumGoroutine()))
		reqCtx.SetStatusCode(fasthttp.StatusOK)
	})

	r.POST("/send", s.routeHandler.Send)
	s.server.Handler = r.Handler
	return s
}

func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(s.config.Address); err != nil {
			s.logger.Fatal("error during start server", zap.Error(err))
		}
	}()
}

func (s *Server) Stop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		if err := s.server.Shutdown(); err != nil {
			s.logger.Error("error stopping server", zap.Error(err))
		}
		cancel()
	}()
	<-ctx.Done()
}
