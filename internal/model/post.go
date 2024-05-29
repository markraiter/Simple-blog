package model

type Post struct {
	ID            int    `json:"id"`
	Title         string `json:"title" validate:"required,min=3,max=50" example:"title"`
	Content       string `json:"content" validate:"required,min=3" example:"lorem ipsum dolor sit amet ..."`
	UserID        int    `json:"user_id"`
	CommentsCount int    `json:"comments_count" example:"0"`
}

type PostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=50" example:"title"`
	Content string `json:"content" validate:"required,min=3" example:"lorem ipsum dolor sit amet ..."`
}
