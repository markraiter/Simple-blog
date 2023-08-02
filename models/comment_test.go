package models_test

import (
	"testing"

	"github.com/markraiter/simple-blog/models"
	"github.com/stretchr/testify/assert"
)

func TestComment_Validate(t *testing.T) {
	c := models.TestUpdateCommentInput(t)
	assert.NoError(t, c.Validate())
}
