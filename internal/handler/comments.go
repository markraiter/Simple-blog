package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

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
