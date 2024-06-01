package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/markraiter/simple-blog/internal/app/storage"
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
	saver     PostSaver
	provider  PostProvider
	processor PostProcessor
}

func (ps *PostService) SavePost(ctx context.Context, userID int, postReq *model.PostRequest) (int, error) {
	const operation = "service.SavePost"

	postModel := model.Post{
		Title:   postReq.Title,
		Content: postReq.Content,
		UserID:  userID,
	}

	id, err := ps.saver.SavePost(ctx, &postModel)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return id, nil
}

func (ps *PostService) Post(ctx context.Context, id int) (*model.Post, error) {
	const operation = "service.Post"

	post, err := ps.provider.Post(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", operation, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return post, nil
}

func (ps *PostService) Posts(ctx context.Context) ([]*model.Post, error) {
	const operation = "service.Posts"

	posts, err := ps.provider.Posts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return posts, nil
}

func (ps *PostService) UpdatePost(ctx context.Context, postID, userID int, postReq *model.PostRequest) error {
	const operation = "service.UpdatePost"

    post, err := ps.provider.Post(ctx, postID)
    if err != nil {
        return fmt.Errorf("%s: %w", operation, err)
    }

    if post.UserID != userID {
        return fmt.Errorf("%s: %w", operation, ErrNotAllowed)
    }

	postModel := model.Post{
		ID:      postID,
		Title:   postReq.Title,
		Content: postReq.Content,
	}

	err = ps.processor.UpdatePost(ctx, &postModel)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("%s: %w", operation, ErrNotFound)
		}

		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (ps *PostService) DeletePost(ctx context.Context, postID, userID int) error {
	const operation = "service.DeletePost"

    post, err := ps.provider.Post(ctx, postID)
    if err != nil {
        return fmt.Errorf("%s: %w", operation, err)
    }

    if post.UserID != userID {
        return fmt.Errorf("%s: %w", operation, ErrNotAllowed)
    }

	err = ps.processor.DeletePost(ctx, postID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("%s: %w", operation, ErrNotFound)
		}

		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
