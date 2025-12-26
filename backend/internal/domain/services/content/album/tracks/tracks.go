package album_tracks_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type AlbumTrackRepository interface {
	GetAllTracks(ctx context.Context, albumID uuid.UUID) ([]*entity.TrackMeta, error)
}

type AlbumTrackService struct {
	albumTrackRepository AlbumTrackRepository
}

func NewAlbumTrackService(albumTrackRepository AlbumTrackRepository) *AlbumTrackService {
	return &AlbumTrackService{albumTrackRepository: albumTrackRepository}
}

func (s *AlbumTrackService) GetAllTracks(ctx context.Context, albumID uuid.UUID) (tracks []*entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllTracks, err)
	}()

	return s.albumTrackRepository.GetAllTracks(ctx, albumID)
}
