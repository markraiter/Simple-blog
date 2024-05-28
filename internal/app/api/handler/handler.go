package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

type AuthService interface {
	Auth
}

type Handler struct {
	Healthcheck
	AuthHandler
}

// The response struct is used to send a message back to the client.
type response struct {
	Message string `json:"message"`
}

func New(log *slog.Logger, validate *validator.Validate, auth AuthService) *Handler {
	return &Handler{
		Healthcheck{log: log},
		AuthHandler{log: log, validate: validate, service: auth},
	}
}

func (h *Handler) Router(ctx context.Context, cfg config.Config) http.Handler {
	m := http.NewServeMux()

	m.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	m.Handle("GET /health", h.APIHealth())
	m.Handle("POST /auth/register", h.RegisterUser(ctx))
	m.Handle("POST /auth/login", h.Login(ctx, cfg.Auth))

	return m
}
