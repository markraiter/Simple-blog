package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	// "fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/internal/app/api/middleware"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostSaver struct {
	mock.Mock
}

func (m *MockPostSaver) SavePost(ctx context.Context, userID int, postReq *model.PostRequest) (int, error) {
	args := m.Called(ctx, userID, postReq)
	return args.Int(0), args.Error(1)
}

func TestCreatePost(t *testing.T) {
	mockSaver := new(MockPostSaver)
	ph := &PostHandler{
        log:       slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
        validate:  validator.New(),
		saver:     mockSaver,
		provider:  nil,
		processor: nil,
	}

	tests := []struct {
		name           string
		userID         string
        addUserIDToCtx bool
		postReq        *model.PostRequest
		mockReturnID   int
		mockReturnErr  error
		expectedStatus int
        expectSavePost bool
	}{
		{
			name: "Success",
			userID: "1",
            addUserIDToCtx: true,
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
            name: "Error getting userID from ctx",
            userID: "",
            addUserIDToCtx: false,
            postReq: &model.PostRequest{
                Title:   "Title",
                Content: "Content",
            },
            mockReturnID:   0,
            mockReturnErr:  nil,
            expectedStatus: http.StatusInternalServerError,
            expectSavePost: false,
        },
        {
            name: "Error parsing userID",
            userID: "a",
            addUserIDToCtx: true,
            postReq: &model.PostRequest{
                Title:   "Title",
                Content: "Content",
            },
            mockReturnID:   0,
            mockReturnErr:  nil,
            expectedStatus: http.StatusInternalServerError,
            expectSavePost: false,
        },
		{
			name: "Invalid request - JSON parsing error",
			userID: "1",
            addUserIDToCtx: true,
			postReq: nil,
			mockReturnID:  0,
			mockReturnErr: nil,
			expectedStatus: http.StatusBadRequest,
            expectSavePost: false,
		},
		{
			name: "Invalid request - validation error",
			userID: "1",
            addUserIDToCtx: true,
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
			userID: "1",
            addUserIDToCtx: true,
			postReq: &model.PostRequest{
				Title:   "Title",
				Content: "Content",
			},
			mockReturnID:   0,
			mockReturnErr:  fmt.Errorf("internal server error"),
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
			if tt.addUserIDToCtx {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UIDKey, tt.userID))
			}
			w := httptest.NewRecorder()

			if tt.expectSavePost {
				mockSaver.On("SavePost", mock.Anything, mock.Anything, tt.postReq).Return(tt.mockReturnID, tt.mockReturnErr).Once()
			}

			handler := ph.CreatePost(context.Background())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSaver.AssertExpectations(t)
		})
	}
}
