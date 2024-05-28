package service

import (
	"context"
	"log/slog"

	"github.com/markraiter/simple-blog/internal/model"
)

type PostSaver interface {
	SavePost(ctx context.Context, post *model.Post) (int, error)
}

type PostProvider interface {
	Post(ctx context.Context, id int) (*model.Post, error)
	Posts(ctx context.Context) ([]*model.Post, error)
}

type PostProcessor interface {
	UpdatePost(ctx context.Context, post *model.Post) error
	DeletePost(ctx context.Context, id int) error
}

type PostService struct {
	log       *slog.Logger
	saver     PostSaver
	provider  PostProvider
	processor PostProcessor
}

// implement the PostService methods here
// SavePost, Post, Posts, UpdatePost, DeletePost
