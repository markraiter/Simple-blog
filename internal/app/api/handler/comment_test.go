package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/markraiter/simple-blog/internal/app/service"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockCommentSaver struct{ mock.Mock }

func (m *MockCommentSaver) SaveComment(ctx context.Context, userID int, commentReq *model.CommentRequest) (int, error) {
	args := m.Called(ctx, userID, commentReq)
	return args.Int(0), args.Error(1)
}

// Tests
func TestCommentHandler_CreateComment(t *testing.T) {
	mockSaver := new(MockCommentSaver)
	h := &CommentHandler{
		log:       log,
		validate:  validator.New(),
		saver:     mockSaver,
		provider:  nil,
		processor: nil,
	}

	tests := []struct {
		name              string
		commentReq        *model.CommentRequest
		mockReturnID      int
		mockReturnErr     error
		expectedStatus    int
		expectSaveComment bool
	}{
		{
			name: "Success",
			commentReq: &model.CommentRequest{
				Content: "content",
				PostID:  1,
			},
			mockReturnID:      1,
			expectedStatus:    http.StatusCreated,
			expectSaveComment: true,
		},
		{
			name:              "Invalid request - JSON parsing error",
			commentReq:        nil,
			mockReturnID:      0,
			expectedStatus:    http.StatusBadRequest,
			expectSaveComment: false,
		},
		{
			name: "Invalid request - Validation error",
			commentReq: &model.CommentRequest{
				Content: "",
				PostID:  1,
			},
			mockReturnID:      0,
			expectedStatus:    http.StatusBadRequest,
			expectSaveComment: false,
		},
		{
			name: "Internal server error",
			commentReq: &model.CommentRequest{
				Content: "content",
				PostID:  1,
			},
			mockReturnID:      0,
			mockReturnErr:     assert.AnError,
			expectedStatus:    http.StatusInternalServerError,
			expectSaveComment: true,
		},
		{
			name: "Post does not exist",
			commentReq: &model.CommentRequest{
				Content: "content",
				PostID:  1,
			},
			mockReturnID:      0,
			mockReturnErr:     service.ErrPostNotExists,
			expectedStatus:    http.StatusBadRequest,
			expectSaveComment: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte

			if tt.commentReq != nil {
				body, _ = json.Marshal(tt.commentReq)
			} else {
				body = []byte("{")
			}

			if tt.expectSaveComment {
				mockSaver.On("SaveComment", mock.Anything, mock.Anything, tt.commentReq).Return(tt.mockReturnID, tt.mockReturnErr).Once()
			}

			r := httptest.NewRequest(http.MethodPost, "/api/comments", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler := h.CreateComment(context.Background())
			handler.ServeHTTP(w, r)

			resp := w.Result()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSaver.AssertExpectations(t)
		})
	}
}
