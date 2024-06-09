package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	st "github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveComment(t *testing.T) {
	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name    string
		comment *model.Comment
		mock    func()
		wantID  int
		wantErr bool
		err     error
	}{
		{
			name: "Success",
			comment: &model.Comment{
				Content: "content",
				PostID:  1,
				UserID:  1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs("content", 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("UPDATE posts SET comments_count = comments_count \\+ 1 WHERE id = \\$1").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantID:  1,
			wantErr: false,
			err:     nil,
		},
		{
			name: "Null value for user_id",
			comment: &model.Comment{
				Content: "content",
				PostID:  1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs("content", 1, 0).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantID:  0,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
		{
			name: "Null value for post_id",
			comment: &model.Comment{
				Content: "content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs("content", 0, 1).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantID:  0,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
		{
			name: "Null value for content",
			comment: &model.Comment{
				PostID: 1,
				UserID: 1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs("", 1, 1).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantID:  0,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			commentID, err := storage.SaveComment(context.Background(), tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				if !errors.Is(err, tt.err) {
					t.Errorf("error = %v, wantErr %v", err, tt.err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, commentID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestComment(t *testing.T) {
	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name        string
		commentID   int
		mock        func()
		wantComment *model.Comment
		wantErr     bool
		err         error
	}{
		{
			name:      "Success",
			commentID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, content, post_id, user_id FROM comments WHERE id = \\$1").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "content", "post_id", "user_id"}).
						AddRow(1, "Test Content", 1, 1))
			},
			wantComment: &model.Comment{
				ID:      1,
				Content: "Test Content",
				PostID:  1,
				UserID:  1,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name:      "Comment not found",
			commentID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, content, post_id, user_id FROM comments WHERE id = \\$1").
					ExpectQuery().
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantComment: nil,
			wantErr:     true,
			err:         st.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ctx := context.Background()

			comment, err := storage.Comment(ctx, tt.commentID)

			if tt.wantErr {
				assert.Error(t, err)
				if !errors.Is(err, tt.err) {
					t.Errorf("error = %v, wantErr %v", err, tt.err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantComment, comment)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCommentsByPost(t *testing.T) {
	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name         string
		postID       int
		mock         func()
		wantComments []*model.Comment
		wantErr      bool
		err          error
	}{
		{
			name:   "Success",
			postID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, content, post_id, user_id FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "content", "post_id", "user_id"}).
						AddRow(1, "Test Content", 1, 1).
						AddRow(2, "Test Content 2", 1, 1))
			},
			wantComments: []*model.Comment{
				{
					ID:      1,
					Content: "Test Content",
					PostID:  1,
					UserID:  1,
				},
				{
					ID:      2,
					Content: "Test Content 2",
					PostID:  1,
					UserID:  1,
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name:   "No comments found",
			postID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, content, post_id, user_id FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "content", "post_id", "user_id"}))
			},
			wantComments: []*model.Comment{},
			wantErr:      false,
			err:          nil,
		},
		{
			name:   "No post found",
			postID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, content, post_id, user_id FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC").
					ExpectQuery().
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantComments: nil,
			wantErr:      true,
			err:          st.ErrNotFound,
		},
		{
			name:   "No postID",
			postID: 0,
			mock: func() {
				mock.ExpectPrepare("SELECT id, content, post_id, user_id FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC").
					ExpectQuery().
					WithArgs(0).
					WillReturnError(sql.ErrNoRows)
			},
			wantComments: nil,
			wantErr:      true,
			err:          st.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ctx := context.Background()

			comments, err := storage.CommentsByPost(ctx, tt.postID)

			if tt.wantErr {
				assert.Error(t, err)
				if !errors.Is(err, tt.err) {
					t.Errorf("error = %v, wantErr %v", err, tt.err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantComments, comments)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name    string
		comment *model.Comment
		mock    func()
		wantErr bool
		err     error
	}{
		{
			name: "Success",
			comment: &model.Comment{
				ID:      1,
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE comments SET content = \\$1 WHERE id = \\$2 AND user_id = \\$3 RETURNING id").
					WithArgs("Updated Content", 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Comment not found",
			comment: &model.Comment{
				ID:      1,
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE comments SET content = \\$1 WHERE id = \\$2 AND user_id = \\$3 RETURNING id").
					WithArgs("Updated Content", 1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM comments WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     st.ErrNotFound,
		},
		{
			name: "Post not found",
			comment: &model.Comment{
				ID:      1,
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE comments SET content = \\$1 WHERE id = \\$2 AND user_id = \\$3 RETURNING id").
					WithArgs("Updated Content", 1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM comments WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     st.ErrNotFound,
		},
		{
			name: "User not allowed",
			comment: &model.Comment{
				ID:      1,
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE comments SET content = \\$1 WHERE id = \\$2 AND user_id = \\$3 RETURNING id").
					WithArgs("Updated Content", 1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM comments WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: true,
			err:     st.ErrNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ctx := context.Background()

			err := storage.UpdateComment(ctx, tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				if !errors.Is(err, tt.err) {
					t.Errorf("error = %v, wantErr %v", err, tt.err)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name      string
		commentID int
		userID    int
		mock      func()
		wantErr   bool
		err       error
	}{
		{
			name:      "Success",
			commentID: 1,
			userID:    1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM comments WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name:      "Comment not found",
			commentID: 1,
			userID:    1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM comments WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM comments WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     st.ErrNotFound,
		},
		{
			name:      "User not allowed",
			commentID: 1,
			userID:    1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM comments WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM comments WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: true,
			err:     st.ErrNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ctx := context.Background()

			err := storage.DeleteComment(ctx, tt.commentID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if !errors.Is(err, tt.err) {
					t.Errorf("error = %v, wantErr %v", err, tt.err)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
