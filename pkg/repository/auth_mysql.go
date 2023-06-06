package repository

import (
	"database/sql"
	"fmt"

	"github.com/markraiter/simple-blog/models"
)

type AuthMySQL struct {
	db *sql.DB
}

func NewAuthMySQL(db *sql.DB) *AuthMySQL {
	return &AuthMySQL{
		db: db,
	}
}

func (r *AuthMySQL) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, password) VALUES (?, ?)", usersTable)
	result, err := r.db.Exec(query, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	id = int(insertID)

	return id, nil
}

func (r *AuthMySQL) GetUser(email, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = ? AND password = ?", usersTable)
	err := r.db.QueryRow(query, email, password).Scan(&user.Email, &user.Password)

	return user, err
}
