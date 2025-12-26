package artist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrUploadAvatar = errors.New("failed to upload artist avatar")
	ErrGetAvatar    = errors.New("failed to get artist avatar")
	ErrDeleteAvatar = errors.New("failed to delete artist avatar")
)

type ArtistAvatarService interface {
	GetCover(ctx context.Context, artistID uuid.UUID) (*entity.Cover, error)
	// Admin
	UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error
	DeleteCover(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) error
}
