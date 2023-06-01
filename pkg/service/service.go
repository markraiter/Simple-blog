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
}

type Comments interface {
}

type Service struct {
	Authorization
	Posts
	Comments
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
