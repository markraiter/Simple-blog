package postgres

import (
	"context"
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
            mock: func() {},
            wantID: 0,
            wantErr: true,
        },
        {
            name: "Null value for title",
            post: &model.Post{
                Content: "Test Content",
                UserID:  1,
            },
            mock: func() {},
            wantID: 0,
            wantErr: true,
        },
        {
            name: "Null value for content",
            post: &model.Post{
                Title:  "Test Title",
                UserID: 1,
            },
            mock: func() {},
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
        })
    }
}
