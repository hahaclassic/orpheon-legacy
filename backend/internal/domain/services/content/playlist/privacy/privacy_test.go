package privacy_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/privacy"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestPlaylistPrivacyChangerSuite(t *testing.T) {
	suite.Run(t, &PlaylistPrivacyChangerSuite{})
}

type PrivacyObjectMother struct{}

func (PrivacyObjectMother) Claims(userID uuid.UUID, accessLvl entity.AccessLevel) *entity.Claims {
	return &entity.Claims{UserID: userID, AccessLvl: accessLvl}
}

type PlaylistPrivacyChangerSuite struct {
	suite.Suite

	ctx     context.Context
	service *privacy.PlaylistPrivacyChanger
	policy  *mocks.PlaylistPolicyService
	favs    *mocks.FavoritesDeletionService
	repo    *mocks.PlaylistPrivacyRepository

	objMother *PrivacyObjectMother
}

func (s *PlaylistPrivacyChangerSuite) SetupTest() {
	s.ctx = context.Background()
	s.policy = mocks.NewPlaylistPolicyService(s.T())
	s.favs = mocks.NewFavoritesDeletionService(s.T())
	s.repo = mocks.NewPlaylistPrivacyRepository(s.T())
	s.service = privacy.NewPlaylistPrivacyChanger(s.policy, s.favs, s.repo)
	s.objMother = &PrivacyObjectMother{}
}

func (s *PlaylistPrivacyChangerSuite) TestMakePrivateSuccess() {
	s.SetupTest()
	userID := uuid.New()
	playlistID := uuid.New()

	s.policy.On("CanEdit", s.ctx, mock.Anything, playlistID).Return(nil)
	s.favs.On("GetUsersWithFavoritePlaylist", s.ctx, mock.Anything, playlistID, false).Return([]uuid.UUID{}, nil)
	s.favs.On("DeleteFromAllFavorites", s.ctx, mock.Anything, playlistID, false).Return(nil)
	s.repo.On("UpdatePrivacy", s.ctx, playlistID, true).Return(nil)

	claims := s.objMother.Claims(userID, 0)
	err := s.service.ChangePrivacy(s.ctx, claims, playlistID, true)
	assert.NoError(s.T(), err)

	s.policy.AssertExpectations(s.T())
	s.favs.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistPrivacyChangerSuite) TestMakePublicSuccess() {
	s.SetupTest()
	userID := uuid.New()
	playlistID := uuid.New()

	s.policy.On("CanEdit", s.ctx, mock.Anything, playlistID).Return(nil)
	s.repo.On("UpdatePrivacy", s.ctx, playlistID, false).Return(nil)

	claims := s.objMother.Claims(userID, 0)
	err := s.service.ChangePrivacy(s.ctx, claims, playlistID, false)
	assert.NoError(s.T(), err)

	s.policy.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *PlaylistPrivacyChangerSuite) TestChangePrivacyForbidden() {
	s.SetupTest()
	userID := uuid.New()
	playlistID := uuid.New()

	s.policy.On("CanEdit", s.ctx, mock.Anything, playlistID).Return(commonerr.ErrForbidden)

	claims := s.objMother.Claims(userID, 0)
	err := s.service.ChangePrivacy(s.ctx, claims, playlistID, true)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)

	s.policy.AssertExpectations(s.T())
}

func (s *PlaylistPrivacyChangerSuite) TestGetUsersError() {
	s.SetupTest()
	userID := uuid.New()
	playlistID := uuid.New()

	s.policy.On("CanEdit", s.ctx, mock.Anything, playlistID).Return(nil)
	s.favs.On("GetUsersWithFavoritePlaylist", s.ctx, mock.Anything, playlistID, false).Return(nil, assert.AnError)

	claims := s.objMother.Claims(userID, 0)
	err := s.service.ChangePrivacy(s.ctx, claims, playlistID, true)
	assert.ErrorIs(s.T(), err, playlist.ErrChangePrivacy)

	s.policy.AssertExpectations(s.T())
	s.favs.AssertExpectations(s.T())
}

func (s *PlaylistPrivacyChangerSuite) TestDeleteFavoritesError() {
	s.SetupTest()
	userID := uuid.New()
	playlistID := uuid.New()

	s.policy.On("CanEdit", s.ctx, mock.Anything, playlistID).Return(nil)
	s.favs.On("GetUsersWithFavoritePlaylist", s.ctx, mock.Anything, playlistID, false).Return([]uuid.UUID{}, nil)
	s.favs.On("DeleteFromAllFavorites", s.ctx, mock.Anything, playlistID, false).Return(assert.AnError)

	claims := s.objMother.Claims(userID, 0)
	err := s.service.ChangePrivacy(s.ctx, claims, playlistID, true)
	assert.ErrorIs(s.T(), err, playlist.ErrChangePrivacy)

	s.policy.AssertExpectations(s.T())
	s.favs.AssertExpectations(s.T())
}

func (s *PlaylistPrivacyChangerSuite) TestUpdatePrivacyErrorRollback() {
	s.SetupTest()
	userID := uuid.New()
	playlistID := uuid.New()

	s.policy.On("CanEdit", s.ctx, mock.Anything, playlistID).Return(nil)
	s.favs.On("GetUsersWithFavoritePlaylist", s.ctx, mock.Anything, playlistID, false).Return([]uuid.UUID{}, nil)
	s.favs.On("DeleteFromAllFavorites", s.ctx, mock.Anything, playlistID, false).Return(nil)
	s.favs.On("AddPlaylistToAllFavorites", s.ctx, mock.Anything, []uuid.UUID{}, playlistID).Return(nil)
	s.repo.On("UpdatePrivacy", s.ctx, playlistID, true).Return(assert.AnError)

	claims := s.objMother.Claims(userID, 0)
	err := s.service.ChangePrivacy(s.ctx, claims, playlistID, true)
	assert.ErrorIs(s.T(), err, playlist.ErrChangePrivacy)

	s.policy.AssertExpectations(s.T())
	s.favs.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}
