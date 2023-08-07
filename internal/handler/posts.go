package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markraiter/simple-blog/models"
)

// @Summary Get All Posts
// @Security ApiKeyAuth
// @Tags posts
// @Description get all posts
// @ID get-all-posts
// @Produce  json
// @Success 200 {array} models.Post "Returns an array of posts."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/posts/all [get]
// getAllPosts is a handler for getting list of all the posts. It returns an array of posts
func (h *Handler) getAllPosts(c *gin.Context) {
	posts, err := h.storage.Posts.GetAll()
	if err != nil {
		log.Printf("error fetching posts: %+v", err)
		c.String(http.StatusInternalServerError, "error fetching posts")
	}

	log.Printf("you successfully got %d posts", len(posts))
	c.JSON(http.StatusOK, posts)
}

// @Summary Filter Posts by User
// @Security ApiKeyAuth
// @Tags posts
// @Description filter posts by user ID
// @ID filter-posts-by-user
// @Param user_id query int true "User ID"
// @Produce json
// @Success 200 {array} models.Post "Returns an array of posts created by the user."
// @Failure 404 {string} string "Post with specified user_id does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/posts [get]
// filterPostsByUser is a handler for filtering posts by user ID. It returns an array of posts associated with the specified user.
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

// @Summary Get Post By ID
// @Security ApiKeyAuth
// @Tags posts
// @Description get post by ID
// @ID get-post-by-id
// @Param id path int true "Post ID"
// @Produce json
// @Success 200 {object} models.Post "Returns the post object."
// @Failure 404 {string} string "Post with the specified ID does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/posts/{id} [get]
// getPostByID is a handler for getting a post by its ID. It returns the requested post if found.
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

// @Summary Create Post
// @Security ApiKeyAuth
// @Tags posts
// @Description create post
// @ID create-post
// @Param user_id query int true "User ID"
// Param input body models.UpdatePostInput true "Post credentials"
// @Accept json
// @Produce json
// @Success 201 {integer} integer "Post successfully created. Returns the newly created post id."
// @Failure 400 {string} string "Invalid request or missing required fields."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/posts [post]
// createPost is a handler for creating new post. It returns the ID of the newly created post.
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

// @Summary Update Post
// @Security ApiKeyAuth
// @Tags posts
// @Description update an existing post
// @ID update-post
// @Param id path int true "Post ID"
// @Accept json
// @Produce json
// @Param post body models.UpdatePostInput true "Post data to update"
// @Success 200 {object} models.UpdatePostInput "Post successfully updated."
// @Failure 400 {string} string "Invalid request or missing required."
// @Failure 404 {string} string "Post with the specified ID does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/posts/{id} [patch]
// updatePost is a handler for updating an existing post. It updates the post with the specified ID.
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

// @Summary Delete Post
// @Security ApiKeyAuth
// @Tags posts
// @Description delete a post by ID
// @ID delete-post
// @Param id path int true "Post ID"
// @Produce json
// @Success 200 {string} string "Post successfully deleted."
// @Failure 404 {string} string "Post with the specified ID does not exist."
// @Failure 500 {string} string "An unexpected error occurred on the server."
// @Router /api/posts/{id} [delete]
// deletePost is a handler for deleting a post by its ID. It deletes the post with the specified ID.
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
