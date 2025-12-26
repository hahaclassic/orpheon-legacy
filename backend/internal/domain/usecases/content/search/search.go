package search

import (
	"context"
	"errors"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrSearchTracks    = errors.New("failed to search tracks")
	ErrSearchAlbums    = errors.New("failed to search albums")
	ErrSearchArtists   = errors.New("failed to search artists")
	ErrSearchPlaylists = errors.New("failed to search playlists")
)

type SearchService interface {
	SearchTracks(ctx context.Context, request *entity.SearchRequest) ([]*entity.TrackMeta, error)
	SearchAlbums(ctx context.Context, request *entity.SearchRequest) ([]*entity.AlbumMeta, error)
	SearchArtists(ctx context.Context, request *entity.SearchRequest) ([]*entity.ArtistMeta, error)
	SearchPlaylists(ctx context.Context, claims *entity.Claims, request *entity.SearchRequest) ([]*entity.PlaylistMeta, error)
}
