package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

// @Summary Get All Comments
// @Security ApiKeyAuth
// @Tags comments
// @Description get all comments
// @ID get-all-comments
// @Produce  json
// @Success 200 {array} models.Comment "Returns an array of comments."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments/all [get]
// getAllComments is a handler for getting list of all the comments. It returns an array of comments
func (h *Handler) getAllComments(c *gin.Context) {
	comments, err := h.storage.Comments.GetAll()
	if err != nil {
		log.Printf("error fetching comments: %v", err)
		c.String(http.StatusBadRequest, "error fetching comments")
		return
	}

	log.Printf("you successfully gor %d comments", len(comments))
	c.JSON(http.StatusOK, comments)
}

// @Summary Filter Comments by Post
// @Security ApiKeyAuth
// @Tags comments
// @Description filter comments by post
// @ID filter-comments-by-post
// @Param post_id query int true "Post ID"
// @Produce json
// @Success 200 {array} models.Comment "Returns an array of comments associated with the post."
// @Failure 404 {string} string "Comment with specified post_id does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments/post [get]
// filterCommentsByPost is a handler for filtering comments by post. It returns an array of comments associated with the post.
func (h *Handler) filterCommentsByPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Query("post_id"))
	if err != nil {
		log.Printf("no such post in the database: %v", err)
		c.String(http.StatusNotFound, "no such post in the database")
		return
	}

	comments, err := h.storage.Comments.FilterByPost(uint(postID))
	if err != nil {
		log.Printf("error fetching comments: %v", err)
		c.String(http.StatusInternalServerError, "error fetching comments")
		return
	}

	log.Printf("you successfully got %d comments", len(comments))
	c.JSON(http.StatusOK, comments)
}

// @Summary Filter Comments by User
// @Security ApiKeyAuth
// @Tags comments
// @Description filter comments by user
// @ID filter-comments-by-user
// @Param user_id query int true "User ID"
// @Produce json
// @Success 200 {array} models.Comment "Returns an array of comments created by the user."
// @Failure 404 {string} string "Comment with specified user_id does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments/user [get]
// filterCommentsByUser is a handler for filtering comments by user. It returns an array of comments created by the user.
func (h *Handler) filterCommentsByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		log.Printf("no such user in the database: %v", err)
		c.String(http.StatusNotFound, "no such user in the database")
		return
	}

	comments, err := h.storage.Comments.FilterByUser(uint(userID))
	if err != nil {
		log.Printf("error fetching comments: %v", err)
		c.String(http.StatusInternalServerError, "error fetching comments")
		return
	}

	log.Printf("you successfully got %d comments", len(comments))
	c.JSON(http.StatusOK, comments)
}

// @Summary Get Comment By ID
// @Security ApiKeyAuth
// @Tags comments
// @Description get comment by ID
// @ID get-comment-by-id
// @Param id path int true "Comment ID"
// @Produce json
// @Success 200 {object} models.Comment "Returns the comment object."
// @Failure 404 {string} string "Comment with the specified ID does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments/{id} [get]
// getCommentByID is a handler for getting a comment by its ID. It returns the requested comment if found.
func (h *Handler) getCommentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("invalid id param: %v", err)
		c.String(http.StatusNotFound, "invalid id param")
		return
	}

	comment, err := h.storage.Comments.Get(uint(id))
	if err != nil {
		log.Printf("incorrect comment id: %v", err)
		c.String(http.StatusInternalServerError, "incorrect comment id")
		return
	}

	log.Printf("you successfully got comment #%d: %+v", comment.ID, comment)
	c.JSON(http.StatusOK, comment)
}

// @Summary Create Comment
// @Security ApiKeyAuth
// @Tags comments
// @Description create comment
// @ID create-comment
// @Param post_id query int true "Post ID"
// @Param user_id query int true "User ID"
// @Param comment body models.UpdateCommentInput true "Comment credentials"
// @Accept json
// @Produce json
// @Success 201 {integer} integer "Comment successfully created. Returns the newly created comment id."
// @Failure 400 {string} string "Invalid request or missing required fields."
// @Failure 404 {string} string "Post or User with the specified ID does not exists."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments [post]
// createComment is a handler for creating new comment. It returns the ID of the newly created comment.
func (h *Handler) createComment(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		log.Printf("no such user in the database: %v\nplease provide correct user_id by query parameter like '?user_id=1'", err)
		c.String(http.StatusNotFound, "no such user in the database\nplease provide correct user_id by query parameter like '?user_id=1'")
		return
	}

	postID, err := strconv.Atoi(c.Query("post_id"))
	if err != nil {
		log.Printf("no such post in the database: %v\nplease provide correct post_id by query parameter like '?post_id=1'", err)
		c.String(http.StatusNotFound, "no such post in the database\nplease provide correct post_id by query parameter like '?post_id=1'")
		return
	}

	comment := new(models.Comment)

	if err := c.Bind(&comment); err != nil {
		log.Printf("invalid comment data: %v", err)
		c.String(http.StatusBadRequest, "invalid comment data")
		return
	}

	id, err := h.storage.Comments.Create(comment, uint(userID), uint(postID))
	if err != nil {
		log.Printf("error creating comment: %v", err)
		c.String(http.StatusInternalServerError, "error creating comment")
		return
	}

	log.Printf("comment successfully created by id %d", id)
	c.String(http.StatusCreated, "comment successfully created by id %d", id)
}

// @Summary Update Comment
// @Security ApiKeyAuth
// @Tags comments
// @Description update an existing comment
// @ID update-comment
// @Param id path int true "Comment ID"
// @Accept json
// @Produce json
// @Param comment body models.UpdateCommentInput true "Comment data to update"
// @Success 200 {object} models.UpdateCommentInput "Comment successfully updated."
// @Failure 400 {string} string "Invalid request or missing required."
// @Failure 404 {string} string "Comment with the specified ID does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments/{id} [patch]
// updateComment is a handler for updating an existing comment. It updates the comment with the specified ID.
func (h *Handler) updateComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("invalid id param: %v", err)
		c.String(http.StatusNotFound, "invalid id param")
		return
	}

	var input models.UpdateCommentInput

	if err := c.Bind(&input); err != nil {
		log.Printf("ivalid comment input: %v", err)
		c.String(http.StatusNotFound, "ivalid comment input")
		return
	}

	if err := h.storage.Comments.Update(uint(id), &input); err != nil {
		log.Printf("error updating comment: %v", err)
		c.String(http.StatusInternalServerError, "error updating comment")
		return
	}

	log.Printf("comment #%d successfully updated:\nemail: %+v\nbody: %+v", id, *input.Email, *input.Body)
	c.String(http.StatusOK, "comment #%d successfully updated:\nemail: %+v\nbody: %+v", id, *input.Email, *input.Body)
}

// @Summary Delete Comment
// @Security ApiKeyAuth
// @Tags comments
// @Description delete a comment by ID
// @ID delete-comment
// @Param id path int true "Comment ID"
// @Produce json
// @Success 200 {string} string "Comment successfully deleted."
// @Failure 404 {string} string "Comment with the specified ID does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/comments/{id} [delete]
// deleteComment is a handler for deleting a comment by its ID. It deletes the comment with the specified ID.
func (h *Handler) deleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("invalid id param: %v", err)
		c.String(http.StatusNotFound, "invalid id param")
		return
	}

	if err := h.storage.Comments.Delete(uint(id)); err != nil {
		log.Printf("error deleting comment: %v", err)
		c.String(http.StatusInternalServerError, "error deleting comment")
		return
	}

	log.Printf("comment #%d successfully deleted", id)
	c.String(http.StatusOK, "comment #%d successfully deleted", id)
}
