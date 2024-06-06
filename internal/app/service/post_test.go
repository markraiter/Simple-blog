package service

import (
	"context"
	"errors"
	"testing"

	"github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockPostSaver struct {
	mock.Mock
}

func (m *MockPostSaver) SavePost(ctx context.Context, post *model.Post) (int, error) {
	args := m.Called(ctx, post)
	return args.Int(0), args.Error(1)
}

type MockPostProvider struct {
    mock.Mock
}

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

// Tests
func TestPostService_SavePost(t *testing.T) {
	mockSaver := new(MockPostSaver)
	postService := &PostService{
		saver: mockSaver,
	}

	tests := []struct {
		name        string
		postReq     *model.PostRequest
		mockReturn  int
		mockError   error
		expectedID  int
		expectedErr error
	}{
		{
			name: "Success",
			postReq: &model.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockReturn:  1,
			mockError:   nil,
			expectedID:  1,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver.On("SavePost", mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockError)

			id, err := postService.SavePost(context.Background(), 1, tt.postReq)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedID, id)
			mockSaver.AssertExpectations(t)
		})
	}
}

func TestPostService_Post(t *testing.T) {
	mockProvider := new(MockPostProvider)
	postService := &PostService{
		provider: mockProvider,
	}

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
			mockError:   nil,
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
			expectedErr:  ErrNotFound, 
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider.On("Post", mock.Anything, tt.id).Return(tt.mockReturn, tt.mockError)

			post, err := postService.Post(context.Background(), tt.id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPost, post)
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestPostService_Posts(t *testing.T) {
    mockProvider := new(MockPostProvider)
    postService := &PostService{
        provider: mockProvider,
    }

    tests := []struct {
        name        string
        mockReturn  []*model.Post
        mockError   error
        expectedPosts []*model.Post
        expectedErr error
    }{
        {
            name: "Success",
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
            mockError:   nil,
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
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockProvider.On("Posts", mock.Anything).Return(tt.mockReturn, tt.mockError)

            posts, err := postService.Posts(context.Background())

            if tt.expectedErr != nil {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedErr.Error())
            } else {
                assert.NoError(t, err)
            }

            assert.Equal(t, tt.expectedPosts, posts)
            mockProvider.AssertExpectations(t)
        })
    }
}
