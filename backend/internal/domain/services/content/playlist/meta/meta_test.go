package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/meta"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestPlaylistMetaServiceSuite(t *testing.T) {
	suite.Run(t, &PlaylistMetaServiceSuite{})
}

type PlaylistMetaObjectMother struct{}

func (PlaylistMetaObjectMother) Claims() *entity.Claims {
	return &entity.Claims{UserID: uuid.New()}
}

func (PlaylistMetaObjectMother) PlaylistMeta(name string) *entity.PlaylistMeta {
	return &entity.PlaylistMeta{Name: name}
}

func (PlaylistMetaObjectMother) PlaylistID() uuid.UUID {
	return uuid.New()
}

type PlaylistMetaServiceSuite struct {
	suite.Suite

	ctx        context.Context
	service    *meta.PlaylistMetaService
	repo       *mocks.PlaylistMetaRepository
	policy     *mocks.PlaylistPolicyService
	accessRepo *mocks.PlaylistAccessMetaDeleter
	objMother  *PlaylistMetaObjectMother
}

func (s *PlaylistMetaServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewPlaylistMetaRepository(s.T())
	s.policy = mocks.NewPlaylistPolicyService(s.T())
	s.accessRepo = mocks.NewPlaylistAccessMetaDeleter(s.T())
	s.service = meta.NewPlaylistMetaService(s.repo, s.policy, s.accessRepo)
	s.objMother = &PlaylistMetaObjectMother{}
}

// --- CreateMeta ---
func (s *PlaylistMetaServiceSuite) TestCreateMeta() {
	claims := s.objMother.Claims()

	s.Run("success", func() {
		s.SetupTest()
		playlist := s.objMother.PlaylistMeta("My Playlist")
		s.repo.On("Create", s.ctx, mock.MatchedBy(func(p *entity.PlaylistMeta) bool {
			return p.Name == "My Playlist" && p.OwnerID == claims.UserID
		})).Return(nil)

		err := s.service.CreateMeta(s.ctx, claims, playlist)
		s.NoError(err)

		s.repo.AssertExpectations(s.T())
	})

	s.Run("empty name", func() {
		s.SetupTest()
		playlist := s.objMother.PlaylistMeta("")
		err := s.service.CreateMeta(s.ctx, claims, playlist)
		s.ErrorIs(err, meta.ErrEmptyPlaylistName)
	})
}

// --- GetMeta ---
func (s *PlaylistMetaServiceSuite) TestGetMeta() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()
	playlist := &entity.PlaylistMeta{ID: playlistID, Name: "Meta"}

	s.Run("success", func() {
		s.SetupTest()
		s.policy.On("CanView", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("GetByID", s.ctx, playlistID).Return(playlist, nil)

		got, err := s.service.GetMeta(s.ctx, claims, playlistID)
		s.NoError(err)
		s.Equal(playlist, got)

		s.policy.AssertExpectations(s.T())
		s.repo.AssertExpectations(s.T())
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		s.policy.On("CanView", s.ctx, claims, playlistID).Return(commonerr.ErrForbidden)

		got, err := s.service.GetMeta(s.ctx, claims, playlistID)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.Nil(got)

		s.policy.AssertExpectations(s.T())
	})
}

// --- GetUserAllPlaylistsMeta ---
func (s *PlaylistMetaServiceSuite) TestGetUserAllPlaylistsMeta() {
	userID := uuid.New()
	otherID := uuid.New()

	playlists := []*entity.PlaylistMeta{
		{ID: uuid.New(), OwnerID: userID, Name: "Public1", IsPrivate: false},
		{ID: uuid.New(), OwnerID: userID, Name: "Private", IsPrivate: true},
		{ID: uuid.New(), OwnerID: userID, Name: "Public2", IsPrivate: false},
	}

	s.Run("owner sees all", func() {
		s.SetupTest()
		claims := &entity.Claims{UserID: userID}
		s.repo.On("GetByUser", s.ctx, userID).Return(playlists, nil)

		got, err := s.service.GetUserAllPlaylistsMeta(s.ctx, claims, userID)
		s.NoError(err)
		s.Len(got, 3)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("not owner sees only public", func() {
		s.SetupTest()
		claims := &entity.Claims{UserID: otherID}
		s.repo.On("GetByUser", s.ctx, userID).Return(playlists, nil)

		got, err := s.service.GetUserAllPlaylistsMeta(s.ctx, claims, userID)
		s.NoError(err)
		s.Len(got, 2)
	})
}

// --- UpdateMeta ---
func (s *PlaylistMetaServiceSuite) TestUpdateMeta() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()
	playlist := &entity.PlaylistMeta{ID: playlistID, Name: "Update"}

	s.Run("success", func() {
		s.SetupTest()
		s.policy.On("CanEdit", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("Update", s.ctx, playlist).Return(nil)

		err := s.service.UpdateMeta(s.ctx, claims, playlist)
		s.NoError(err)

		s.policy.AssertExpectations(s.T())
		s.repo.AssertExpectations(s.T())
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		s.policy.On("CanEdit", s.ctx, claims, playlistID).Return(commonerr.ErrForbidden)

		err := s.service.UpdateMeta(s.ctx, claims, playlist)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.policy.AssertExpectations(s.T())
	})
}

// --- DeleteMeta ---
func (s *PlaylistMetaServiceSuite) TestDeleteMeta() {
	claims := s.objMother.Claims()
	playlistID := s.objMother.PlaylistID()

	s.Run("success", func() {
		s.SetupTest()
		s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(nil)
		s.accessRepo.On("DeleteAccessMeta", s.ctx, playlistID).Return(nil)
		s.repo.On("Delete", s.ctx, playlistID).Return(nil)

		err := s.service.DeleteMeta(s.ctx, claims, playlistID)
		s.NoError(err)

		s.policy.AssertExpectations(s.T())
		s.accessRepo.AssertExpectations(s.T())
		s.repo.AssertExpectations(s.T())
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(commonerr.ErrForbidden)

		err := s.service.DeleteMeta(s.ctx, claims, playlistID)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.policy.AssertExpectations(s.T())
	})

	s.Run("access repo error", func() {
		s.SetupTest()
		s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(nil)
		s.accessRepo.On("DeleteAccessMeta", s.ctx, playlistID).Return(errors.New("access error"))

		err := s.service.DeleteMeta(s.ctx, claims, playlistID)
		s.ErrorIs(err, usecase.ErrDeleteMeta)

		s.policy.AssertExpectations(s.T())
		s.accessRepo.AssertExpectations(s.T())
	})
}
