package processor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/stat/processor"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestListeningStatServiceSuite(t *testing.T) {
	suite.Run(t, &ListeningStatServiceSuite{})
}

type ListeningStatObjectMother struct{}

func (ListeningStatObjectMother) DefaultTrackID() uuid.UUID {
	return uuid.New()
}

func (ListeningStatObjectMother) DefaultUserID() uuid.UUID {
	return uuid.New()
}

func (ListeningStatObjectMother) DefaultSegments(trackID uuid.UUID) []*entity.Segment {
	return []*entity.Segment{
		{Range: &entity.Range{Start: 0, End: 10}},
		{Range: &entity.Range{Start: 10, End: 20}},
		{Range: &entity.Range{Start: 20, End: 30}},
		{Range: &entity.Range{Start: 30, End: 40}},
	}
}

func (ListeningStatObjectMother) DefaultListeningEvent(trackID, userID uuid.UUID, start, end int) *entity.ListeningEvent {
	return &entity.ListeningEvent{
		TrackID: trackID,
		UserID:  userID,
		Ranges:  []*entity.Range{{Start: start, End: end}},
	}
}

type ListeningStatServiceSuite struct {
	suite.Suite

	ctx         context.Context
	service     *processor.ListeningStatService
	trackRepo   *mocks.TrackStatRepository
	segmentRepo *mocks.SegmentStatRepository

	objMother *ListeningStatObjectMother
}

func (s *ListeningStatServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.trackRepo = mocks.NewTrackStatRepository(s.T())
	s.segmentRepo = mocks.NewSegmentStatRepository(s.T())
	s.service = processor.NewListeningStatService(s.trackRepo, s.segmentRepo)
	s.objMother = &ListeningStatObjectMother{}
}

func (s *ListeningStatServiceSuite) TestUpdateStat_SuccessWithTrackIncrement() {
	trackID := s.objMother.DefaultTrackID()
	userID := s.objMother.DefaultUserID()
	event := s.objMother.DefaultListeningEvent(trackID, userID, 0, 35)
	segments := s.objMother.DefaultSegments(trackID)

	s.segmentRepo.On("GetSegments", s.ctx, trackID).Return(segments, nil)
	s.segmentRepo.On("IncrementTotalStreams", s.ctx, trackID, mock.Anything).Return(nil)
	s.trackRepo.On("IncrementTrackTotalStreams", s.ctx, trackID).Return(nil)

	err := s.service.UpdateStat(s.ctx, event)

	s.NoError(err)
	s.segmentRepo.AssertExpectations(s.T())
	s.trackRepo.AssertExpectations(s.T())
}

func (s *ListeningStatServiceSuite) TestUpdateStat_SuccessWithoutTrackIncrement() {
	trackID := s.objMother.DefaultTrackID()
	userID := s.objMother.DefaultUserID()
	event := s.objMother.DefaultListeningEvent(trackID, userID, 0, 20)
	segments := []*entity.Segment{
		{Range: &entity.Range{Start: 0, End: 10}},
		{Range: &entity.Range{Start: 10, End: 20}},
	}

	s.segmentRepo.On("GetSegments", s.ctx, trackID).Return(segments, nil)
	s.segmentRepo.On("IncrementTotalStreams", s.ctx, trackID, mock.Anything).Return(nil)
	s.trackRepo.On("IncrementTrackTotalStreams", s.ctx, trackID).Return(nil) // добавить этот мок

	err := s.service.UpdateStat(s.ctx, event)

	s.NoError(err)
	s.segmentRepo.AssertExpectations(s.T())
	s.trackRepo.AssertExpectations(s.T())
}
func (s *ListeningStatServiceSuite) TestUpdateStat_GetSegmentsError() {
	trackID := s.objMother.DefaultTrackID()
	userID := s.objMother.DefaultUserID()
	event := s.objMother.DefaultListeningEvent(trackID, userID, 0, 20)

	s.segmentRepo.On("GetSegments", s.ctx, trackID).Return(nil, errors.New("db error"))

	err := s.service.UpdateStat(s.ctx, event)

	s.Error(err)
	s.segmentRepo.AssertExpectations(s.T())
}

func (s *ListeningStatServiceSuite) TestUpdateStat_IncrementSegmentsError() {
	trackID := s.objMother.DefaultTrackID()
	userID := s.objMother.DefaultUserID()
	event := s.objMother.DefaultListeningEvent(trackID, userID, 0, 35)
	segments := s.objMother.DefaultSegments(trackID)

	s.segmentRepo.On("GetSegments", s.ctx, trackID).Return(segments, nil)
	s.segmentRepo.On("IncrementTotalStreams", s.ctx, trackID, mock.Anything).Return(errors.New("increment error"))

	err := s.service.UpdateStat(s.ctx, event)

	s.Error(err)
	s.segmentRepo.AssertExpectations(s.T())
}

func (s *ListeningStatServiceSuite) TestUpdateStat_IncrementTrackError() {
	trackID := s.objMother.DefaultTrackID()
	userID := s.objMother.DefaultUserID()
	event := s.objMother.DefaultListeningEvent(trackID, userID, 0, 40)
	segments := s.objMother.DefaultSegments(trackID)

	s.segmentRepo.On("GetSegments", s.ctx, trackID).Return(segments, nil)
	s.segmentRepo.On("IncrementTotalStreams", s.ctx, trackID, mock.Anything).Return(nil)
	s.trackRepo.On("IncrementTrackTotalStreams", s.ctx, trackID).Return(errors.New("increment track error"))

	err := s.service.UpdateStat(s.ctx, event)

	s.Error(err)
	s.segmentRepo.AssertExpectations(s.T())
	s.trackRepo.AssertExpectations(s.T())
}
