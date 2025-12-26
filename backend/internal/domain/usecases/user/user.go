package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrCreateUser = errors.New("failed to create user")
	ErrGetUser    = errors.New("failed to get user")
	ErrUpdateUser = errors.New("failed to update user")
	ErrDeleteUser = errors.New("failed to delete user")
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.UserInfo) (uuid.UUID, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.UserInfo, error)
	UpdateUser(ctx context.Context, claims *entity.Claims, user *entity.UserInfo) error
	DeleteUser(ctx context.Context, claims *entity.Claims, userID uuid.UUID) error
}
