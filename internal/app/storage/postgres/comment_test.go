package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
		{
			name: "Database error",
			comment: &model.Comment{
				Content: "content",
				PostID:  1,
				UserID:  1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs("content", 1, 1).
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
