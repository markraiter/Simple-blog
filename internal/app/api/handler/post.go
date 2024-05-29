package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	"github.com/markraiter/simple-blog/internal/lib/sl"
	"github.com/markraiter/simple-blog/internal/model"
)

type PostSaver interface {
	SavePost(ctx context.Context, userID int, postReq *model.PostRequest) (int, error)
}

type PostProvider interface {
	Post(ctx context.Context, id int) (*model.Post, error)
	Posts(ctx context.Context) ([]*model.Post, error)
}

type PostProcessor interface {
	UpdatePost(ctx context.Context, id int, postReq *model.PostRequest) error
	DeletePost(ctx context.Context, id int) error
}

type PostHandler struct {
	log       *slog.Logger
	validate  *validator.Validate
	saver     PostSaver
	provider  PostProvider
	processor PostProcessor
}

// @Summary Create a post
// @Description Create a post
// @Security ApiKeyAuth
// @Tags posts
// @Accept json
// @Produce json
// @Param post body model.PostRequest true "Post object that needs to be created"
// @Success 201 {string} string "Post created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/posts [post]
func (ph *PostHandler) CreatePost(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const operation = "handler.CreatePost"

		log := ph.log.With(slog.String("operation", operation))

		var postReq model.PostRequest
		userIDStr, ok := r.Context().Value(middleware.UIDKey).(string)
		if !ok {
			log.Warn("error getting userID from context")
			http.Error(w, "error getting userID from context", http.StatusInternalServerError)

			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Warn("error parsing userID", sl.Err(err))
			http.Error(w, "error parsing userID", http.StatusInternalServerError)

			return
		}

		if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
			log.Warn("error parsing request", sl.Err(err))
			http.Error(w, "error parsing request", http.StatusBadRequest)

			return
		}

		if err := ph.validate.Struct(postReq); err != nil {
			log.Warn("error validating post", sl.Err(err))
			http.Error(w, "error validating post", http.StatusBadRequest)

			return
		}

		id, err := ph.saver.SavePost(ctx, userID, &postReq)
		if err != nil {
			log.Error("error saving post", sl.Err(err))
			http.Error(w, "error saving post", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(id)))
	}
}
