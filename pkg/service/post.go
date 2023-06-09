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

func (s *PostService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetByID(postID int) (models.Post, error) {
	return s.repo.GetByID(postID)
}

func (s *PostService) Update(postID int, input models.UpdatePostInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(postID, input)
}

func (s *PostService) Delete(postID int) error {
	return s.repo.Delete(postID)
}
