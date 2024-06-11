package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	st "github.com/markraiter/simple-blog/internal/app/storage"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
)

func prepareStorage(t *testing.T) (*Storage, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	storage := &Storage{PostgresDB: db}

	closeFunc := func() {
		db.Close()
	}

	return storage, mock, closeFunc
}

func TestPostStorage_SavePost(t *testing.T) {
	const operation = "storage.SavePost"
	var err = errors.New("error")

	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name       string
		ctx        context.Context
		post       *model.Post
		mock       func()
		mockReturn int
		mockError  error
		wantID     int
		wantErr    error
	}{
		{
			name: "Success",
			ctx:  context.Background(),
			post: &model.Post{
				Title:   "Test Title",
				Content: "Test Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs("Test Title", "Test Content", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			mockReturn: 1,
			wantID:     1,
			wantErr:    nil,
		},
		{
			name: "Null value for userID",
			ctx:  context.Background(),
			post: &model.Post{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs("Test Title", "Test Content", 0).
					WillReturnError(sql.ErrNoRows)
			},
			mockReturn: 0,
			mockError:  sql.ErrNoRows,
			wantID:     0,
			wantErr:    fmt.Errorf("%s: %w", operation, sql.ErrNoRows),
		},
		{
			name: "Null value for title",
			ctx:  context.Background(),
			post: &model.Post{
				Content: "Test Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs("", "Test Content", 1).
					WillReturnError(sql.ErrNoRows)
			},
			mockReturn: 0,
			mockError:  sql.ErrNoRows,
			wantID:     0,
			wantErr:    fmt.Errorf("%s: %w", operation, sql.ErrNoRows),
		},
		{
			name: "Null value for content",
			ctx:  context.Background(),
			post: &model.Post{
				Title:  "Test Title",
				UserID: 1,
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs("Test Title", "", 1).
					WillReturnError(sql.ErrNoRows)
			},
			mockReturn: 0,
			mockError:  sql.ErrNoRows,
			wantID:     0,
			wantErr:    fmt.Errorf("%s: %w", operation, sql.ErrNoRows),
		},
		{
			name: "Error",
			ctx:  nil,
			post: &model.Post{
				Title:   "Test Title",
				Content: "Test Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").
					WithArgs("Test Title", "Test Content", 1).
					WillReturnError(err)
			},
			mockReturn: 0,
			mockError:  err,
			wantID:     0,
			wantErr:    fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			_, err := storage.SavePost(context.Background(), tt.post)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.mockReturn, tt.wantID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostStorage_Post(t *testing.T) {
	const operation = "storage.Post"
	var err = errors.New("error")

	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name       string
		postID     int
		mock       func()
		mockReturn *model.Post
		mockErr    error
		wantPost   *model.Post
		wantErr    error
	}{
		{
			name:   "Success",
			postID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = \\$1").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_count"}).
						AddRow(1, "Test Title", "Test Content", 1, 0))
			},
			mockReturn: &model.Post{
				ID:            1,
				Title:         "Test Title",
				Content:       "Test Content",
				UserID:        1,
				CommentsCount: 0,
			},
			mockErr: nil,
			wantPost: &model.Post{
				ID:            1,
				Title:         "Test Title",
				Content:       "Test Content",
				UserID:        1,
				CommentsCount: 0,
			},
			wantErr: nil,
		},
		{
			name:   "Post not found",
			postID: 2,
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = \\$1").
					ExpectQuery().
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			mockReturn: nil,
			mockErr:    sql.ErrNoRows,
			wantPost:   nil,
			wantErr:    fmt.Errorf("%s: %w", operation, st.ErrNotFound),
		},
		{
			name:   "Error",
			postID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = \\$1").
					ExpectQuery().
					WithArgs(1).
					WillReturnError(err)
			},
			mockReturn: nil,
			mockErr:    err,
			wantPost:   nil,
			wantErr:    fmt.Errorf("%s: %w", operation, err),
		},
		{
			name:   "Prepare error",
			postID: 1,
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = \\$1").
					WillReturnError(err)
			},
			mockReturn: nil,
			mockErr:    err,
			wantPost:   nil,
			wantErr:    fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := storage.Post(context.Background(), tt.postID)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.mockReturn, tt.wantPost)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostStorage_Posts(t *testing.T) {
	const operation = "storage.Posts"
	var (
		err     = errors.New("error")
		scanErr = errors.New("sql: Scan error")
	)

	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name       string
		mock       func()
		mockReturn []*model.Post
		mockErr    error
		wantPosts  []*model.Post
		wantErr    error
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts ORDER BY created_at DESC").
					ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_count"}).
						AddRow(1, "Test Title 1", "Test Content 1", 1, 0).
						AddRow(2, "Test Title 2", "Test Content 2", 2, 0))
			},
			mockReturn: []*model.Post{
				{
					ID:            1,
					Title:         "Test Title 1",
					Content:       "Test Content 1",
					UserID:        1,
					CommentsCount: 0,
				},
				{
					ID:            2,
					Title:         "Test Title 2",
					Content:       "Test Content 2",
					UserID:        2,
					CommentsCount: 0,
				},
			},
			mockErr: nil,
			wantPosts: []*model.Post{
				{
					ID:            1,
					Title:         "Test Title 1",
					Content:       "Test Content 1",
					UserID:        1,
					CommentsCount: 0,
				},
				{
					ID:            2,
					Title:         "Test Title 2",
					Content:       "Test Content 2",
					UserID:        2,
					CommentsCount: 0,
				},
			},
			wantErr: nil,
		},
		{
			name: "No posts",
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts ORDER BY created_at DESC").
					ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_count"}))
			},
			mockReturn: []*model.Post{},
			mockErr:    nil,
			wantPosts:  []*model.Post{},
			wantErr:    nil,
		},
		{
			name: "Error",
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts ORDER BY created_at DESC").
					ExpectQuery().
					WillReturnError(err)
			},
			mockReturn: nil,
			mockErr:    err,
			wantPosts:  nil,
			wantErr:    fmt.Errorf("%s: %w", operation, err),
		},
		{
			name: "Prepare error",
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts ORDER BY created_at DESC").
					WillReturnError(err)
			},
			mockReturn: nil,
			mockErr:    err,
			wantPosts:  nil,
			wantErr:    fmt.Errorf("%s: %w", operation, err),
		},
		{
			name: "Scan error",
			mock: func() {
				mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts ORDER BY created_at DESC").
					ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_count"}).
						AddRow("invalid_id", "Test Title", "Test Content", 1, 0))
			},
			mockReturn: nil,
			mockErr:    scanErr,
			wantPosts:  nil,
			wantErr:    fmt.Errorf("%s: %w", operation, scanErr),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := storage.Posts(context.Background())

			if tt.wantErr != nil && tt.mockErr != scanErr {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else if tt.mockErr == scanErr {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.mockReturn, tt.wantPosts)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostStorage_UpdatePost(t *testing.T) {
	const operation = "storage.UpdatePost"
	var err = errors.New("error")

	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name    string
		post    *model.Post
		mock    func()
		wantErr error
	}{
		{
			name: "Success",
			post: &model.Post{
				ID:      1,
				Title:   "Updated Title",
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE posts SET title = \\$1, content = \\$2 WHERE id = \\$3 AND user_id = \\$4 RETURNING id").
					WithArgs("Updated Title", "Updated Content", 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: nil,
		},
		{
			name: "Post Not Found",
			post: &model.Post{
				ID:      2,
				Title:   "Title",
				Content: "Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE posts SET title = \\$1, content = \\$2 WHERE id = \\$3 AND user_id = \\$4 RETURNING id").
					WithArgs("Title", "Content", 2, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM posts WHERE id = \\$1").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: fmt.Errorf("%s: %w", operation, st.ErrNotFound),
		},
		{
			name: "User Not Allowed",
			post: &model.Post{
				ID:      3,
				Title:   "Another Title",
				Content: "Another Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE posts SET title = \\$1, content = \\$2 WHERE id = \\$3 AND user_id = \\$4 RETURNING id").
					WithArgs("Another Title", "Another Content", 3, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM posts WHERE id = \\$1").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
			},
			wantErr: fmt.Errorf("%s: %w", operation, st.ErrNotAllowed),
		},
		{
			name: "Error",
			post: &model.Post{
				ID:      1,
				Title:   "Updated Title",
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE posts SET title = \\$1, content = \\$2 WHERE id = \\$3 AND user_id = \\$4 RETURNING id").
					WithArgs("Updated Title", "Updated Content", 1, 1).
					WillReturnError(err)
			},
			wantErr: fmt.Errorf("%s: %w", operation, err),
		},
		{
			name: "QueryRawContext error",
			post: &model.Post{
				ID:      1,
				Title:   "Updated Title",
				Content: "Updated Content",
				UserID:  1,
			},
			mock: func() {
				mock.ExpectQuery("UPDATE posts SET title = \\$1, content = \\$2 WHERE id = \\$3 AND user_id = \\$4 RETURNING id").
					WithArgs("Updated Title", "Updated Content", 1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM posts WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(err)
			},
			wantErr: fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ctx := context.Background()
			err := storage.UpdatePost(ctx, tt.post)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostStorage_DeletePost(t *testing.T) {
	const operation = "storage.DeletePost"
	var err = errors.New("error")

	storage, mock, closeDB := prepareStorage(t)
	defer closeDB()

	tests := []struct {
		name    string
		postID  int
		userID  int
		mock    func()
		wantErr error
	}{
		{
			name:   "Success",
			postID: 1,
			userID: 1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM posts WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: nil,
		},
		{
			name:   "Post Not Found",
			postID: 2,
			userID: 1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM posts WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(2, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM posts WHERE id = \\$1").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: fmt.Errorf("%s: %w", operation, st.ErrNotFound),
		},
		{
			name:   "User Not Allowed",
			postID: 3,
			userID: 1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM posts WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(3, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM posts WHERE id = \\$1").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
			},
			wantErr: fmt.Errorf("%s: %w", operation, st.ErrNotAllowed),
		},
		{
			name:   "Error",
			postID: 1,
			userID: 1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM posts WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(1, 1).
					WillReturnError(err)
			},
			wantErr: fmt.Errorf("%s: %w", operation, err),
		},
		{
			name:   "QueryRawContext error",
			postID: 1,
			userID: 1,
			mock: func() {
				mock.ExpectQuery("DELETE FROM posts WHERE id = \\$1 AND user_id = \\$2 RETURNING id").
					WithArgs(1, 1).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("SELECT id FROM posts WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(err)
			},
			wantErr: fmt.Errorf("%s: %w", operation, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ctx := context.Background()
			err := storage.DeletePost(ctx, tt.postID, tt.userID)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
