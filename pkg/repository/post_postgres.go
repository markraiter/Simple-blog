package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/markraiter/simple-blog/models"
)

type PostPostgres struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{
		db: db,
	}
}

func (r *PostPostgres) Create(userID int, post models.Post) (int, error) {
	var id int
	createPostQuery := fmt.Sprintf("INSERT INTO %s (user_id, title, body) VALUES ($1, $2, $3) RETURNING id", postsTable)

	row := r.db.QueryRow(createPostQuery, post.UserID, post.Title, post.Body)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostPostgres) GetAll() ([]models.Post, error) {
	var posts []models.Post

	query := fmt.Sprintf("SELECT * FROM %s", postsTable)
	err := r.db.Select(&posts, query)

	return posts, err
}

func (r *PostPostgres) GetByID(postID int) (models.Post, error) {
	var post models.Post

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postsTable)
	err := r.db.Get(&post, query, postID)

	return post, err
}

func (r *PostPostgres) Update(postID int, input models.UpdatePostInput) error {
	// Проверяем валидность входных данных
	if err := input.Validate(); err != nil {
		return err
	}

	// Подготавливаем запрос на обновление
	updatePostQuery := fmt.Sprintf("UPDATE %s SET", postsTable)

	// Список значений для обновления
	values := make(map[string]interface{})
	if input.Title != nil {
		updatePostQuery += " title = :title,"
		values["title"] = *input.Title
	}
	if input.Body != nil {
		updatePostQuery += " body = :body,"
		values["body"] = *input.Body
	}

	// Удаляем последнюю запятую из запроса
	updatePostQuery = strings.TrimSuffix(updatePostQuery, ",")

	// Добавляем условие для конкретного поста
	updatePostQuery += " WHERE id = :id"

	// Добавляем postID в список значений
	values["id"] = postID

	// Выполняем запрос на обновление
	_, err := r.db.NamedExec(updatePostQuery, values)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostPostgres) Delete(postID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postsTable)

	_, err := r.db.Exec(query, postID)

	return err
}
