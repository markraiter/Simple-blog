package jwt

import (
	"testing"
	"time"

	"github.com/markraiter/simple-blog/config"
	"github.com/markraiter/simple-blog/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	cfg := config.Auth{
		SigningKey: "testKey",
	}

	user := model.User{
		ID:       111,
		Username: "testUser",
		Email:    "test@test.com",
	}
	duration := time.Minute

	token, err := NewToken(cfg, &user, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
