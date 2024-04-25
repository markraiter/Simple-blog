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

	token, err := NewToken(cfg, user, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestNewTokenPair(t *testing.T) {
	cfg := config.Auth{
		SigningKey: "testKey",
		AccessTTL:  time.Minute,
		RefreshTTL: time.Minute,
	}

	user := &model.User{
		ID: 111,
	}

	tokenPair, err := NewTokenPair(cfg, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
}
