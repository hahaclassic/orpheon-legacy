package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetPlaylistsAggregated      = errors.New("failed to get playlists aggregated")
	ErrGetPlaylistsAggregatedByIDs = errors.New("failed to get playlists aggregated by ids")
)

type PlaylistAggregator interface {
	GetPlaylists(ctx context.Context, claims *entity.Claims, playlists ...*entity.PlaylistMeta) ([]*entity.PlaylistMetaAggregated, error)
	GetPlaylistsByIDs(ctx context.Context, claims *entity.Claims, playlistIDs ...uuid.UUID) ([]*entity.PlaylistMetaAggregated, error)
}
