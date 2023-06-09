package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markraiter/simple-blog/models"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{
		db: db,
	}
}

func (r *CommentPostgres) Create(postID int, comment models.Comment) (int, error) {
	var id int
	createCommentQuery := fmt.Sprintf("INSERT INTO %s (post_id, name, email, body) VALUES ($1, $2, $3, $4) RETURNING id", commentsTable)

	row := r.db.QueryRow(createCommentQuery, comment.PostID, comment.Name, comment.Email, comment.Body)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CommentPostgres) GetAll() ([]models.Comment, error) {
	var comments []models.Comment

	query := fmt.Sprintf("SELECT * FROM %s", commentsTable)
	err := r.db.Select(&comments, query)

	return comments, err
}

func (r *CommentPostgres) GetByID(commentID int) (models.Comment, error) {
	var comment models.Comment

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", commentsTable)
	err := r.db.Get(&comment, query, commentID)

	return comment, err
}

func (r *CommentPostgres) Update(commentID int, input models.UpdateCommentInput) error {
	// Проверяем валидность входных данных
	if err := input.Validate(); err != nil {
		return err
	}

	// Подготавливаем запрос на обновление
	updateCommentQuery := fmt.Sprintf("UPDATE %s SET", commentsTable)

	// Список значений для обновления
	values := make(map[string]interface{})
	if input.Name != nil {
		updateCommentQuery += " name = :name,"
		values["name"] = *input.Name
	}
	if input.Email != nil {
		updateCommentQuery += " email = :email,"
		values["email"] = *input.Email
	}
	if input.Body != nil {
		updateCommentQuery += " body = :body,"
		values["body"] = *input.Body
	}

	// Удаляем последнюю запятую из запроса
	updateCommentQuery = strings.TrimSuffix(updateCommentQuery, ",")

	// Добавляем условие для конкретного поста
	updateCommentQuery += " WHERE id = :id"

	// Добавляем commentID в список значений
	values["id"] = commentID

	// Выполняем запрос на обновление
	_, err := r.db.NamedExec(updateCommentQuery, values)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommentPostgres) Delete(commentID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", commentsTable)

	_, err := r.db.Exec(query, commentID)

	return err
}
