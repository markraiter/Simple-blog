package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	"github.com/markraiter/simple-blog/internal/app/service"
	"github.com/markraiter/simple-blog/internal/lib/sl"
	"github.com/markraiter/simple-blog/internal/model"
)

type CommentSaver interface {
	SaveComment(ctx context.Context, userID int, commentReq *model.CommentRequest) (int, error)
}

type CommentProvider interface{}

type CommentProcessor interface{}

type CommentHandler struct {
	log       *slog.Logger
	validate  *validator.Validate
	saver     CommentSaver
	provider  CommentProvider
	processor CommentProcessor
}

// @Summary Create a comment
// @Description Create a comment
// @Security ApiKeyAuth
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body model.CommentRequest true "Comment object that needs to be created"
// @Success 201 {string} string "Comment created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handler.CreateComment"
		log := h.log.With(slog.String("operation", operation))

		var commentReq model.CommentRequest
		userID := middleware.GetUserIDFromCtx(r.Context())

		if err := json.NewDecoder(r.Body).Decode(&commentReq); err != nil {
			log.Warn("error parsing request", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := h.validate.Struct(commentReq); err != nil {
			log.Warn("error validating comment", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		id, err := h.saver.SaveComment(ctx, userID, &commentReq)
		if err != nil {
			if errors.Is(err, service.ErrPostNotExists) {
				log.Warn("error saving comment", sl.Err(err))
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			log.Error("error saving comment", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(strconv.Itoa(id))) //nolint:errcheck
	}
}
