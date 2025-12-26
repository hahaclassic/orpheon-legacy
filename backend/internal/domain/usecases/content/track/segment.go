package track

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type TrackSegmentService interface {
	GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
	CreateSegments(ctx context.Context, trackID uuid.UUID, trackDuration int) error
	IncrementTotalStreams(ctx context.Context, trackID uuid.UUID, segmentsIdxs []int) error
	DeleteSegments(ctx context.Context, trackID uuid.UUID) error
}
