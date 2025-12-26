package artist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetArtistMeta    = errors.New("get artist meta error")
	ErrGetAllArtistMeta = errors.New("get all artist meta error")
	ErrCreateArtistMeta = errors.New("create artist meta error")
	ErrUpdateArtistMeta = errors.New("update artist meta error")
	ErrDeleteArtistMeta = errors.New("delete artist meta error")
)

type ArtistMetaService interface {
	GetArtistMeta(ctx context.Context, artistID uuid.UUID) (*entity.ArtistMeta, error)
	GetAllArtistMeta(ctx context.Context) ([]*entity.ArtistMeta, error)
	CreateArtistMeta(ctx context.Context, claims *entity.Claims, artist *entity.ArtistMeta) error
	UpdateArtistMeta(ctx context.Context, claims *entity.Claims, artist *entity.ArtistMeta) error
	DeleteArtistMeta(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) error
}
