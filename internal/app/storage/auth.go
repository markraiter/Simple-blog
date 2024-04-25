package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/markraiter/simple-blog/internal/model"
)

func (s *Storage) SaveUser(ctx context.Context, user *model.User) (int, error) {
	const operation = "Storage.SaveUser"

	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id"
	err := s.PostgresDB.QueryRow(query, user.Username, user.Password, user.Email).Scan(&user.ID)
	if err != nil {
		var pgErr *pq.Error

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", operation, model.ErrUserAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return user.ID, nil
}

func (s *Storage) User(ctx context.Context, email string) (*model.User, error) {
	const operation = "Storage.UserByEmail"

	query, err := s.PostgresDB.Prepare("SELECT id, username, password, email FROM users WHERE email = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	row := query.QueryRowContext(ctx, email)

	user := &model.User{}
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", operation, model.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return user, nil
}
