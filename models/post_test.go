package models_test

import (
	"testing"

	"github.com/markraiter/simple-blog/models"
	"github.com/stretchr/testify/assert"
)

func TestPost_Validate(t *testing.T) {
	p := models.TestUpdatePostInput(t)
	assert.NoError(t, p.Validate())
}
