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
type MockPostSaver struct{ mock.Mock }

func (m *MockPostSaver) SavePost(ctx context.Context, post *model.Post) (int, error) {
	args := m.Called(ctx, post)
	return args.Int(0), args.Error(1)
}

type MockPostProvider struct{ mock.Mock }

func (m *MockPostProvider) Post(ctx context.Context, id int) (*model.Post, error) {
	args := m.Called(ctx, id)
	post := args.Get(0)
	if post == nil {
		return nil, args.Error(1)
	}
	return post.(*model.Post), args.Error(1)
}

func (m *MockPostProvider) Posts(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

type MockPostProcessor struct{ mock.Mock }

func (m *MockPostProcessor) UpdatePost(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostProcessor) DeletePost(ctx context.Context, postID, userID int) error {
	args := m.Called(ctx, postID, userID)
	return args.Error(0)
}

// Tests
func TestPostService_SavePost(t *testing.T) {
	const operation = "service.SavePost"
	var err = errors.New("error")

	mockSaver := new(MockPostSaver)
	postService := &PostService{saver: mockSaver}

	tests := []struct {
		name       string
		postReq    *model.PostRequest
		userID     int
		mockReturn int
		mockError  error
		wantError  error
	}{
		{
			name: "Success",
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			userID:     1,
			mockReturn: 1,
			mockError:  nil,
			wantError:  nil,
		},
		{
			name: "Error",
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			userID:     0,
			mockReturn: 0,
			mockError:  err,
			wantError:  fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver.On("SavePost", mock.Anything, &model.Post{
				Title:   tt.postReq.Title,
				Content: tt.postReq.Content,
				UserID:  tt.userID,
			}).Return(tt.mockReturn, tt.mockError)

			_, err := postService.SavePost(context.Background(), tt.userID, tt.postReq)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockSaver.AssertExpectations(t)
		})
	}
}

func TestPostService_Post(t *testing.T) {
	const operation = "service.Post"
	var err = errors.New("error")

	mockProvider := new(MockPostProvider)
	postService := &PostService{provider: mockProvider}

	tests := []struct {
		name         string
		id           int
		mockReturn   *model.Post
		mockError    error
		expectedPost *model.Post
		expectedErr  error
	}{
		{
			name: "Success",
			id:   1,
			mockReturn: &model.Post{
				ID:      1,
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockError: nil,
			expectedPost: &model.Post{
				ID:      1,
				Title:   "Test Title",
				Content: "Test Content",
			},
			expectedErr: nil,
		},
		{
			name:         "Post not found",
			id:           2,
			mockReturn:   nil,
			mockError:    storage.ErrNotFound,
			expectedPost: nil,
			expectedErr:  fmt.Errorf("%s: %w", operation, ErrNotFound),
		},
		{
			name:         "Error",
			id:           3,
			mockReturn:   nil,
			mockError:    err,
			expectedPost: nil,
			expectedErr:  fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider.On("Post", mock.Anything, tt.id).Return(tt.mockReturn, tt.mockError)

			post, err := postService.Post(context.Background(), tt.id)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPost, post)
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestPostService_Posts(t *testing.T) {
	const operation = "service.Posts"
	var err = errors.New("error")

	mockProvider := new(MockPostProvider)
	postService := &PostService{provider: mockProvider}

	tests := []struct {
		name          string
		ctx           context.Context
		mockReturn    []*model.Post
		mockError     error
		expectedPosts []*model.Post
		expectedErr   error
	}{
		{
			name: "Success",
			ctx:  context.Background(),
			mockReturn: []*model.Post{
				{
					ID:      1,
					Title:   "Test Title",
					Content: "Test Content",
				},
				{
					ID:      2,
					Title:   "Test Title 2",
					Content: "Test Content 2",
				},
			},
			mockError: nil,
			expectedPosts: []*model.Post{
				{
					ID:      1,
					Title:   "Test Title",
					Content: "Test Content",
				},
				{
					ID:      2,
					Title:   "Test Title 2",
					Content: "Test Content 2",
				},
			},
			expectedErr: nil,
		},
		{
			name:          "Error",
			ctx:           nil,
			mockReturn:    nil,
			mockError:     err,
			expectedPosts: nil,
			expectedErr:   fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider.On("Posts", tt.ctx).Return(tt.mockReturn, tt.mockError)

			posts, err := postService.Posts(tt.ctx)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPosts, posts)
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	const operation = "service.UpdatePost"
	var err = errors.New("error")

	mockProcessor := new(MockPostProcessor)
	postService := &PostService{processor: mockProcessor}

	tests := []struct {
		name        string
		ctx         context.Context
		postID      int
		userID      int
		postReq     *model.PostRequest
		mockError   error
		expectedErr error
	}{
		{
			name:   "Success",
			ctx:    context.Background(),
			postID: 1,
			userID: 1,
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:   "Post not found",
			ctx:    context.Background(),
			postID: 2,
			userID: 1,
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockError:   storage.ErrNotFound,
			expectedErr: fmt.Errorf("%s: %w", operation, ErrNotFound),
		},
		{
			name:   "Not allowed",
			ctx:    context.Background(),
			postID: 1,
			userID: 2,
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockError:   storage.ErrNotAllowed,
			expectedErr: fmt.Errorf("%s: %w", operation, ErrNotAllowed),
		},
		{
			name:   "Error",
			ctx:    nil,
			postID: 1,
			userID: 1,
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockError:   err,
			expectedErr: fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor.On("UpdatePost", tt.ctx, mock.MatchedBy(func(post *model.Post) bool {
				return post.ID == tt.postID && post.UserID == tt.userID
			})).Return(tt.mockError)

			err := postService.UpdatePost(tt.ctx, tt.postID, tt.userID, tt.postReq)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockProcessor.AssertExpectations(t)
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	const operation = "service.DeletePost"
	var err = errors.New("error")

	mockProcessor := new(MockPostProcessor)
	postService := &PostService{processor: mockProcessor}

	tests := []struct {
		name        string
		ctx         context.Context
		postID      int
		userID      int
		mockError   error
		expectedErr error
	}{
		{
			name:        "Success",
			ctx:         context.Background(),
			postID:      1,
			userID:      1,
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Post not found",
			ctx:         context.Background(),
			postID:      2,
			userID:      1,
			mockError:   storage.ErrNotFound,
			expectedErr: fmt.Errorf("%s: %w", operation, ErrNotFound),
		},
		{
			name:        "Not allowed",
			ctx:         context.Background(),
			postID:      1,
			userID:      2,
			mockError:   storage.ErrNotAllowed,
			expectedErr: fmt.Errorf("%s: %w", operation, ErrNotAllowed),
		},
		{
			name:        "Error",
			ctx:         nil,
			postID:      1,
			userID:      1,
			mockError:   err,
			expectedErr: fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor.On("DeletePost", tt.ctx, tt.postID, tt.userID).Return(tt.mockError)

			err := postService.DeletePost(tt.ctx, tt.postID, tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockProcessor.AssertExpectations(t)
		})
	}
}
