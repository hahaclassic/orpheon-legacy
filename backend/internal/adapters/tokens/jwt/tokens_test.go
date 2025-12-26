package jwttokens

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessToken(t *testing.T) {
	config := AccessTokenConfig{
		TTL:       time.Minute * 15,
		Jitter:    time.Second * 5,
		SecretKey: []byte("supersecretkey"),
	}
	tokenService := New(config)

	claims := &entity.Claims{
		UserID:    uuid.New(),
		AccessLvl: entity.AccessLevel(1),
	}

	token, err := tokenService.GenerateAccessToken(claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedClaims, err := tokenService.ParseAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, claims.UserID, parsedClaims.UserID)
	assert.Equal(t, claims.AccessLvl, parsedClaims.AccessLvl)
}

func TestTokenExpired(t *testing.T) {
	config := AccessTokenConfig{
		TTL:       time.Second,
		Jitter:    time.Second * 2,
		SecretKey: []byte("supersecretkey"),
	}
	tokenService := New(config)

	claims := &entity.Claims{
		UserID:    uuid.New(),
		AccessLvl: entity.AccessLevel(1),
	}

	token, err := tokenService.GenerateAccessToken(claims)
	assert.NoError(t, err)

	time.Sleep(time.Second * 3)

	_, err = tokenService.ParseAccessToken(token)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrExpired))
}

func TestInvalidClaims(t *testing.T) {
	config := AccessTokenConfig{
		TTL:       time.Minute * 15,
		Jitter:    time.Second * 5,
		SecretKey: []byte("supersecretkey"),
	}
	tokenService := New(config)

	claims := &entity.Claims{
		UserID:    uuid.New(),
		AccessLvl: entity.AccessLevel(1),
	}

	token, err := tokenService.GenerateAccessToken(claims)
	assert.NoError(t, err)

	parsedClaims, err := tokenService.ParseAccessToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)
	assert.Equal(t, claims.UserID, parsedClaims.UserID)
	assert.Equal(t, claims.AccessLvl, parsedClaims.AccessLvl)
}

func TestTokenParsingError(t *testing.T) {
	config := AccessTokenConfig{
		TTL:       time.Minute * 15,
		Jitter:    time.Second * 5,
		SecretKey: []byte("supersecretkey"),
	}
	tokenService := New(config)

	_, err := tokenService.ParseAccessToken("invalid.token.string")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrTokenParsing))
}

func TestTokenFromDifferentSource(t *testing.T) {
	config := AccessTokenConfig{
		TTL:       time.Minute * 15,
		Jitter:    time.Second * 5,
		SecretKey: []byte("supersecretkey"),
	}
	tokenService := New(config)

	claims := &entity.Claims{
		UserID:    uuid.New(),
		AccessLvl: entity.AccessLevel(1),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		UserID:    claims.UserID,
		AccessLvl: claims.AccessLvl,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TTL)),
		},
	})

	otherSecretKey := []byte("anothersecretkey")
	tokenStr, err := token.SignedString(otherSecretKey)
	assert.NoError(t, err)

	_, err = tokenService.ParseAccessToken(tokenStr)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrTokenParsing))
}
