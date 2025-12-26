package genre

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrAssignGenreToTrack     = errors.New("assign genre to track")
	ErrUnassignGenreFromTrack = errors.New("unassign genre from track")
	ErrAssignGenreToAlbum     = errors.New("assign genre to album")
	ErrUnassignGenreFromAlbum = errors.New("unassign genre from album")
)

type GenreAssignService interface {
	AssignGenreToAlbum(ctx context.Context, claims *entity.Claims, genreID uuid.UUID, albumID uuid.UUID) error
	UnassignGenreFromAlbum(ctx context.Context, claims *entity.Claims, genreID uuid.UUID, albumID uuid.UUID) error
}
