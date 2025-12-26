package genre_assign

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

// --- Object Mother ---

type GenreAssignObjectMother struct{}

func (GenreAssignObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (GenreAssignObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

func (GenreAssignObjectMother) DefaultIDs() (uuid.UUID, uuid.UUID) {
	return uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		uuid.MustParse("22222222-2222-2222-2222-222222222222")
}

// --- Suite ---

type GenreAssignServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.GenreAssignRepository
	service   *GenreAssignService
	objMother *GenreAssignObjectMother
}

func TestGenreAssignServiceSuite(t *testing.T) {
	suite.Run(t, new(GenreAssignServiceSuite))
}

func (s *GenreAssignServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewGenreAssignRepository(s.T())
	s.service = NewGenreAssignService(s.repo)
	s.objMother = &GenreAssignObjectMother{}
}

// --- AssignGenreToAlbum ---

func (s *GenreAssignServiceSuite) TestAssignGenreToAlbum_Success() {
	admin := s.objMother.AdminClaims()
	genreID, albumID := s.objMother.DefaultIDs()

	s.repo.On("AssignGenreToAlbum", s.ctx, genreID, albumID).Return(nil)

	err := s.service.AssignGenreToAlbum(s.ctx, admin, genreID, albumID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *GenreAssignServiceSuite) TestAssignGenreToAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	genreID, albumID := s.objMother.DefaultIDs()

	err := s.service.AssignGenreToAlbum(s.ctx, user, genreID, albumID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "AssignGenreToAlbum", s.ctx, genreID, albumID)
}

func (s *GenreAssignServiceSuite) TestAssignGenreToAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	genreID, albumID := s.objMother.DefaultIDs()

	s.repo.On("AssignGenreToAlbum", s.ctx, genreID, albumID).Return(errors.New("db fail"))

	err := s.service.AssignGenreToAlbum(s.ctx, admin, genreID, albumID)

	s.ErrorIs(err, usecase.ErrAssignGenreToAlbum)
	s.repo.AssertExpectations(s.T())
}

// --- UnassignGenreFromAlbum ---

func (s *GenreAssignServiceSuite) TestUnassignGenreFromAlbum_Success() {
	admin := s.objMother.AdminClaims()
	genreID, albumID := s.objMother.DefaultIDs()

	s.repo.On("UnassignGenreFromAlbum", s.ctx, genreID, albumID).Return(nil)

	err := s.service.UnassignGenreFromAlbum(s.ctx, admin, genreID, albumID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *GenreAssignServiceSuite) TestUnassignGenreFromAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	genreID, albumID := s.objMother.DefaultIDs()

	err := s.service.UnassignGenreFromAlbum(s.ctx, user, genreID, albumID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "UnassignGenreFromAlbum", s.ctx, genreID, albumID)
}

func (s *GenreAssignServiceSuite) TestUnassignGenreFromAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	genreID, albumID := s.objMother.DefaultIDs()

	s.repo.On("UnassignGenreFromAlbum", s.ctx, genreID, albumID).Return(errors.New("delete fail"))

	err := s.service.UnassignGenreFromAlbum(s.ctx, admin, genreID, albumID)

	s.ErrorIs(err, usecase.ErrUnassignGenreFromAlbum)
	s.repo.AssertExpectations(s.T())
}
