package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
)

type CommentSaver interface {
	SaveComment(ctx context.Context, comment *model.Comment) (int, error)
}

type CommentProvider interface {
	Comment(ctx context.Context, id int) (*model.Comment, error)
	CommentsByPost(ctx context.Context, postID int) ([]*model.Comment, error)
}

type CommentProcessor interface {
	// UpdateComment(ctx context.Context, comment *model.Comment) error
	// DeleteComment(ctx context.Context, commentID, userID int) error
}

type CommentService struct {
	saver     CommentSaver
	provider  CommentProvider
	processor CommentProcessor
}

func (s *CommentService) SaveComment(ctx context.Context, userID int, commentReq *model.CommentRequest) (int, error) {
	const operation = "service.SaveComment"

	commentModel := model.Comment{
		Content: commentReq.Content,
		PostID:  commentReq.PostID,
		UserID:  userID,
	}

	id, err := s.saver.SaveComment(ctx, &commentModel)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return id, nil
}

func (s *CommentService) Comment(ctx context.Context, id int) (*model.Comment, error) {
	const operation = "service.Comment"

	comment, err := s.provider.Comment(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", operation, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return comment, nil
}

func (s *CommentService) CommentsByPost(ctx context.Context, postID int) ([]*model.Comment, error) {
	const operation = "service.CommentsByPost"

	comments, err := s.provider.CommentsByPost(ctx, postID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", operation, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return comments, nil
}
