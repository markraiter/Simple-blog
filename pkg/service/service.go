package service

import (
	"github.com/markraiter/simple-blog/models"
	"github.com/markraiter/simple-blog/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Posts interface {
	Create(userID int, post models.Post) (int, error)
	GetAll() ([]models.Post, error)
	GetByID(postID int) (models.Post, error)
	Update(postID int, input models.UpdatePostInput) error
	Delete(postID int) error
}

type Comments interface {
}

type Service struct {
	Authorization
	Posts
	Comments
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Posts:         NewPostService(repos.Posts),
	}
}
