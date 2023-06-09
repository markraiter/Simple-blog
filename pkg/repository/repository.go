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
	GetAll() ([]models.Post, error)
	GetByID(postID int) (models.Post, error)
	Update(postID int, input models.UpdatePostInput) error
	Delete(postID int) error
}

type Comments interface {
	Create(postID int, comment models.Comment) (int, error)
	GetAll() ([]models.Comment, error)
	GetByID(commentID int) (models.Comment, error)
	Update(commentID int, input models.UpdateCommentInput) error
	Delete(commentID int) error
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
		Comments:      NewCommentPostgres(db),
	}
}
