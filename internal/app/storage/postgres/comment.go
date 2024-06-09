package postgres

import (
	"context"
	"fmt"

	"github.com/markraiter/simple-blog/internal/model"
)

func (s *Storage) SaveComment(ctx context.Context, comment *model.Comment) (int, error) {
	const operation = "storage.SaveComment"

	query := "INSERT INTO comments (content, post_id, user_id) VALUES ($1, $2, $3) RETURNING id"
	err := s.PostgresDB.QueryRow(query, comment.Content, comment.PostID, comment.UserID).Scan(&comment.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return comment.ID, nil
}
