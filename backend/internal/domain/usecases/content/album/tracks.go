package album

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetAllTracks = errors.New("get all tracks error")
)

type AlbumTrackService interface {
	GetAllTracks(ctx context.Context, albumID uuid.UUID) ([]*entity.TrackMeta, error)
}
