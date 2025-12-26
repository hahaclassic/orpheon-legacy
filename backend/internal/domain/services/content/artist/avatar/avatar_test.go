package avatar_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/avatar"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

// --- Object Mother ---

type ArtistAvatarObjectMother struct{}

func (ArtistAvatarObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (ArtistAvatarObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

func (ArtistAvatarObjectMother) DefaultArtistID() uuid.UUID {
	return uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
}

func (ArtistAvatarObjectMother) DefaultCover() *entity.Cover {
	return &entity.Cover{
		ObjectID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		Data:     []byte{0xFF, 0xD8, 0xFF}, // JPEG header bytes
	}
}

// --- Suite ---

type ArtistAvatarServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.ArtistAvatarRepository
	service   *avatar.ArtistCoverService
	objMother *ArtistAvatarObjectMother
}

func TestArtistAvatarServiceSuite(t *testing.T) {
	suite.Run(t, new(ArtistAvatarServiceSuite))
}

func (s *ArtistAvatarServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewArtistAvatarRepository(s.T())
	s.service = avatar.NewArtistCoverService(s.repo)
	s.objMother = &ArtistAvatarObjectMother{}
}

// --- GetCover ---

func (s *ArtistAvatarServiceSuite) TestGetCover_Success() {
	artistID := s.objMother.DefaultArtistID()
	expected := s.objMother.DefaultCover()

	s.repo.On("GetCover", s.ctx, artistID).Return(expected, nil)

	got, err := s.service.GetCover(s.ctx, artistID)

	s.NoError(err)
	s.Equal(expected, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAvatarServiceSuite) TestGetCover_RepoError() {
	artistID := s.objMother.DefaultArtistID()

	s.repo.On("GetCover", s.ctx, artistID).Return(nil, errors.New("repo error"))

	got, err := s.service.GetCover(s.ctx, artistID)

	s.ErrorIs(err, usecase.ErrGetAvatar)
	s.Nil(got)
	s.repo.AssertExpectations(s.T())
}

// --- UploadCover ---

func (s *ArtistAvatarServiceSuite) TestUploadCover_Success() {
	admin := s.objMother.AdminClaims()
	cover := s.objMother.DefaultCover()

	s.repo.On("SaveCover", s.ctx, cover).Return(nil)

	err := s.service.UploadCover(s.ctx, admin, cover)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAvatarServiceSuite) TestUploadCover_Forbidden() {
	user := s.objMother.UserClaims()
	cover := s.objMother.DefaultCover()

	err := s.service.UploadCover(s.ctx, user, cover)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "SaveCover", s.ctx, cover)
}

func (s *ArtistAvatarServiceSuite) TestUploadCover_RepoError() {
	admin := s.objMother.AdminClaims()
	cover := s.objMother.DefaultCover()

	s.repo.On("SaveCover", s.ctx, cover).Return(errors.New("db error"))

	err := s.service.UploadCover(s.ctx, admin, cover)

	s.ErrorIs(err, usecase.ErrUploadAvatar)
	s.repo.AssertExpectations(s.T())
}

// --- DeleteCover ---

func (s *ArtistAvatarServiceSuite) TestDeleteCover_Success() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()

	s.repo.On("DeleteCover", s.ctx, artistID).Return(nil)

	err := s.service.DeleteCover(s.ctx, admin, artistID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAvatarServiceSuite) TestDeleteCover_Forbidden() {
	user := s.objMother.UserClaims()
	artistID := s.objMother.DefaultArtistID()

	err := s.service.DeleteCover(s.ctx, user, artistID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "DeleteCover", s.ctx, artistID)
}

func (s *ArtistAvatarServiceSuite) TestDeleteCover_RepoError() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()

	s.repo.On("DeleteCover", s.ctx, artistID).Return(errors.New("repo error"))

	err := s.service.DeleteCover(s.ctx, admin, artistID)

	s.ErrorIs(err, usecase.ErrDeleteAvatar)
	s.repo.AssertExpectations(s.T())
}
