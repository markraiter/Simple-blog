package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/markraiter/simple-blog/internal/model"
)

type Healthcheck struct {
	log *slog.Logger
}

func (h *Healthcheck) APIHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := model.Response{Message: "healthy"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}
}
