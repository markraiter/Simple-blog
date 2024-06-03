package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestSavePost(t *testing.T) {
    storage, mock, closeDB := prepareStorage(t)
    defer closeDB()

    tests := []struct {
        name    string
        post    *model.Post
        mock    func()
        wantID  int
        wantErr bool
    }{
        {
            name: "Success",
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
            wantID: 1,
            wantErr: false,
        },
        {
            name: "Null value for userID",
            post: &model.Post{
                Title:   "Test Title",
                Content: "Test Content",
            },
            mock: func() {
                mock.ExpectQuery("INSERT INTO posts").
                    WithArgs("Test Title", "Test Content", 0).
                    WillReturnError(sql.ErrNoRows)
            },
            wantID: 0,
            wantErr: true,
        },
        {
            name: "Null value for title",
            post: &model.Post{
                Content: "Test Content",
                UserID:  1,
            },
            mock: func() {
                mock.ExpectQuery("INSERT INTO posts").
                    WithArgs("", "Test Content", 1).
                    WillReturnError(sql.ErrNoRows)
            },
            wantID: 0,
            wantErr: true,
        },
        {
            name: "Null value for content",
            post: &model.Post{
                Title:  "Test Title",
                UserID: 1,
            },
            mock: func() {
                mock.ExpectQuery("INSERT INTO posts").
                    WithArgs("Test Title", "", 1).
                    WillReturnError(sql.ErrNoRows)
            },
            wantID: 0,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            postID, err := storage.SavePost(context.Background(), tt.post)

            assert.Equal(t, tt.wantID, postID)
            assert.Equal(t, tt.wantErr, err != nil)

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("there were unfulfilled expectations: %s", err)
            }
        })
    }
}

func TestGetPost(t *testing.T) {
    storage, mock, closeDB := prepareStorage(t)
    defer closeDB()

    tests := []struct {
        name    string
        postID  int
        mock    func()
        wantPost *model.Post
        wantErr  bool
    }{
        {
            name: "Success",
            postID: 1,
            mock: func() {
                mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = \\$1").
                    ExpectQuery().
                    WithArgs(1).
                    WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_count"}).
                        AddRow(1, "Test Title", "Test Content", 1, 0))
            },
            wantPost: &model.Post{
                ID:      1,
                Title:   "Test Title",
                Content: "Test Content",
                UserID:  1,
                CommentsCount: 0,
            },
            wantErr: false,
        },
        {
            name: "Post not found",
            postID: 999999,
            mock: func() {
                mock.ExpectPrepare("SELECT id, title, content, user_id, comments_count FROM posts WHERE id = \\$1").
                    ExpectQuery().
                    WithArgs(999999).
                    WillReturnError(sql.ErrNoRows)
            },
            wantPost: nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()
            ctx := context.Background()
            post, err := storage.Post(ctx, tt.postID)

            assert.Equal(t, tt.wantPost, post)
            assert.Equal(t, tt.wantErr, err != nil)

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("there were unfulfilled expectations: %s", err)
            }
        })
    }
}

























