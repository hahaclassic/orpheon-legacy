package cover_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/cover"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestCoverServiceSuite(t *testing.T) {
	suite.Run(t, new(CoverServiceSuite))
}

// --- Object Mother ---

type CoverObjectMother struct{}

func (CoverObjectMother) DefaultAlbumID() uuid.UUID {
	return uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
}

func (m CoverObjectMother) DefaultCover() *entity.Cover {
	return &entity.Cover{ObjectID: m.DefaultAlbumID()}
}

func (CoverObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (CoverObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

// --- Suite ---

type CoverServiceSuite struct {
	suite.Suite

	ctx       context.Context
	service   *cover.AlbumCoverService
	repo      *mocks.AlbumCoverRepository
	objMother *CoverObjectMother
}

func (s *CoverServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewAlbumCoverRepository(s.T())
	s.service = cover.New(s.repo)
	s.objMother = &CoverObjectMother{}
}

// --- GetCover ---

func (s *CoverServiceSuite) TestGetCover_Success() {
	albumID := s.objMother.DefaultAlbumID()
	expected := s.objMother.DefaultCover()

	s.repo.On("GetCover", s.ctx, albumID).Return(expected, nil)

	result, err := s.service.GetCover(s.ctx, albumID)

	s.NoError(err)
	s.Equal(expected, result)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestGetCover_RepoError() {
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("GetCover", s.ctx, albumID).Return(nil, errors.New("db error"))

	result, err := s.service.GetCover(s.ctx, albumID)

	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, usecase.ErrGetCover)
	s.repo.AssertExpectations(s.T())
}

// --- UploadCover ---

func (s *CoverServiceSuite) TestUploadCover_AdminSuccess() {
	claims := s.objMother.AdminClaims()
	cov := s.objMother.DefaultCover()

	s.repo.On("SaveCover", s.ctx, cov).Return(nil)

	err := s.service.UploadCover(s.ctx, claims, cov)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestUploadCover_AdminRepoError() {
	claims := s.objMother.AdminClaims()
	cov := s.objMother.DefaultCover()

	s.repo.On("SaveCover", s.ctx, cov).Return(errors.New("db error"))

	err := s.service.UploadCover(s.ctx, claims, cov)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrUploadCover)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestUploadCover_NonAdminForbidden() {
	claims := s.objMother.UserClaims()
	cov := s.objMother.DefaultCover()

	err := s.service.UploadCover(s.ctx, claims, cov)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestUploadCover_NilClaimsForbidden() {
	cov := s.objMother.DefaultCover()

	err := s.service.UploadCover(s.ctx, nil, cov)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertExpectations(s.T())
}

// --- DeleteCover ---

func (s *CoverServiceSuite) TestDeleteCover_AdminSuccess() {
	claims := s.objMother.AdminClaims()
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("DeleteCover", s.ctx, albumID).Return(nil)

	err := s.service.DeleteCover(s.ctx, claims, albumID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestDeleteCover_AdminRepoError() {
	claims := s.objMother.AdminClaims()
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("DeleteCover", s.ctx, albumID).Return(errors.New("db error"))

	err := s.service.DeleteCover(s.ctx, claims, albumID)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrDeleteCover)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestDeleteCover_NonAdminForbidden() {
	claims := s.objMother.UserClaims()
	albumID := s.objMother.DefaultAlbumID()

	err := s.service.DeleteCover(s.ctx, claims, albumID)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertExpectations(s.T())
}

func (s *CoverServiceSuite) TestDeleteCover_NilClaimsForbidden() {
	albumID := s.objMother.DefaultAlbumID()

	err := s.service.DeleteCover(s.ctx, nil, albumID)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertExpectations(s.T())
}
