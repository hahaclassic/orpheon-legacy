package artist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetArtistAlbums  = errors.New("failed to get artist albums")
	ErrGetArtistTracks  = errors.New("failed to get artist tracks")
	ErrGetArtistByAlbum = errors.New("failed to get artists by album")
	ErrGetArtistByTrack = errors.New("failed to get artists by track")

	ErrAssignArtistOnTrack     = errors.New("failed to assign artist on track")
	ErrAssignArtistOnAlbum     = errors.New("failed to assign artist on album")
	ErrUnassignArtistFromTrack = errors.New("failed to unassign artist from track")
	ErrUnassignArtistFromAlbum = errors.New("failed to unassign artist from album")
)

type ArtistAssignService interface {
	GetArtistAlbums(ctx context.Context, artistID uuid.UUID) ([]*entity.AlbumMeta, error)
	GetArtistTracks(ctx context.Context, artistID uuid.UUID) ([]*entity.TrackMeta, error)
	GetArtistByAlbum(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error)
	GetArtistByTrack(ctx context.Context, trackID uuid.UUID) ([]*entity.ArtistMeta, error)

	AssignArtistToTrack(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, trackID uuid.UUID) error
	AssignArtistToAlbum(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, albumID uuid.UUID) error
	UnassignArtistFromTrack(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, trackID uuid.UUID) error
	UnassignArtistFromAlbum(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, albumID uuid.UUID) error
}
