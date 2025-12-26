package jwttokens

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrExpired              = errors.New("access token expired error")
	ErrInvalidClaims        = errors.New("invalid claims error")
	ErrInvalidSigningMethod = errors.New("invalid signing method error")
	ErrTokenParsing         = errors.New("token parsing error")
	ErrInvalidAccessToken   = errors.New("invalid access token")
)

type jwtClaims struct {
	UserID    uuid.UUID          `json:"user_id"`
	AccessLvl entity.AccessLevel `json:"access_level"`
	jwt.RegisteredClaims
}

type AccessTokenConfig = config.AccessTokenConfig

// type AccessTokenConfig struct {
// 	TTL    time.Duration
// 	Jitter time.Duration
// }

type JWTTokenService struct {
	cfg AccessTokenConfig
}

func New(config AccessTokenConfig) *JWTTokenService {
	return &JWTTokenService{
		cfg: config,
	}
}

func (s *JWTTokenService) GenerateAccessToken(claims *entity.Claims) (string, error) {
	jitter := time.Duration(rand.Int63n(int64(s.cfg.Jitter)))
	exp := time.Now().Add(s.cfg.TTL + jitter)

	jwtClaims := jwtClaims{
		UserID:    claims.UserID,
		AccessLvl: claims.AccessLvl,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	return token.SignedString(s.cfg.SecretKey)
}

func (s *JWTTokenService) ParseAccessToken(tokenStr string) (*entity.Claims, error) {
	var jwtClaims jwtClaims

	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return s.cfg.SecretKey, nil
	})

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, ErrExpired
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return nil, ErrInvalidClaims
	case err != nil:
		return nil, ErrTokenParsing
	case !token.Valid:
		return nil, ErrInvalidAccessToken
	}

	return &entity.Claims{
		UserID:    jwtClaims.UserID,
		AccessLvl: jwtClaims.AccessLvl,
	}, nil
}

func (s *JWTTokenService) GenerateRefreshToken() (string, error) {
	newRefresh, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return newRefresh.String(), nil
}
