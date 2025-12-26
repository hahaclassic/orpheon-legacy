package tracks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/tracks"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestPlaylistTrackServiceSuite(t *testing.T) {
	suite.Run(t, &PlaylistTrackServiceSuite{})
}

type PlaylistTrackServiceSuite struct {
	suite.Suite
	ctx           context.Context
	service       *tracks.PlaylistTrackService
	repo          *mocks.PlaylistTracksRepository
	policy        *mocks.PlaylistPolicyService
	userID        uuid.UUID
	playlistID    uuid.UUID
	trackID       uuid.UUID
	claims        *entity.Claims
	playlistTrack *entity.PlaylistTrack
}

func (s *PlaylistTrackServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewPlaylistTracksRepository(s.T())
	s.policy = mocks.NewPlaylistPolicyService(s.T())
	s.service = tracks.NewPlaylistTrackService(s.repo, s.policy)
	s.userID = uuid.New()
	s.playlistID = uuid.New()
	s.trackID = uuid.New()
	s.claims = &entity.Claims{UserID: s.userID}
	s.playlistTrack = &entity.PlaylistTrack{PlaylistID: s.playlistID, TrackID: s.trackID}
}

func (s *PlaylistTrackServiceSuite) TestAddTrackSuccess() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("AddTrackToPlaylist", s.ctx, s.playlistTrack).Return(nil)

	err := s.service.AddTrack(s.ctx, s.claims, s.playlistTrack)
	assert.NoError(s.T(), err)
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestAddTrackForbidden() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(commonerr.ErrForbidden)

	err := s.service.AddTrack(s.ctx, s.claims, s.playlistTrack)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)
	s.policy.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestAddTrackRepoError() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("AddTrackToPlaylist", s.ctx, s.playlistTrack).Return(errors.New("db error"))

	err := s.service.AddTrack(s.ctx, s.claims, s.playlistTrack)
	assert.ErrorContains(s.T(), err, "db error")
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestGetAllTracksSuccess() {
	s.SetupTest()
	expected := []*entity.TrackMeta{{ID: uuid.New()}}
	s.policy.On("CanView", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("GetAllPlaylistTracks", s.ctx, s.playlistID).Return(expected, nil)

	tracksRes, err := s.service.GetAllTracks(s.ctx, s.claims, s.playlistID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, tracksRes)
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestGetAllTracksForbidden() {
	s.SetupTest()
	s.policy.On("CanView", s.ctx, s.claims, s.playlistID).Return(commonerr.ErrForbidden)

	tracksRes, err := s.service.GetAllTracks(s.ctx, s.claims, s.playlistID)
	assert.Nil(s.T(), tracksRes)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)
	s.policy.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestGetAllTracksRepoError() {
	s.SetupTest()
	s.policy.On("CanView", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("GetAllPlaylistTracks", s.ctx, s.playlistID).Return(nil, errors.New("db error"))

	tracksRes, err := s.service.GetAllTracks(s.ctx, s.claims, s.playlistID)
	assert.Nil(s.T(), tracksRes)
	assert.ErrorContains(s.T(), err, "db error")
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestDeleteTrackSuccess() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("DeleteTrackFromPlaylist", s.ctx, s.playlistTrack).Return(nil)

	err := s.service.DeleteTrack(s.ctx, s.claims, s.playlistTrack)
	assert.NoError(s.T(), err)
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestDeleteTrackForbidden() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(commonerr.ErrForbidden)

	err := s.service.DeleteTrack(s.ctx, s.claims, s.playlistTrack)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)
	s.policy.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestDeleteTrackRepoError() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("DeleteTrackFromPlaylist", s.ctx, s.playlistTrack).Return(errors.New("db error"))

	err := s.service.DeleteTrack(s.ctx, s.claims, s.playlistTrack)
	assert.ErrorContains(s.T(), err, "db error")
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestDeleteAllTracksSuccess() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("DeleteAllTracksFromPlaylist", s.ctx, s.playlistID).Return(nil)

	err := s.service.DeleteAllTracks(s.ctx, s.claims, s.playlistID)
	assert.NoError(s.T(), err)
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestDeleteAllTracksForbidden() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(commonerr.ErrForbidden)

	err := s.service.DeleteAllTracks(s.ctx, s.claims, s.playlistID)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)
	s.policy.AssertExpectations(s.T())
}

func (s *PlaylistTrackServiceSuite) TestDeleteAllTracksRepoError() {
	s.SetupTest()
	s.policy.On("CanEdit", s.ctx, s.claims, s.playlistID).Return(nil)
	s.repo.On("DeleteAllTracksFromPlaylist", s.ctx, s.playlistID).Return(errors.New("db error"))

	err := s.service.DeleteAllTracks(s.ctx, s.claims, s.playlistID)
	assert.ErrorContains(s.T(), err, "db error")
	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}
