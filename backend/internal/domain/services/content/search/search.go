package search

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type SearchRepository interface {
	SearchTracks(ctx context.Context, req *entity.SearchRequest) ([]*entity.TrackMeta, error)
	SearchAlbums(ctx context.Context, req *entity.SearchRequest) ([]*entity.AlbumMeta, error)
	SearchArtists(ctx context.Context, req *entity.SearchRequest) ([]*entity.ArtistMeta, error)
	SearchPlaylists(ctx context.Context, req *entity.SearchRequest) ([]*entity.PlaylistMeta, error)
}

type SearchService struct {
	repo SearchRepository
}

func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

func (s *SearchService) SearchTracks(ctx context.Context, req *entity.SearchRequest) (_ []*entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrSearchTracks, err)
	}()

	return s.repo.SearchTracks(ctx, req)
}

func (s *SearchService) SearchAlbums(ctx context.Context, req *entity.SearchRequest) (_ []*entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrSearchAlbums, err)
	}()

	return s.repo.SearchAlbums(ctx, req)
}

func (s *SearchService) SearchArtists(ctx context.Context, req *entity.SearchRequest) (_ []*entity.ArtistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrSearchArtists, err)
	}()

	return s.repo.SearchArtists(ctx, req)
}

func (s *SearchService) SearchPlaylists(ctx context.Context, claims *entity.Claims, req *entity.SearchRequest) (_ []*entity.PlaylistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrSearchPlaylists, err)
	}()

	if claims == nil {
		claims = &entity.Claims{
			UserID: uuid.Nil,
		}
	}

	playlists, err := s.repo.SearchPlaylists(ctx, req)
	if err != nil {
		return nil, err
	}

	availablePlaylists := make([]*entity.PlaylistMeta, 0, len(playlists))
	for i := range playlists {
		if !playlists[i].IsPrivate || playlists[i].OwnerID == claims.UserID {
			availablePlaylists = append(availablePlaylists, playlists[i])
		}
	}

	return availablePlaylists, nil
}
