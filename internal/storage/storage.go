package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/markraiter/simple-blog/internal/storage/postgres"
	"github.com/markraiter/simple-blog/models"
)

type Authentication interface {
	Create(user *models.User) (uint, error)
	GetEmail(email string) string
	GetUserByEmail(email, password string) (*models.User, error)
	GenerateToken(email, password string) (string, error)
}

type Posts interface {
	GetAll() ([]models.Post, error)
	Filter(userID uint) ([]models.Post, error)
	Get(id uint) (*models.Post, error)
	Create(post *models.Post, userID uint) (uint, error)
	Update(id uint, input *models.UpdatePostInput) error
	Delete(id uint) error
}

type Comments interface {
	GetAll() ([]models.Comment, error)
	FilterByPost(postID uint) ([]models.Comment, error)
	FilterByUser(userID uint) ([]models.Comment, error)
	Get(id uint) (*models.Comment, error)
	Create(comment *models.Comment, userID, postID uint) (uint, error)
	Update(id uint, input *models.UpdateCommentInput) error
	Delete(id uint) error
}

type Storage struct {
	Authentication
	Posts
	Comments
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		postgres.NewAuth(db),
		postgres.NewPost(db),
		postgres.NewComment(db),
	}
}
