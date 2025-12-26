package favorites_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/favorites"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/suite"
)

func TestPlaylistFavoriteServiceSuite(t *testing.T) {
	suite.Run(t, &PlaylistFavoriteServiceSuite{})
}

type PlaylistFavoriteObjectMother struct{}

func (PlaylistFavoriteObjectMother) Claims() *entity.Claims {
	return &entity.Claims{UserID: uuid.New()}
}

func (PlaylistFavoriteObjectMother) PlaylistID() uuid.UUID {
	return uuid.New()
}

func (PlaylistFavoriteObjectMother) UserIDs() []uuid.UUID {
	return []uuid.UUID{uuid.New(), uuid.New()}
}

type PlaylistFavoriteServiceSuite struct {
	suite.Suite

	ctx     context.Context
	service *favorites.PlaylistFavoriteService

	repo   *mocks.PlaylistFavoriteRepository
	policy *mocks.PlaylistPolicyService

	objMother *PlaylistFavoriteObjectMother
}

func (s *PlaylistFavoriteServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewPlaylistFavoriteRepository(s.T())
	s.policy = mocks.NewPlaylistPolicyService(s.T())
	s.service = favorites.NewPlaylistFavoriteService(s.repo, s.policy)
	s.objMother = &PlaylistFavoriteObjectMother{}
}

// --- AddToUserFavorites ---
func (s *PlaylistFavoriteServiceSuite) TestAddToUserFavorites() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{"success", nil, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, commonerr.ErrForbidden},
		{"repo error", nil, errors.New("repo error"), usecase.ErrAddToUserFavorites},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.policy.On("CanView", s.ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				s.repo.On("AddToFavorites", s.ctx, claims.UserID, playlistID).Return(tt.repoErr)
			}

			err := s.service.AddToUserFavorites(s.ctx, claims, playlistID)
			s.ErrorIs(err, tt.wantErr)

			s.policy.AssertExpectations(s.T())
			s.repo.AssertExpectations(s.T())
		})
	}
}

// --- GetUserFavorites ---
func (s *PlaylistFavoriteServiceSuite) TestGetUserFavorites() {
	claims := s.objMother.Claims()
	mockResult := []*entity.PlaylistMeta{{ID: uuid.New()}}

	tests := []struct {
		name    string
		repoRes []*entity.PlaylistMeta
		repoErr error
		wantErr error
	}{
		{"success", mockResult, nil, nil},
		{"repo error", nil, errors.New("repo error"), usecase.ErrGetUserFavorites},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.repo.On("GetUserFavorites", s.ctx, claims.UserID).Return(tt.repoRes, tt.repoErr)

			res, err := s.service.GetUserFavorites(s.ctx, claims)
			s.Equal(tt.repoRes, res)
			s.ErrorIs(err, tt.wantErr)

			s.repo.AssertExpectations(s.T())
		})
	}
}

// --- DeleteFromUserFavorites ---
func (s *PlaylistFavoriteServiceSuite) TestDeleteFromUserFavorites() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()

	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{"success", nil, nil},
		{"repo error", errors.New("repo error"), usecase.ErrDeleteFromUserFavorites},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.repo.On("DeleteFromUserFavorites", s.ctx, claims.UserID, playlistID).Return(tt.repoErr)

			err := s.service.DeleteFromUserFavorites(s.ctx, claims, playlistID)
			s.ErrorIs(err, tt.wantErr)

			s.repo.AssertExpectations(s.T())
		})
	}
}

// --- GetUsersWithFavoritePlaylist ---
func (s *PlaylistFavoriteServiceSuite) TestGetUsersWithFavoritePlaylist() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()
	userList := s.objMother.UserIDs()

	tests := []struct {
		name      string
		policyErr error
		repoRes   []uuid.UUID
		repoErr   error
		wantErr   error
	}{
		{"success", nil, userList, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, nil, commonerr.ErrForbidden},
		{"repo error", nil, nil, errors.New("repo error"), usecase.ErrGetUsersWithFavoritePlaylist},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.policy.On("CanView", s.ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				s.repo.On("GetUsersWithFavoritePlaylist", s.ctx, playlistID, false).Return(tt.repoRes, tt.repoErr)
			}

			res, err := s.service.GetUsersWithFavoritePlaylist(s.ctx, claims, playlistID, true)
			s.Equal(tt.repoRes, res)
			s.ErrorIs(err, tt.wantErr)

			s.policy.AssertExpectations(s.T())
			s.repo.AssertExpectations(s.T())
		})
	}
}

// --- DeleteFromAllFavorites ---
func (s *PlaylistFavoriteServiceSuite) TestDeleteFromAllFavorites() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{"success", nil, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, commonerr.ErrForbidden},
		{"repo error", nil, errors.New("repo error"), usecase.ErrDeleteFromAllFavorites},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				s.repo.On("DeleteFromAllFavorites", s.ctx, playlistID, true).Return(tt.repoErr)
			}

			err := s.service.DeleteFromAllFavorites(s.ctx, claims, playlistID, true)
			s.ErrorIs(err, tt.wantErr)

			s.policy.AssertExpectations(s.T())
			s.repo.AssertExpectations(s.T())
		})
	}
}

// --- AddPlaylistToAllFavorites ---
func (s *PlaylistFavoriteServiceSuite) TestAddPlaylistToAllFavorites() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()
	userIDs := s.objMother.UserIDs()

	tests := []struct {
		name      string
		policyErr error
		repoErr   error
		wantErr   error
	}{
		{"success", nil, nil, nil},
		{"policy error", commonerr.ErrForbidden, nil, commonerr.ErrForbidden},
		{"repo error", nil, errors.New("repo error"), usecase.ErrAddPlaylistToAllFavorites},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(tt.policyErr)
			if tt.policyErr == nil {
				s.repo.On("RestoreAllFavorites", s.ctx, userIDs, playlistID).Return(tt.repoErr)
			}

			err := s.service.AddPlaylistToAllFavorites(s.ctx, claims, userIDs, playlistID)
			s.ErrorIs(err, tt.wantErr)

			s.policy.AssertExpectations(s.T())
			s.repo.AssertExpectations(s.T())
		})
	}
}
