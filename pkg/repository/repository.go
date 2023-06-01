package repository

import (
	"database/sql"

	"github.com/markraiter/simple-blog/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (models.User, error)
}

type Posts interface {
}

type Comments interface {
}

type Repository struct {
	Authorization
	Posts
	Comments
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMySQL(db),
	}
}
