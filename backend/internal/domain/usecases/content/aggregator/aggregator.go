package aggregator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetTracksAggregated      = errors.New("failed to get tracks aggregated")
	ErrGetTracksAggregatedByIDs = errors.New("failed to get tracks aggregated by ids")
	ErrGetAlbumsAggregated      = errors.New("failed to get albums aggregated")
	ErrGetAlbumsAggregatedByIDs = errors.New("failed to get albums aggregated by ids")
)

type ContentAggregator interface {
	GetTracksByIDs(ctx context.Context, trackIDs ...uuid.UUID) (_ []*entity.TrackMetaAggregated, err error)
	GetAlbumsByIDs(ctx context.Context, albumIDs ...uuid.UUID) (_ []*entity.AlbumMetaAggregated, err error)
	GetAlbums(ctx context.Context, albums ...*entity.AlbumMeta) (_ []*entity.AlbumMetaAggregated, err error)
	GetTracks(ctx context.Context, tracks ...*entity.TrackMeta) (_ []*entity.TrackMetaAggregated, err error)
}
