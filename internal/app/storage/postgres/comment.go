package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
)

func (s *Storage) SaveComment(ctx context.Context, comment *model.Comment) (int, error) {
	const operation = "storage.SaveComment"

	tx, err := s.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	query := "INSERT INTO comments (content, post_id, user_id) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.UserID).Scan(&comment.ID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	updateQuery := "UPDATE posts SET comments_count = comments_count + 1 WHERE id = $1"
	_, err = tx.ExecContext(ctx, updateQuery, comment.PostID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return comment.ID, nil
}

// Comment returns any comment by its ID.
func (s *Storage) Comment(ctx context.Context, id int) (*model.Comment, error) {
	const operation = "storage.Comment"

	query, err := s.PostgresDB.Prepare("SELECT id, content, post_id, user_id FROM comments WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	row := query.QueryRowContext(ctx, id)

	comment := &model.Comment{}
	err = row.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return comment, nil
}

// CommentsByPost returns all comments for provided post.
func (s *Storage) CommentsByPost(ctx context.Context, postID int) ([]*model.Comment, error) {
	const operation = "storage.CommentsByPost"

	query, err := s.PostgresDB.Prepare("SELECT id, content, post_id, user_id FROM comments WHERE post_id = $1 ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	rows, err := query.QueryContext(ctx, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	defer rows.Close()

	comments := make([]*model.Comment, 0)
	for rows.Next() {
		comment := &model.Comment{}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// UpdateComment updates a comment by its ID.
//
// If the comment does not exist it returns storage.ErrNotFound.
// If the user is not the author of the comment it returns storage.ErrNotAllowed.
// If the post does not exist it returns storage.ErrNotFound.
func (s *Storage) UpdateComment(ctx context.Context, comment *model.Comment) error {
	const operation = "storage.UpdateComment"

	query := `
        UPDATE comments 
        SET content = $1 
        WHERE id = $2 AND user_id = $3
        RETURNING id
    `

	var updatedCommentID int
	err := s.PostgresDB.QueryRowContext(ctx, query, comment.Content, comment.ID, comment.UserID).Scan(&updatedCommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			commentExistsQuery := "SELECT id FROM comments WHERE id = $1"
			var existsCommentID int
			err = s.PostgresDB.QueryRowContext(ctx, commentExistsQuery, comment.ID).Scan(&existsCommentID)
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

// func (s *Storage) DeleteComment(ctx context.Context, commentID, userID int) error {
// 	const operation = "storage.DeleteComment"

// 	query := `
//         DELETE FROM comments
//         WHERE id = $1 AND user_id = $2
//         RETURNING id
//     `
// 	var deletedCommentID int
// 	err := s.PostgresDB.QueryRowContext(ctx, query, commentID, userID).Scan(&deletedCommentID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			postExistsQuery := "SELECT id FROM comments WHERE id = $1"
// 			var existsCommentID int
// 			err := s.PostgresDB.QueryRowContext(ctx, postExistsQuery, commentID).Scan(&existsCommentID)
// 			if err != nil {
// 				if errors.Is(err, sql.ErrNoRows) {
// 					return fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
// 				}
// 				return fmt.Errorf("%s: %w", operation, err)
// 			}
// 			return fmt.Errorf("%s: %w", operation, storage.ErrNotAllowed)
// 		}
// 		return fmt.Errorf("%s: %w", operation, err)
// 	}

// 	return nil
// }

func (s *Storage) DeleteComment(ctx context.Context, commentID, userID int) error {
	const operation = "storage.DeleteComment"

	tx, err := s.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
		RETURNING id
	`
	var deletedCommentID int
	err = tx.QueryRowContext(ctx, query, commentID, userID).Scan(&deletedCommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			commentExistsQuery := "SELECT id FROM comments WHERE id = $1"
			var existsCommentID int

			err = s.PostgresDB.QueryRowContext(ctx, commentExistsQuery, commentID).Scan(&existsCommentID)
			if err != nil {
				tx.Rollback()
				if errors.Is(err, sql.ErrNoRows) {
					tx.Rollback()

					return fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
				}
				tx.Rollback()

				return fmt.Errorf("%s: %w", operation, err)
			}
			tx.Rollback()

			return fmt.Errorf("%s: %w", operation, storage.ErrNotAllowed)
		}
		tx.Rollback()

		return fmt.Errorf("%s: %w", operation, err)
	}

	updateQuery := "UPDATE posts SET comments_count = comments_count - 1 WHERE id = (SELECT post_id FROM comments WHERE id = $1)"
	_, err = tx.ExecContext(ctx, updateQuery, commentID)
	if err != nil {
		tx.Rollback()

		return fmt.Errorf("%s: %w", operation, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
