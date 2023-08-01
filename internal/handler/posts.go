package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

func (h *Handler) getAllPosts(c *gin.Context) {
	posts, err := h.storage.Posts.GetAll()
	if err != nil {
		log.Printf("error fetching posts: %+v", err)
		c.String(http.StatusInternalServerError, "error fetching posts")
	}

	log.Printf("you successfully got %d posts", len(posts))
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) filterPostsByUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		log.Printf("no such user in the database: %v\nplease provide correct user_id by query parameter like '?user_id=1'", err)
		c.String(http.StatusNotFound, "no such user in the database\nplease provide correct user_id by query parameter like '?user_id=1'")
		return
	}

	posts, err := h.storage.Posts.Filter(uint(id))
	if err != nil {
		log.Printf("error fetching posts: %v", err)
		c.String(http.StatusInternalServerError, "error fetching posts")
		return
	}

	log.Printf("you successfully got %d posts", len(posts))
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) getPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("invalid id param: %v", err)
		c.String(http.StatusNotFound, "invalid id param")
		return
	}

	post, err := h.storage.Posts.Get(uint(id))
	if err != nil {
		log.Printf("incorrect post id: %v", err)
		c.String(http.StatusInternalServerError, "incorrect post id")
		return
	}

	log.Printf("you successfully got post #%d: %+v", post.ID, post)
	c.JSON(http.StatusOK, post)
}

func (h *Handler) createPost(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		log.Printf("no such user in the database: %v\nplease provide correct user_id by query parameter like '?user_id=1'", err)
		c.String(http.StatusBadRequest, "no such user in the database\nplease provide correct user_id by query parameter like '?user_id=1'")
		return
	}

	post := new(models.Post)

	if err := c.Bind(&post); err != nil {
		log.Printf("invalid post data: %v", err)
		c.String(http.StatusBadRequest, "invalid post data")
		return
	}

	id, err := h.storage.Posts.Create(post, uint(userID))
	if err != nil {
		log.Printf("error creating post: %v", err)
		c.String(http.StatusInternalServerError, "error creating post")
		return
	}

	log.Printf("post successfully created by id %d", id)
	c.String(http.StatusCreated, "post successfully created by id %d", id)
}

func (h *Handler) updatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("invalid id param: %v", err)
		c.String(http.StatusNotFound, "invalid id param")
		return
	}

	var input models.UpdatePostInput

	if err := c.Bind(&input); err != nil {
		log.Printf("ivalid post input: %v", err)
		c.String(http.StatusNotFound, "ivalid post input")
		return
	}

	if err := h.storage.Posts.Update(uint(id), &input); err != nil {
		log.Printf("error updating post: %v", err)
		c.String(http.StatusInternalServerError, "error updating post")
		return
	}

	log.Printf("post #%d successfully updated:\ntitle: %+v\nbody: %+v", id, *input.Title, *input.Body)
	c.String(http.StatusOK, "post #%d successfully updated:\ntitle: %+v\nbody: %+v", id, *input.Title, *input.Body)
}

func (h *Handler) deletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("invalid id param: %v", err)
		c.String(http.StatusNotFound, "invalid id param")
		return
	}

	if err := h.storage.Posts.Delete(uint(id)); err != nil {
		log.Printf("error deleting post: %v", err)
		c.String(http.StatusInternalServerError, "error deleting post")
		return
	}

	log.Printf("post #%d successfully deleted", id)
	c.String(http.StatusOK, "post #%d successfully deleted", id)
}
