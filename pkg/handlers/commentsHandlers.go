package handlers

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/models"
	"gorm.io/gorm"
	_ "github.com/markraiter/simple-blog/docs"
)

// @Summary Get all comments
// @Description Get all comments from the database.
// @Tags Comments
// @Accept  json
// @Produce  json
// @Produce  xml
// @Security ApiKeyAuth
// @Success 200 {array} models.Comment
// @Failure 500 {string} error fetching comments
// @Router /api/v1/comments [get]
func GetComments(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		responseType := c.QueryParam("type")
		if responseType != "xml" {
			responseType = "json"
		}
		
		comments := []models.Comment{}

		if err := db.Find(&comments).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error retreiving comments")
		}

		if responseType == "xml" {
			xmlResponse, err := xml.MarshalIndent(comments, "", " ")
			if err != nil {
				return c.String(http.StatusInternalServerError, "error encoding xml")
			}

			c.Response().Header().Set("content-type", "application/xml")

			return c.Blob(http.StatusOK, "application/xml", xmlResponse)
			
		} else {
			c.Response().Header().Set("content-type", "application/json")

			return c.JSON(http.StatusOK, comments)
		}
	}
}

// @Summary Get comment by ID
// @Description Get particular comment from the database by unique ID
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Comment
// @Failure 400 {string} invalid comment ID
// @Failure 404 {string} comment not found
// @Router /api/v1/comments/{id} [get] 
func GetCommentByID(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid comment ID")
		}

		comment := new(models.Comment)
		
		if err := db.First(comment, commentID).Error; err != nil {
			return c.String(http.StatusNotFound, "comment not found")
		}

		return c.JSON(http.StatusOK, comment)
	}
}

// @Summary Create comment
// @Description Create new comment in the database
// @Tags Comments
// @Accept json
// @Produce json
// @Param post body models.Comment true "Comment data"
// @Security ApiKeyAuth
// @Success 201 {object} models.Comment
// @Failure 400 {string} invalid comment data
// @Failure 500 {string} error creating comment
// @Router /api/v1/comments [post]
func CreateComment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		comment := new(models.Comment)

		if err := c.Bind(comment); err != nil {
			return c.String(http.StatusBadRequest, "invalid comment data")
		}

		if err := db.Create(comment).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error creating comment")
		}

		return c.JSON(http.StatusCreated, comment)
	}
}

// @Summary Update comment
// @Description Update information in the particular comment from the database by unique ID
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param post body models.Comment true "Comment data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Comment
// @Failure 400 {string} invalid comment data
// @Failure 404 {string} comment not found
// @Failure 500 {string} error updating comment
// @Router /api/v1/comments/{id} [put]
func UpdateComment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid comment ID")
		}

		comment := new(models.Comment)

		if err := db.First(comment, commentID).Error; err != nil {
			return c.String(http.StatusNotFound, "comment not found")
		}

		if err := c.Bind(comment); err != nil {
			return c.String(http.StatusBadRequest, "invalid comment data")
		}

		if err := db.Save(comment).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error updating comment")
		}

		return c.JSON(http.StatusOK, comment)
	}
}

// @Summary Delete comment
// @Description Delete particular comment from the database by unique ID
// @Tags Comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 204 {string} comment deleted successfully
// @Failure 400 {string} invalid comment ID
// @Failure 404 {string} comment not found
// @Failure 500 {string} error deleting comment
// @Router /api/v1/comments/{id} [delete]
func DeleteComment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusNotFound, "invalid comment id")
		}

		comment := new(models.Comment)

		if err := db.First(comment, commentID).Error; err != nil {
			return c.String(http.StatusNotFound, "comment not found")
		}

		if err := db.Delete(comment, commentID).Error; err != nil {
			return c.String(http.StatusInternalServerError, "error deleting comment")
		}

		return c.NoContent(http.StatusNoContent)
	}
}