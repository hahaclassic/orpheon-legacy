package cover

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"

	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/suite"
)

type PlaylistCoverObjectMother struct{}

func (PlaylistCoverObjectMother) Claims() *entity.Claims {
	return &entity.Claims{UserID: uuid.New()}
}

func (PlaylistCoverObjectMother) Cover(playlistID uuid.UUID) *entity.Cover {
	return &entity.Cover{ObjectID: playlistID, Data: []byte("image")}
}

type PlaylistCoverServiceSuite struct {
	suite.Suite

	ctx    context.Context
	policy *mocks.PlaylistPolicyService
	repo   *mocks.PlaylistCoverRepository
	svc    *PlaylistCoverService
	mother PlaylistCoverObjectMother
}

func TestPlaylistCoverServiceSuite(t *testing.T) {
	suite.Run(t, new(PlaylistCoverServiceSuite))
}

func (s *PlaylistCoverServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.policy = mocks.NewPlaylistPolicyService(s.T())
	s.repo = mocks.NewPlaylistCoverRepository(s.T())
	s.svc = New(s.repo, s.policy)
	s.mother = PlaylistCoverObjectMother{}
}

// --- GetCover ---

func (s *PlaylistCoverServiceSuite) TestGetCover() {
	claims := s.mother.Claims()
	playlistID := uuid.New()
	cov := s.mother.Cover(playlistID)

	s.Run("success", func() {
		s.SetupTest()
		s.policy.On("CanView", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("GetCover", s.ctx, playlistID).Return(cov, nil)

		result, err := s.svc.GetCover(s.ctx, claims, playlistID)
		s.NoError(err)
		s.Equal(cov, result)
	})

	s.Run("policy denied", func() {
		s.SetupTest()
		s.policy.On("CanView", s.ctx, claims, playlistID).Return(commonerr.ErrForbidden)

		result, err := s.svc.GetCover(s.ctx, claims, playlistID)
		s.Nil(result)
		s.ErrorIs(err, usecase.ErrGetCover)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.policy.On("CanView", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("GetCover", s.ctx, playlistID).Return(nil, errors.New("repo fail"))

		result, err := s.svc.GetCover(s.ctx, claims, playlistID)
		s.Nil(result)
		s.ErrorIs(err, usecase.ErrGetCover)
	})
}

// --- UploadCover ---

func (s *PlaylistCoverServiceSuite) TestUploadCover() {
	claims := s.mother.Claims()
	playlistID := uuid.New()
	cov := s.mother.Cover(playlistID)

	s.Run("success", func() {
		s.SetupTest()
		s.policy.On("CanEdit", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("SaveCover", s.ctx, cov).Return(nil)

		err := s.svc.UploadCover(s.ctx, claims, cov)
		s.NoError(err)
	})

	s.Run("policy denied", func() {
		s.SetupTest()
		s.policy.On("CanEdit", s.ctx, claims, playlistID).Return(commonerr.ErrForbidden)

		err := s.svc.UploadCover(s.ctx, claims, cov)
		s.ErrorIs(err, usecase.ErrUploadCover)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.policy.On("CanEdit", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("SaveCover", s.ctx, cov).Return(errors.New("save fail"))

		err := s.svc.UploadCover(s.ctx, claims, cov)
		s.ErrorIs(err, usecase.ErrUploadCover)
	})
}

// --- DeleteCover ---

func (s *PlaylistCoverServiceSuite) TestDeleteCover() {
	claims := s.mother.Claims()
	playlistID := uuid.New()

	s.Run("success", func() {
		s.SetupTest()
		s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("DeleteCover", s.ctx, playlistID).Return(nil)

		err := s.svc.DeleteCover(s.ctx, claims, playlistID)
		s.NoError(err)
	})

	s.Run("policy denied", func() {
		s.SetupTest()
		s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(commonerr.ErrForbidden)

		err := s.svc.DeleteCover(s.ctx, claims, playlistID)
		s.ErrorIs(err, usecase.ErrDeleteCover)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.policy.On("CanDelete", s.ctx, claims, playlistID).Return(nil)
		s.repo.On("DeleteCover", s.ctx, playlistID).Return(errors.New("delete fail"))

		err := s.svc.DeleteCover(s.ctx, claims, playlistID)
		s.ErrorIs(err, usecase.ErrDeleteCover)
	})
}
