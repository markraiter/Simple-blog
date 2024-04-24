package handler

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	Healthcheck
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		Healthcheck{log: log},
	}
}

func (h *Handler) Router() http.Handler {
	m := http.NewServeMux()

	m.Handle("GET /health", h.APIHealth())

	return m
}
