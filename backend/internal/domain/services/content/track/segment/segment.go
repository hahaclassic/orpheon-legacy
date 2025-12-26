package tracksegment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

const (
	// SegmentSizeForUnder3min = 2 * time.Second
	// SegmentSizeFor3to7min   = 3 * time.Second
	// SegmentSizeForOver7min  = 5 * time.Second

	defaultSegmentCount = 60
)

var (
	ErrInvalidTrackID  = errors.New("invalid track ID")
	ErrInvalidDuration = errors.New("track duration must be positive")

	ErrCreateSegments   = errors.New("failed to create segments")
	ErrDeleteSegments   = errors.New("failed to delete segments")
	ErrGetSegments      = errors.New("failed to get segments")
	ErrIncrementStreams = errors.New("failed to increment total streams")
)

type TrackSegmentRepository interface {
	GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error)
	CreateSegments(ctx context.Context, trackID uuid.UUID, segments []*entity.Segment) error
	DeleteSegments(ctx context.Context, trackID uuid.UUID) error
	IncrementTotalStreams(ctx context.Context, trackID uuid.UUID, segmentsIdxs []int) error
}

type Service struct {
	repo TrackSegmentRepository
}

func NewTrackSegmentService(repo TrackSegmentRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error) {
	if trackID == uuid.Nil {
		return nil, ErrInvalidTrackID
	}
	segments, err := s.repo.GetSegments(ctx, trackID)
	if err != nil {
		return nil, ErrGetSegments
	}
	return segments, nil
}

func (s *Service) CreateSegments(ctx context.Context, trackID uuid.UUID, trackDuration int) error {
	if trackID == uuid.Nil {
		return ErrInvalidTrackID
	}
	if trackDuration <= 0 {
		return ErrInvalidDuration
	}

	segmentDuration := trackDuration / defaultSegmentCount
	segments := make([]*entity.Segment, 0, defaultSegmentCount)

	for i := range defaultSegmentCount {
		start := i * segmentDuration
		end := start + segmentDuration

		// last segmemt can be longer than others
		if i == defaultSegmentCount-1 {
			end = trackDuration
		}
		segments = append(segments, &entity.Segment{
			TrackID:      trackID,
			Idx:          i,
			Range:        &entity.Range{Start: int(start), End: int(end)},
			TotalStreams: 0,
		})
	}

	if err := s.repo.CreateSegments(ctx, trackID, segments); err != nil {
		return ErrCreateSegments
	}
	return nil
}

func (s *Service) DeleteSegments(ctx context.Context, trackID uuid.UUID) error {
	if trackID == uuid.Nil {
		return ErrInvalidTrackID
	}
	if err := s.repo.DeleteSegments(ctx, trackID); err != nil {
		return ErrDeleteSegments
	}
	return nil
}

func (s *Service) IncrementTotalStreams(ctx context.Context, trackID uuid.UUID, segmentsIdxs []int) error {
	if trackID == uuid.Nil {
		return ErrInvalidTrackID
	}
	if err := s.repo.IncrementTotalStreams(ctx, trackID, segmentsIdxs); err != nil {
		return ErrIncrementStreams
	}
	return nil
}
