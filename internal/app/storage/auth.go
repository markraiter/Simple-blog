package storage

import (
	"fmt"

	"github.com/markraiter/simple-blog/internal/model"
)

func (s *Storage) SaveUser(user *model.User) (int, error) {
	const operation = "Storage.SaveUser"

	query := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id`
	err := s.PostgresDB.QueryRow(query, user.Username, user.Password, user.Email).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return user.ID, nil
}

func (s *Storage) UserByEmail(email string) (*model.User, error) {
	const operation = "Storage.UserByEmail"

	query := `SELECT id, username, password, email FROM users WHERE email = $1`

	user := &model.User{}

	err := s.PostgresDB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return user, nil
}
