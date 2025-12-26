package track

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrCreateTrackMeta = errors.New("failed to create track meta")
	ErrGetTrackMeta    = errors.New("failed to get track meta")
	ErrUpdateTrackMeta = errors.New("failed to update track meta")
	ErrDeleteTrackMeta = errors.New("failed to delete track meta")
)

type TrackMetaService interface {
	GetTrackMeta(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error)
	CreateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (uuid.UUID, error)
	UpdateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) error
	DeleteTrackMeta(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error
}
