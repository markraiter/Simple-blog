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

func TestParseToken(t *testing.T) {
	cfg := config.Auth{
		SigningKey: "testKey",
	}

	user := model.User{
		ID:       111,
		Username: "testUser",
		Email:    "test@test.com",
	}

	duration := time.Minute

	tests := []struct {
		name       string
		user       *model.User
		duration   time.Duration
		wantClaims *TokenClaims
		wantErr    error
	}{
		{
			name:       "valid token",
			user:       &user,
			duration:   duration,
			wantClaims: &TokenClaims{UID: "111", Username: "testUser", Email: "test@test.com", Exp: time.Now().Add(duration).Unix()},
			wantErr:    nil,
		},
		{
			name:       "token expired",
			user:       &user,
			duration:   -time.Minute,
			wantClaims: nil,
			wantErr:    ErrTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(cfg, tt.user, tt.duration)
			assert.NoError(t, err)

			claims, err := ParseToken(token, cfg.SigningKey)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.wantClaims != nil {
				assert.Equal(t, tt.wantClaims.UID, claims.UID)
				assert.Equal(t, tt.wantClaims.Username, claims.Username)
				assert.Equal(t, tt.wantClaims.Email, claims.Email)
			}
		})
	}
}
