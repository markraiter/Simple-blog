package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockCommentSaver struct{ mock.Mock }

func (m *MockCommentSaver) SaveComment(ctx context.Context, comment *model.Comment) (int, error) {
	args := m.Called(ctx, comment)
	return args.Int(0), args.Error(1)
}

type MockCommentProvider struct{ mock.Mock }

func (m *MockCommentProvider) Comment(ctx context.Context, id int) (*model.Comment, error) {
	args := m.Called(ctx, id)
	comment := args.Get(0)
	if comment == nil {
		return nil, args.Error(1)
	}
	return comment.(*model.Comment), args.Error(1)
}

func (m *MockCommentProvider) CommentsByPost(ctx context.Context, postID int) ([]*model.Comment, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

// Tests
func TestCommentService_SaveComment(t *testing.T) {
	const operation = "service.SaveComment"
	var err = errors.New("error")

	mockSaver := new(MockCommentSaver)
	commentService := &CommentService{saver: mockSaver}

	tests := []struct {
		name       string
		commentReq *model.CommentRequest
		userID     int
		mockReturn int
		mockError  error
		wantError  error
	}{
		{
			name: "Success",
			commentReq: &model.CommentRequest{
				Content: "Test Content",
				PostID:  1,
			},
			userID:     1,
			mockReturn: 1,
			mockError:  nil,
			wantError:  nil,
		},
		{
			name: "Error",
			commentReq: &model.CommentRequest{
				Content: "Test Content",
				PostID:  1,
			},
			userID:     0,
			mockReturn: 0,
			mockError:  err,
			wantError:  fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver.On("SaveComment", mock.Anything, &model.Comment{
				Content: tt.commentReq.Content,
				PostID:  tt.commentReq.PostID,
				UserID:  tt.userID,
			}).Return(tt.mockReturn, tt.mockError)

			_, err := commentService.SaveComment(context.Background(), tt.userID, tt.commentReq)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockSaver.AssertExpectations(t)
		})
	}
}

func TestCommentService_Comment(t *testing.T) {
	const operation = "service.Comment"
	var err = errors.New("error")

	mockProvider := new(MockCommentProvider)
	commentService := &CommentService{provider: mockProvider}

	tests := []struct {
		name        string
		id          int
		mockReturn  *model.Comment
		mockError   error
		wantComment *model.Comment
		wantError   error
	}{
		{
			name: "Success",
			id:   1,
			mockReturn: &model.Comment{
				ID:      1,
				Content: "Test Content",
				PostID:  1,
				UserID:  1,
			},
			mockError: nil,
			wantComment: &model.Comment{
				ID:      1,
				Content: "Test Content",
				PostID:  1,
				UserID:  1,
			},
			wantError: nil,
		},
		{
			name:        "Comment Not Found",
			id:          2,
			mockReturn:  nil,
			mockError:   storage.ErrNotFound,
			wantComment: nil,
			wantError:   fmt.Errorf("%s: %w", operation, ErrNotFound),
		},
		{
			name:        "Error",
			id:          3,
			mockReturn:  nil,
			mockError:   err,
			wantComment: nil,
			wantError:   fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider.On("Comment", mock.Anything, tt.id).Return(tt.mockReturn, tt.mockError)

			_, err := commentService.Comment(context.Background(), tt.id)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantComment, tt.mockReturn)
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestCommentService_CommentsByPost(t *testing.T) {
	const operation = "service.CommentsByPost"
	var err = errors.New("error")

	mockProvider := new(MockCommentProvider)
	commentService := &CommentService{provider: mockProvider}

	tests := []struct {
		name         string
		ctx          context.Context
		postID       int
		mockReturn   []*model.Comment
		mockError    error
		wantComments []*model.Comment
		wantError    error
	}{
		{
			name:   "Success",
			ctx:    context.Background(),
			postID: 1,
			mockReturn: []*model.Comment{
				{
					ID:      1,
					Content: "Test Content 1",
					PostID:  1,
					UserID:  1,
				},
				{
					ID:      2,
					Content: "Test Content 2",
					PostID:  1,
					UserID:  2,
				},
			},
			mockError: nil,
			wantComments: []*model.Comment{
				{
					ID:      1,
					Content: "Test Content 1",
					PostID:  1,
					UserID:  1,
				},
				{
					ID:      2,
					Content: "Test Content 2",
					PostID:  1,
					UserID:  2,
				},
			},
			wantError: nil,
		},
		{
			name:         "Post not found",
			ctx:          context.Background(),
			postID:       0,
			mockReturn:   nil,
			mockError:    storage.ErrNotFound,
			wantComments: nil,
			wantError:    fmt.Errorf("%s: %w", operation, ErrNotFound),
		},
		{
			name:         "Error",
			ctx:          nil,
			postID:       1,
			mockReturn:   nil,
			mockError:    err,
			wantComments: nil,
			wantError:    fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider.On("CommentsByPost", tt.ctx, tt.postID).Return(tt.mockReturn, tt.mockError)

			_, err := commentService.CommentsByPost(tt.ctx, tt.postID)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.mockReturn, tt.wantComments)
			mockProvider.AssertExpectations(t)
		})
	}
}
