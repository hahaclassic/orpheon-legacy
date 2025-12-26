package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/meta"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestTrackMetaServiceSuite(t *testing.T) {
	suite.Run(t, &TrackMetaServiceSuite{})
}

type TrackMetaObjectMother struct{}

func (TrackMetaObjectMother) DefaultTrackMeta() *entity.TrackMeta {
	return &entity.TrackMeta{
		ID:       uuid.New(),
		Name:     "Track 1",
		Duration: 180,
	}
}

func (TrackMetaObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{UserID: uuid.New(), AccessLvl: entity.Admin}
}

func (TrackMetaObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{UserID: uuid.New(), AccessLvl: entity.User}
}

type TrackMetaServiceSuite struct {
	suite.Suite

	ctx            context.Context
	service        *meta.TrackMetaService
	repo           *mocks.TrackMetaRepository
	segmentService *mocks.TrackSegmentService

	objMother *TrackMetaObjectMother
}

func (s *TrackMetaServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewTrackMetaRepository(s.T())
	s.segmentService = mocks.NewTrackSegmentService(s.T())
	s.service = meta.NewTrackMetaService(s.repo, s.segmentService)
	s.objMother = &TrackMetaObjectMother{}
}

func (s *TrackMetaServiceSuite) TestGetTrackMeta_Success() {
	track := s.objMother.DefaultTrackMeta()
	s.repo.On("GetByID", s.ctx, track.ID).Return(track, nil)

	result, err := s.service.GetTrackMeta(s.ctx, track.ID)

	s.NoError(err)
	s.Equal(track, result)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackMetaServiceSuite) TestGetTrackMeta_Error() {
	trackID := uuid.New()
	s.repo.On("GetByID", s.ctx, trackID).Return(nil, errors.New("db error"))

	result, err := s.service.GetTrackMeta(s.ctx, trackID)

	s.Error(err)
	s.Nil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackMetaServiceSuite) TestCreateTrackMeta_Success() {
	track := s.objMother.DefaultTrackMeta()
	claims := s.objMother.AdminClaims()

	s.repo.On("Create", s.ctx, mock.Anything).Return(nil)
	s.segmentService.On("CreateSegments", s.ctx, mock.Anything, track.Duration).Return(nil)

	id, err := s.service.CreateTrackMeta(s.ctx, claims, track)

	s.NoError(err)
	s.NotEqual(uuid.Nil, id)
	s.repo.AssertExpectations(s.T())
	s.segmentService.AssertExpectations(s.T())
}

func (s *TrackMetaServiceSuite) TestCreateTrackMeta_Forbidden() {
	track := s.objMother.DefaultTrackMeta()
	claims := s.objMother.UserClaims()

	id, err := s.service.CreateTrackMeta(s.ctx, claims, track)

	s.Error(err)
	s.Equal(uuid.Nil, id)
}

func (s *TrackMetaServiceSuite) TestUpdateTrackMeta_Success() {
	track := s.objMother.DefaultTrackMeta()
	claims := s.objMother.AdminClaims()

	s.repo.On("Update", s.ctx, track).Return(nil)

	err := s.service.UpdateTrackMeta(s.ctx, claims, track)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackMetaServiceSuite) TestUpdateTrackMeta_Forbidden() {
	track := s.objMother.DefaultTrackMeta()
	claims := s.objMother.UserClaims()

	err := s.service.UpdateTrackMeta(s.ctx, claims, track)

	s.Error(err)
}

func (s *TrackMetaServiceSuite) TestDeleteTrackMeta_Success() {
	track := s.objMother.DefaultTrackMeta()
	claims := s.objMother.AdminClaims()

	s.segmentService.On("DeleteSegments", s.ctx, track.ID).Return(nil)
	s.repo.On("Delete", s.ctx, track.ID).Return(nil)

	err := s.service.DeleteTrackMeta(s.ctx, claims, track.ID)

	s.NoError(err)
	s.segmentService.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *TrackMetaServiceSuite) TestDeleteTrackMeta_Forbidden() {
	track := s.objMother.DefaultTrackMeta()
	claims := s.objMother.UserClaims()

	err := s.service.DeleteTrackMeta(s.ctx, claims, track.ID)

	s.Error(err)
}
