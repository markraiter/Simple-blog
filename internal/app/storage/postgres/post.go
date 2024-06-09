package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
)

func (s *Storage) SavePost(ctx context.Context, post *model.Post) (int, error) {
	const operation = "storage.SavePost"

	query := "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id"

	err := s.PostgresDB.QueryRow(query, post.Title, post.Content, post.UserID).Scan(&post.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return post.ID, nil
}

func (s *Storage) Post(ctx context.Context, id int) (*model.Post, error) {
	const operation = "storage.Post"

	query, err := s.PostgresDB.Prepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	row := query.QueryRowContext(ctx, id)

	post := &model.Post{}

	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CommentsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return post, nil
}

func (s *Storage) Posts(ctx context.Context) ([]*model.Post, error) {
	const operation = "storage.Posts"

	query, err := s.PostgresDB.Prepare("SELECT id, title, content, user_id, comments_count FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	defer rows.Close()

	posts := make([]*model.Post, 0)
	for rows.Next() {
		post := &model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CommentsCount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost updates post by its ID.
//
// If the post does not exist it returns storage.ErrNotFound.
// If the post does not belong to the user it returns storage.ErrNotAllowed.
func (s *Storage) UpdatePost(ctx context.Context, post *model.Post) error {
	const operation = "storage.UpdatePost"

	query := `
        UPDATE posts 
        SET title = $1, content = $2 
        WHERE id = $3 AND user_id = $4
        RETURNING id
    `

	var updatedPostID int
	err := s.PostgresDB.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.UserID).Scan(&updatedPostID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			postExistsQuery := "SELECT id FROM posts WHERE id = $1"
			var existsPostID int
			err := s.PostgresDB.QueryRowContext(ctx, postExistsQuery, post.ID).Scan(&existsPostID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
				}
				return fmt.Errorf("%s: %w", operation, err)
			}
			return fmt.Errorf("%s: %w", operation, storage.ErrNotAllowed)
		}
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (s *Storage) DeletePost(ctx context.Context, postID, userID int) error {
	const operation = "storage.DeletePost"

	query := `
        DELETE FROM posts 
        WHERE id = $1 AND user_id = $2 
        RETURNING id
    `

	var deletedPostID int
	err := s.PostgresDB.QueryRowContext(ctx, query, postID, userID).Scan(&deletedPostID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			postExistsQuery := "SELECT id FROM posts WHERE id = $1"
			var existsPostID int
			err := s.PostgresDB.QueryRowContext(ctx, postExistsQuery, postID).Scan(&existsPostID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
				}
				return fmt.Errorf("%s: %w", operation, err)
			}
			return fmt.Errorf("%s: %w", operation, storage.ErrNotAllowed)
		}
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

// func (s *Storage) UpdatePost(ctx context.Context, post *model.Post) error {
// 	const operation = "storage.UpdatePost"
//
// 	query := "UPDATE posts SET title = $1, content = $2 WHERE id = $3"
//
// 	_, err := s.PostgresDB.Exec(query, post.Title, post.Content, post.ID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
// 		}
//
// 		return fmt.Errorf("%s: %w", operation, err)
// 	}
//
// 	return nil
// }
//
// func (s *Storage) DeletePost(ctx context.Context, id int) error {
// 	const operation = "storage.DeletePost"
//
// 	query := "DELETE FROM posts WHERE id = $1"
//
// 	_, err := s.PostgresDB.Exec(query, id)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
// 		}
//
// 		return fmt.Errorf("%s: %w", operation, err)
// 	}
//
// 	return nil
// }
