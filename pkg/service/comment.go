package service

import (
	"github.com/markraiter/simple-blog/models"
	"github.com/markraiter/simple-blog/pkg/repository"
)

type CommentService struct {
	repo repository.Comments
}

func NewCommentService(repo repository.Comments) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (s *CommentService) Create(postID int, comment models.Comment) (int, error) {
	return s.repo.Create(postID, comment)
}
func (s *CommentService) GetAll() ([]models.Comment, error) {
	return s.repo.GetAll()
}
func (s *CommentService) GetByID(commentID int) (models.Comment, error) {
	return s.repo.GetByID(commentID)
}
func (s *CommentService) Update(commentID int, input models.UpdateCommentInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(commentID, input)
}
func (s *CommentService) Delete(commentID int) error {
	return s.repo.Delete(commentID)
}
