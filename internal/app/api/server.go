package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/markraiter/simple-blog/config"
)

type Server struct {
	HTTPServer *http.Server
	logger     *slog.Logger
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	s.HTTPServer = &http.Server{
		Addr:           cfg.Server.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
	}

	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
