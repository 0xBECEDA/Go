package server

import (
	"context"
	"messanger/server/handlers/authorize"
	"messanger/server/handlers/register"
	"messanger/server/handlers/route"
	"runtime"
	"strconv"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Server struct {
	config           Config
	server           *fasthttp.Server
	routeHandler     *route.Handler
	regHandler       *register.Handler
	authorizeHandler *authorize.Handler
	logger           *zap.Logger
}

func New(cfg Config, logger *zap.Logger, handler *route.Handler, handler2 *register.Handler, handler3 *authorize.Handler) *Server {
	s := &Server{
		config:           cfg,
		routeHandler:     handler,
		regHandler:       handler2,
		authorizeHandler: handler3,
		server: &fasthttp.Server{
			TCPKeepalivePeriod: 20 * time.Second,
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
	r.POST("/reg", s.regHandler.RegisterNewUser)
	r.POST("/authorize", s.authorizeHandler.Authorize)

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
