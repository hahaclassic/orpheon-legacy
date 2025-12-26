package album

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrUploadCover = errors.New("failed to upload album cover")
	ErrGetCover    = errors.New("failed to get album cover")
	ErrDeleteCover = errors.New("failed to delete album cover")
)

type AlbumCoverService interface {
	GetCover(ctx context.Context, albumID uuid.UUID) (*entity.Cover, error)
	// Admin
	UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error
	DeleteCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) error
}
