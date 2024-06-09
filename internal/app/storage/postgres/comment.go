package postgres

import (
	"context"
	"fmt"

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
