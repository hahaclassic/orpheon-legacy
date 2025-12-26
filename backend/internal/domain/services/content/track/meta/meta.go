package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGenerateTrackID = errors.New("generate track id error")
)

type TrackMetaRepository interface {
	GetByID(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error)
	Create(ctx context.Context, track *entity.TrackMeta) error
	Update(ctx context.Context, track *entity.TrackMeta) error
	Delete(ctx context.Context, trackID uuid.UUID) error
}

type TrackMetaService struct {
	repo           TrackMetaRepository
	segmentService track.TrackSegmentService
}

func NewTrackMetaService(repo TrackMetaRepository, segmentService track.TrackSegmentService) *TrackMetaService {
	return &TrackMetaService{repo: repo, segmentService: segmentService}
}

func (s *TrackMetaService) GetTrackMeta(ctx context.Context, trackID uuid.UUID) (_ *entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetTrackMeta, err)
	}()

	track, err := s.repo.GetByID(ctx, trackID)
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (s *TrackMetaService) CreateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (_ uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateTrackMeta, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return uuid.Nil, commonerr.ErrForbidden
	}

	track.ID, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, ErrGenerateTrackID
	}

	if err = s.repo.Create(ctx, track); err != nil {
		return uuid.Nil, err
	}

	if err = s.segmentService.CreateSegments(ctx, track.ID, track.Duration); err != nil {
		return uuid.Nil, err
	}

	return track.ID, nil
}

func (s *TrackMetaService) UpdateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateTrackMeta, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.Update(ctx, track)
}

func (s *TrackMetaService) DeleteTrackMeta(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteTrackMeta, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	if err = s.segmentService.DeleteSegments(ctx, trackID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, trackID)
}
