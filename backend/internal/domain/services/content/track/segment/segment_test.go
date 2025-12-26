package tracksegment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestTrackSegmentServiceSuite(t *testing.T) {
	suite.Run(t, &TrackSegmentServiceSuite{})
}

type TrackSegmentObjectMother struct{}

func (TrackSegmentObjectMother) DefaultTrackID() uuid.UUID {
	return uuid.New()
}

func (TrackSegmentObjectMother) DefaultSegments(trackID uuid.UUID) []*entity.Segment {
	segments := make([]*entity.Segment, 60)
	for i := 0; i < 60; i++ {
		segments[i] = &entity.Segment{
			TrackID: trackID,
			Idx:     i,
			Range:   &entity.Range{Start: i * 3, End: (i + 1) * 3},
		}
	}
	return segments
}

type TrackSegmentServiceSuite struct {
	suite.Suite

	ctx     context.Context
	service *Service
	repo    *mocks.TrackSegmentRepository

	objMother *TrackSegmentObjectMother
}

func (s *TrackSegmentServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewTrackSegmentRepository(s.T())
	s.service = NewTrackSegmentService(s.repo)
	s.objMother = &TrackSegmentObjectMother{}
}

func (s *TrackSegmentServiceSuite) TestGetSegments_Success() {
	trackID := s.objMother.DefaultTrackID()
	expected := s.objMother.DefaultSegments(trackID)
	s.repo.On("GetSegments", s.ctx, trackID).Return(expected, nil)

	result, err := s.service.GetSegments(s.ctx, trackID)

	s.NoError(err)
	s.Equal(expected, result)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackSegmentServiceSuite) TestGetSegments_InvalidTrackID() {
	result, err := s.service.GetSegments(s.ctx, uuid.Nil)

	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, ErrInvalidTrackID)
}

func (s *TrackSegmentServiceSuite) TestCreateSegments_Success() {
	trackID := s.objMother.DefaultTrackID()
	duration := 180

	s.repo.On("CreateSegments", s.ctx, trackID, mock.Anything).Return(nil)

	err := s.service.CreateSegments(s.ctx, trackID, duration)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackSegmentServiceSuite) TestCreateSegments_InvalidDuration() {
	trackID := s.objMother.DefaultTrackID()
	err := s.service.CreateSegments(s.ctx, trackID, 0)

	s.Error(err)
	s.ErrorIs(err, ErrInvalidDuration)
}

func (s *TrackSegmentServiceSuite) TestDeleteSegments_Success() {
	trackID := s.objMother.DefaultTrackID()
	s.repo.On("DeleteSegments", s.ctx, trackID).Return(nil)

	err := s.service.DeleteSegments(s.ctx, trackID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackSegmentServiceSuite) TestDeleteSegments_InvalidTrackID() {
	err := s.service.DeleteSegments(s.ctx, uuid.Nil)

	s.Error(err)
	s.ErrorIs(err, ErrInvalidTrackID)
}

func (s *TrackSegmentServiceSuite) TestIncrementTotalStreams_Success() {
	trackID := s.objMother.DefaultTrackID()
	idxs := []int{0, 1, 2}
	s.repo.On("IncrementTotalStreams", s.ctx, trackID, idxs).Return(nil)

	err := s.service.IncrementTotalStreams(s.ctx, trackID, idxs)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *TrackSegmentServiceSuite) TestIncrementTotalStreams_InvalidTrackID() {
	err := s.service.IncrementTotalStreams(s.ctx, uuid.Nil, []int{0})

	s.Error(err)
	s.ErrorIs(err, ErrInvalidTrackID)
}
