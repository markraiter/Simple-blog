package models_test

import (
	"testing"

	"github.com/markraiter/simple-blog/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := models.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.PasswordHash)
}
