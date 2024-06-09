package model

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content" validate:"required,min=3" example:"lorem ipsum dolor sit amet ..."`
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
}

type CommentRequest struct {
	Content string `json:"content" validate:"required,min=3" example:"lorem ipsum dolor sit amet ..."`
	PostID  int    `json:"post_id" validate:"required" example:"1"`
}
