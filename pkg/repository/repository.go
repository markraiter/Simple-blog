package repository

import (
	"github.com/jmoiron/sqlx"

	"github.com/markraiter/simple-blog/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (models.User, error)
}

type Posts interface {
	Create(userID int, post models.Post) (int, error)
}

type Comments interface {
	Create(postID int, comment models.Comment) (int, error)
}

type Repository struct {
	Authorization
	Posts
	Comments
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Posts:         NewPostPostgres(db),
	}
}
