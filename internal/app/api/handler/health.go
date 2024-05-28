package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Healthcheck struct {
	log *slog.Logger
}

// @Summary Healthcheck
// @Description Healthcheck
// @Tags health
// @Produce json
// @Success 200 {object} response "healthy"
// @Router /health [get]
func (h *Healthcheck) APIHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{Message: "healthy"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}
}
