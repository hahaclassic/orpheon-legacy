package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrRegisterUser   = errors.New("register user error")
	ErrLogin          = errors.New("login error")
	ErrLogout         = errors.New("logout error")
	ErrGetClaims      = errors.New("get claims error")
	ErrUpdatePassword = errors.New("update password error")
	ErrRefreshTokens  = errors.New("refresh tokens errors")
)

type AuthService interface {
	RegisterUser(ctx context.Context, credentials *entity.UserCredentials) (*entity.AuthTokens, error)
	Login(ctx context.Context, credentials *entity.UserCredentials) (*entity.AuthTokens, error)

	Logout(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*entity.AuthTokens, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwords *entity.UserPasswords) error

	GetClaims(ctx context.Context, accessToken string) (*entity.Claims, error)
}
