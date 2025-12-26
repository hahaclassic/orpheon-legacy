package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrUploadCover = errors.New("failed to upload playlist cover")
	ErrGetCover    = errors.New("failed to get playlist cover")
	ErrDeleteCover = errors.New("failed to delete playlist cover")
)

type PlaylistCoverService interface {
	GetCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (*entity.Cover, error)
	UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error
	DeleteCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
