package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	"github.com/markraiter/simple-blog/internal/app/service"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

// mocks
type MockPostSaver struct{ mock.Mock }

func (m *MockPostSaver) SavePost(ctx context.Context, userID int, postReq *model.PostRequest) (int, error) {
	args := m.Called(ctx, userID, postReq)
	return args.Int(0), args.Error(1)
}

type MockPostProvider struct{ mock.Mock }

func (m *MockPostProvider) Post(ctx context.Context, postID int) (*model.Post, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Post), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockPostProvider) Posts(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

type MockPostProcessor struct{ mock.Mock }

func (m *MockPostProcessor) UpdatePost(ctx context.Context, posID, userID int, postReq *model.PostRequest) error {
	args := m.Called(ctx, posID, userID, postReq)
	return args.Error(0)
}

func (m *MockPostProcessor) DeletePost(ctx context.Context, posID, userID int) error {
	args := m.Called(ctx, posID, userID)
	return args.Error(0)
}

// tests
func TestPostHandler_CreatePost(t *testing.T) {
	mockSaver := new(MockPostSaver)
	h := &PostHandler{
		log:       log,
		validate:  validator.New(),
		saver:     mockSaver,
		provider:  nil,
		processor: nil,
	}

	tests := []struct {
		name           string
		postReq        *model.PostRequest
		mockReturnID   int
		mockReturnErr  error
		expectedStatus int
		expectSavePost bool
	}{
		{
			name: "Success",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnID:   1,
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectSavePost: true,
		},
		{
			name:           "Invalid request - JSON parsing error",
			postReq:        nil,
			mockReturnID:   0,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectSavePost: false,
		},
		{
			name: "Invalid request - validation error",
			postReq: &model.PostRequest{
				Title:   "",
				Content: "",
			},
			mockReturnID:   0,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectSavePost: false,
		},
		{
			name: "Internal server error",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnID:   0,
			mockReturnErr:  assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectSavePost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.postReq != nil {
				body, _ = json.Marshal(tt.postReq)
			} else {
				body = []byte("{")
			}

			req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			if tt.expectSavePost {
				mockSaver.On("SavePost", mock.Anything, mock.Anything, tt.postReq).Return(tt.mockReturnID, tt.mockReturnErr).Once()
			}

			handler := h.CreatePost(context.Background())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSaver.AssertExpectations(t)
		})
	}
}

func TestPostHandler_Post(t *testing.T) {
	mockProvider := new(MockPostProvider)
	h := &PostHandler{
		log:       log,
		validate:  validator.New(),
		saver:     nil,
		provider:  mockProvider,
		processor: nil,
	}

	tests := []struct {
		name           string
		postID         string
		mockReturnPost *model.Post
		mockReturnErr  error
		expectedStatus int
		expectGetPost  bool
	}{
		{
			name:   "Success",
			postID: "1",
			mockReturnPost: &model.Post{
				ID:      1,
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectGetPost:  true,
		},
		{
			name:           "Error getting postID from query",
			postID:         "",
			mockReturnPost: nil,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectGetPost:  false,
		},
		{
			name:           "Error parsing postID",
			postID:         "a",
			mockReturnPost: nil,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectGetPost:  false,
		},
		{
			name:           "Post not found",
			postID:         "2",
			mockReturnPost: nil,
			mockReturnErr:  service.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectGetPost:  true,
		},
		{
			name:           "Internal server error",
			postID:         "1",
			mockReturnPost: nil,
			mockReturnErr:  fmt.Errorf("internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectGetPost:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/posts?id="+tt.postID, nil)
			w := httptest.NewRecorder()

			if tt.expectGetPost {
				postID, _ := strconv.Atoi(tt.postID)
				mockProvider.On("Post", mock.Anything, postID).Return(tt.mockReturnPost, tt.mockReturnErr).Once()
			}

			handler := h.Post(context.Background())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestPostHandler_Posts(t *testing.T) {
	mockProvider := new(MockPostProvider)
	h := &PostHandler{
		log:       log,
		validate:  validator.New(),
		saver:     nil,
		provider:  mockProvider,
		processor: nil,
	}

	tests := []struct {
		name            string
		mockReturnPosts []*model.Post
		mockReturnErr   error
		expectedStatus  int
		expectGetPosts  bool
	}{
		{
			name: "Success",
			mockReturnPosts: []*model.Post{
				{
					ID:      1,
					Title:   "Title",
					Content: "Content",
				},
				{
					ID:      2,
					Title:   "Title",
					Content: "Content",
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectGetPosts: true,
		},
		{
			name:            "Internal server error",
			mockReturnPosts: nil,
			mockReturnErr:   fmt.Errorf("internal server error"),
			expectedStatus:  http.StatusInternalServerError,
			expectGetPosts:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/posts", nil)
			w := httptest.NewRecorder()

			if tt.expectGetPosts {
				mockProvider.On("Posts", mock.Anything).Return(tt.mockReturnPosts, tt.mockReturnErr).Once()
			}

			handler := h.Posts(context.Background())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestPostHandler_UpdatePost(t *testing.T) {
	mockProcessor := new(MockPostProcessor)
	h := &PostHandler{
		log:       log,
		validate:  validator.New(),
		saver:     nil,
		provider:  nil,
		processor: mockProcessor,
	}

	tests := []struct {
		name             string
		userID           string
		postID           string
		postReq          *model.PostRequest
		mockReturnErr    error
		expectedStatus   int
		expectUpdatePost bool
	}{
		{
			name:   "Success",
			userID: "1",
			postID: "1",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:    nil,
			expectedStatus:   http.StatusOK,
			expectUpdatePost: true,
		},
		{
			name:             "Invalid request - JSON parsing error",
			userID:           "1",
			postID:           "1",
			postReq:          nil,
			mockReturnErr:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectUpdatePost: false,
		},
		{
			name:   "Invalid request - validation error",
			userID: "1",
			postID: "1",
			postReq: &model.PostRequest{
				Title:   "",
				Content: "",
			},
			mockReturnErr:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectUpdatePost: false,
		},
		{
			name:   "Error getting postID from query",
			userID: "1",
			postID: "",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error parsing postID",
			userID: "1",
			postID: "a",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "User is not allowed to perform this operation",
			userID: "1",
			postID: "1",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:    service.ErrNotAllowed,
			expectedStatus:   http.StatusForbidden,
			expectUpdatePost: true,
		},
		{
			name:   "Post not found",
			userID: "1",
			postID: "1",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:    service.ErrNotFound,
			expectedStatus:   http.StatusNotFound,
			expectUpdatePost: true,
		},
		{
			name:   "Internal server error",
			userID: "1",
			postID: "1",
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnErr:    fmt.Errorf("internal server error"),
			expectedStatus:   http.StatusInternalServerError,
			expectUpdatePost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.postReq != nil {
				body, _ = json.Marshal(tt.postReq)
			} else {
				body = []byte("{")
			}

			req := httptest.NewRequest("PUT", "/api/posts?id="+tt.postID, bytes.NewBuffer(body))
			req = req.WithContext(context.WithValue(req.Context(), middleware.UIDKey, tt.userID))
			w := httptest.NewRecorder()

			if tt.expectUpdatePost {
				postID, _ := strconv.Atoi(tt.postID)
				mockProcessor.On("UpdatePost", mock.Anything, postID, 1, tt.postReq).Return(tt.mockReturnErr).Once()
			}

			handler := h.UpdatePost(context.Background())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockProcessor.AssertExpectations(t)
		})
	}
}

func TestPostHandler_DeletePost(t *testing.T) {
	mockProcessor := new(MockPostProcessor)
	h := &PostHandler{
		log:       log,
		validate:  validator.New(),
		saver:     nil,
		provider:  nil,
		processor: mockProcessor,
	}

	tests := []struct {
		name             string
		userID           string
		postID           string
		mockReturnErr    error
		expectedStatus   int
		expectDeletePost bool
	}{
		{
			name:             "Success",
			userID:           "1",
			postID:           "1",
			mockReturnErr:    nil,
			expectedStatus:   http.StatusOK,
			expectDeletePost: true,
		},
		{
			name:           "Error getting postID from query",
			userID:         "1",
			postID:         "",
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error parsing postID",
			userID:         "1",
			postID:         "a",
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:             "User is not allowed to perform this operation",
			userID:           "1",
			postID:           "1",
			mockReturnErr:    service.ErrNotAllowed,
			expectedStatus:   http.StatusForbidden,
			expectDeletePost: true,
		},
		{
			name:             "Post not found",
			userID:           "1",
			postID:           "1",
			mockReturnErr:    service.ErrNotFound,
			expectedStatus:   http.StatusNotFound,
			expectDeletePost: true,
		},
		{
			name:             "Internal server error",
			userID:           "1",
			postID:           "1",
			mockReturnErr:    fmt.Errorf("internal server error"),
			expectedStatus:   http.StatusInternalServerError,
			expectDeletePost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/posts?id="+tt.postID, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.UIDKey, tt.userID))
			w := httptest.NewRecorder()

			if tt.expectDeletePost {
				postID, _ := strconv.Atoi(tt.postID)
				mockProcessor.On("DeletePost", mock.Anything, postID, 1).Return(tt.mockReturnErr).Once()
			}

			handler := h.DeletePost(context.Background())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockProcessor.AssertExpectations(t)
		})
	}
}
