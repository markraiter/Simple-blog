package handlers

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/markraiter/simple-blog/internal/models"
	"gorm.io/gorm"
)

// Getting all comments
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

// Getting comment by ID
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

// Creating new comment
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

// Updating comment
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

// Deleting comment
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