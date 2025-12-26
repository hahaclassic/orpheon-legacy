package policy_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/policy"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/suite"
)

func TestPlaylistPolicyServiceSuite(t *testing.T) {
	suite.Run(t, &PlaylistPolicyServiceSuite{})
}

type PolicyObjectMother struct{}

func (PolicyObjectMother) Claims(userID uuid.UUID, accessLvl entity.AccessLevel) *entity.Claims {
	return &entity.Claims{UserID: userID, AccessLvl: accessLvl}
}

func (PolicyObjectMother) AccessMeta(ownerID uuid.UUID, isPrivate bool) *entity.PlaylistAccessMeta {
	return &entity.PlaylistAccessMeta{OwnerID: ownerID, IsPrivate: isPrivate}
}

type PlaylistPolicyServiceSuite struct {
	suite.Suite

	ctx     context.Context
	service *policy.PlaylistPolicyService
	repo    *mocks.PlaylistAccessRepository

	objMother *PolicyObjectMother
}

func (s *PlaylistPolicyServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewPlaylistAccessRepository(s.T())
	s.service = policy.New(s.repo)
	s.objMother = &PolicyObjectMother{}
}

// --- CanView ---

func (s *PlaylistPolicyServiceSuite) TestCanView() {
	userID := uuid.New()
	playlistID := uuid.New()

	s.Run("owner can view private", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(userID, true)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanView(s.ctx, claims, playlistID)
		s.NoError(err)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("public playlist", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(uuid.New(), false)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanView(s.ctx, claims, playlistID)
		s.NoError(err)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("not owner and private", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(uuid.New(), true)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanView(s.ctx, claims, playlistID)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.repo.AssertExpectations(s.T())
	})
}

// --- CanEdit ---

func (s *PlaylistPolicyServiceSuite) TestCanEdit() {
	userID := uuid.New()
	playlistID := uuid.New()

	s.Run("owner can edit", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(userID, true)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanEdit(s.ctx, claims, playlistID)
		s.NoError(err)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("not owner cannot edit", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(uuid.New(), true)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanEdit(s.ctx, claims, playlistID)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.repo.AssertExpectations(s.T())
	})
}

// --- CanDelete ---

func (s *PlaylistPolicyServiceSuite) TestCanDelete() {
	userID := uuid.New()
	playlistID := uuid.New()

	s.Run("owner can delete", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(userID, true)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanDelete(s.ctx, claims, playlistID)
		s.NoError(err)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("admin can delete public playlist", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(uuid.New(), false)
		claims := s.objMother.Claims(uuid.New(), entity.Admin)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanDelete(s.ctx, claims, playlistID)
		s.NoError(err)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("admin cannot delete private playlist", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(uuid.New(), true)
		claims := s.objMother.Claims(uuid.New(), entity.Admin)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanDelete(s.ctx, claims, playlistID)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.repo.AssertExpectations(s.T())
	})

	s.Run("user cannot delete others' playlist", func() {
		s.SetupTest()
		meta := s.objMother.AccessMeta(uuid.New(), false)
		claims := s.objMother.Claims(userID, 0)

		s.repo.On("GetAccessMeta", s.ctx, playlistID).Return(meta, nil)

		err := s.service.CanDelete(s.ctx, claims, playlistID)
		s.ErrorIs(err, commonerr.ErrForbidden)
		s.repo.AssertExpectations(s.T())
	})
}
