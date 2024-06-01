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
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
			log.Warn("error parsing request", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := ph.validate.Struct(postReq); err != nil {
			log.Warn("error validating post", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		id, err := ph.saver.SavePost(ctx, userID, &postReq)
		if err != nil {
			log.Error("error saving post", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(strconv.Itoa(id))) //nolint:errcheck
	}
}

// @Summary Get a Post
// @Description Get a Post
// @Tags posts
// @Accept json
// @Produce json
// @Param id query int true "Post ID"
// @Success 200 {object} model.Post
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/posts/{id} [get]
func (ph *PostHandler) Post(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const operation = "handler.GetPost"

        log := ph.log.With(slog.String("operation", operation))

        idStr := r.URL.Query().Get("id")
        if idStr == "" {
            log.Warn("error getting id from query")
            http.Error(w, "error getting id from query", http.StatusBadRequest)

            return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
            log.Warn("error parsing id", sl.Err(err))
            http.Error(w, err.Error(), http.StatusBadRequest)

            return
        }

        post, err := ph.provider.Post(ctx, id)
        if err != nil {
            if errors.Is(err, service.ErrNotFound) {
                log.Warn("post not found", sl.Err(err))
                http.Error(w, err.Error(), http.StatusNotFound)

                return
            }

            log.Error("error getting post", sl.Err(err))
            http.Error(w, err.Error(), http.StatusInternalServerError)

            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        if err := json.NewEncoder(w).Encode(post); err != nil {
            log.Error("error encoding post", sl.Err(err))
            http.Error(w, err.Error(), http.StatusInternalServerError)

            return
        }
    }
}

// @Summary Get Posts
// @Description Get Posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} model.Post
// @Failure 500 {string} string "Internal server error"
// @Router /api/posts [get]
func (ph *PostHandler) Posts(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const operation = "handler.GetPosts"

        log := ph.log.With(slog.String("operation", operation))

        posts, err := ph.provider.Posts(ctx)
        if err != nil {
            log.Error("error getting posts", sl.Err(err))
            http.Error(w, err.Error(), http.StatusInternalServerError)

            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        if err := json.NewEncoder(w).Encode(posts); err != nil {
            log.Error("error encoding posts", sl.Err(err))
            http.Error(w, err.Error(), http.StatusInternalServerError)

            return
        }
    }
}

// @Summary Update a UpdatePost
// @Description Update a UpdatePost
// @Security ApiKeyAuth
// @Tags posts
// @Accept json
// @Produce json
// @Param id query int true "Post ID"
// @Param post body model.PostRequest true "Post object that needs to be updated"
// @Success 200 {string} string "Post updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/posts/{id} [put]
func (ph *PostHandler) UpdatePost(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const operation = "handler.UpdatePost"

        log := ph.log.With(slog.String("operation", operation))

        idStr := r.URL.Query().Get("id")
        if idStr == "" {
            log.Warn("error getting id from query")
            http.Error(w, "error getting id from query", http.StatusBadRequest)

            return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
            log.Warn("error parsing id", sl.Err(err))
            http.Error(w, err.Error(), http.StatusBadRequest)

            return
        }

        var postReq model.PostRequest
        if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
            log.Warn("error parsing request", sl.Err(err))
            http.Error(w, err.Error(), http.StatusBadRequest)

            return
        }

        if err := ph.validate.Struct(postReq); err != nil {
            log.Warn("error validating post", sl.Err(err))
            http.Error(w, err.Error(), http.StatusBadRequest)

            return
        }

        err = ph.processor.UpdatePost(ctx, id, &postReq)
        if err != nil {
            if errors.Is(err, service.ErrNotFound) {
                log.Warn("post not found", sl.Err(err))
                http.Error(w, err.Error(), http.StatusNotFound)

                return
            }

            log.Error("error updating post", sl.Err(err))
            http.Error(w, err.Error(), http.StatusInternalServerError)

            return
        }

        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("Post updated")) //nolint:errcheck
    }
}
