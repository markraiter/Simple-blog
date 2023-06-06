package service

import (
	"github.com/markraiter/simple-blog/models"
	"github.com/markraiter/simple-blog/pkg/repository"
)

type PostService struct {
	repo repository.Posts
}

func NewPostService(repo repository.Posts) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) Create(userID int, post models.Post) (int, error) {
	return s.repo.Create(userID, post)
}
