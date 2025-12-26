package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/meta"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

// --- Object Mother ---

type MetaObjectMother struct{}

func (MetaObjectMother) DefaultAlbumID() uuid.UUID {
	return uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
}

func (m MetaObjectMother) DefaultAlbum() *entity.AlbumMeta {
	return &entity.AlbumMeta{ID: m.DefaultAlbumID(), Title: "Test Album"}
}

func (MetaObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (MetaObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

// --- Suite ---

type MetaServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.AlbumRepository
	service   *meta.AlbumService
	objMother *MetaObjectMother
}

func TestMetaServiceSuite(t *testing.T) {
	suite.Run(t, new(MetaServiceSuite))
}

func (s *MetaServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewAlbumRepository(s.T())
	s.service = meta.New(s.repo)
	s.objMother = &MetaObjectMother{}
}

// --- CreateAlbum ---

func (s *MetaServiceSuite) TestCreateAlbum_Success() {
	admin := s.objMother.AdminClaims()
	album := s.objMother.DefaultAlbum()

	s.repo.On("CreateAlbum", s.ctx, mock.AnythingOfType("*entity.AlbumMeta")).Return(nil)

	id, err := s.service.CreateAlbum(s.ctx, admin, album)

	s.NoError(err)
	s.NotEqual(uuid.Nil, id)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestCreateAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	album := s.objMother.DefaultAlbum()

	id, err := s.service.CreateAlbum(s.ctx, user, album)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.Equal(uuid.Nil, id)
	s.repo.AssertNotCalled(s.T(), "CreateAlbum")
}

func (s *MetaServiceSuite) TestCreateAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	album := s.objMother.DefaultAlbum()

	s.repo.On("CreateAlbum", s.ctx, mock.AnythingOfType("*entity.AlbumMeta")).Return(errors.New("db error"))

	id, err := s.service.CreateAlbum(s.ctx, admin, album)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrCreateAlbum)
	s.Equal(uuid.Nil, id)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestCreateAlbum_NilClaims() {
	album := s.objMother.DefaultAlbum()

	id, err := s.service.CreateAlbum(s.ctx, nil, album)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.Equal(uuid.Nil, id)
	s.repo.AssertNotCalled(s.T(), "CreateAlbum")
}

// --- GetAlbum ---

func (s *MetaServiceSuite) TestGetAlbum_Success() {
	id := s.objMother.DefaultAlbumID()
	expected := s.objMother.DefaultAlbum()

	s.repo.On("GetAlbum", s.ctx, id).Return(expected, nil)

	res, err := s.service.GetAlbum(s.ctx, id)

	s.NoError(err)
	s.Equal(expected, res)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestGetAlbum_RepoError() {
	id := s.objMother.DefaultAlbumID()

	s.repo.On("GetAlbum", s.ctx, id).Return(nil, errors.New("repo error"))

	res, err := s.service.GetAlbum(s.ctx, id)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrGetAlbum)
	s.Nil(res)
	s.repo.AssertExpectations(s.T())
}

// --- GetAllAlbums ---

func (s *MetaServiceSuite) TestGetAllAlbums_Success() {
	albums := []*entity.AlbumMeta{s.objMother.DefaultAlbum()}

	s.repo.On("GetAllAlbums", s.ctx).Return(albums, nil)

	res, err := s.service.GetAllAlbums(s.ctx)

	s.NoError(err)
	s.Equal(albums, res)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestGetAllAlbums_RepoError() {
	s.repo.On("GetAllAlbums", s.ctx).Return(nil, errors.New("db error"))

	res, err := s.service.GetAllAlbums(s.ctx)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrGetAllAlbums)
	s.Nil(res)
	s.repo.AssertExpectations(s.T())
}

// --- UpdateAlbum ---

func (s *MetaServiceSuite) TestUpdateAlbum_Success() {
	admin := s.objMother.AdminClaims()
	album := s.objMother.DefaultAlbum()

	s.repo.On("UpdateAlbum", s.ctx, album).Return(nil)

	err := s.service.UpdateAlbum(s.ctx, admin, album)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestUpdateAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	album := s.objMother.DefaultAlbum()

	err := s.service.UpdateAlbum(s.ctx, user, album)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "UpdateAlbum")
}

func (s *MetaServiceSuite) TestUpdateAlbum_NilClaims() {
	album := s.objMother.DefaultAlbum()

	err := s.service.UpdateAlbum(s.ctx, nil, album)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "UpdateAlbum")
}

func (s *MetaServiceSuite) TestUpdateAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	album := s.objMother.DefaultAlbum()

	s.repo.On("UpdateAlbum", s.ctx, album).Return(errors.New("repo error"))

	err := s.service.UpdateAlbum(s.ctx, admin, album)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrUpdateAlbum)
	s.repo.AssertExpectations(s.T())
}

// --- DeleteAlbum ---

func (s *MetaServiceSuite) TestDeleteAlbum_Success() {
	admin := s.objMother.AdminClaims()
	id := s.objMother.DefaultAlbumID()

	s.repo.On("DeleteAlbum", s.ctx, id).Return(nil)

	err := s.service.DeleteAlbum(s.ctx, admin, id)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestDeleteAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	id := s.objMother.DefaultAlbumID()

	s.repo.On("DeleteAlbum", s.ctx, id).Return(errors.New("db error"))

	err := s.service.DeleteAlbum(s.ctx, admin, id)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrDeleteAlbum)
	s.repo.AssertExpectations(s.T())
}

func (s *MetaServiceSuite) TestDeleteAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	id := s.objMother.DefaultAlbumID()

	err := s.service.DeleteAlbum(s.ctx, user, id)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "DeleteAlbum")
}

func (s *MetaServiceSuite) TestDeleteAlbum_NilClaims() {
	id := s.objMother.DefaultAlbumID()

	err := s.service.DeleteAlbum(s.ctx, nil, id)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "DeleteAlbum")
}
