package service

import (
	"context"
	"errors"
	"testing"

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

// Tests
func TestCommentService_SaveComment(t *testing.T) {
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
			mockError:  errors.New("error"),
			wantError:  errors.New("service.SaveComment: error"),
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
