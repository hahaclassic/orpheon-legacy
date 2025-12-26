package assign_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/assign"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

// --- Object Mother ---

type ArtistAssignObjectMother struct{}

func (ArtistAssignObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (ArtistAssignObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

func (ArtistAssignObjectMother) DefaultArtistID() uuid.UUID {
	return uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
}

func (ArtistAssignObjectMother) DefaultAlbumID() uuid.UUID {
	return uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
}

func (ArtistAssignObjectMother) DefaultTrackID() uuid.UUID {
	return uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
}

func (ArtistAssignObjectMother) DefaultAlbums() []*entity.AlbumMeta {
	return []*entity.AlbumMeta{
		{ID: uuid.New(), Title: "Album 1"},
		{ID: uuid.New(), Title: "Album 2"},
	}
}

func (ArtistAssignObjectMother) DefaultTracks() []*entity.TrackMeta {
	return []*entity.TrackMeta{
		{ID: uuid.New(), Name: "Track 1"},
		{ID: uuid.New(), Name: "Track 2"},
	}
}

func (ArtistAssignObjectMother) DefaultArtists() []*entity.ArtistMeta {
	return []*entity.ArtistMeta{
		{ID: uuid.New(), Name: "Artist 1"},
		{ID: uuid.New(), Name: "Artist 2"},
	}
}

// --- Suite ---

type ArtistAssignServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.ArtistAssignRepository
	service   *assign.ArtistAssignService
	objMother *ArtistAssignObjectMother
}

func TestArtistAssignServiceSuite(t *testing.T) {
	suite.Run(t, new(ArtistAssignServiceSuite))
}

func (s *ArtistAssignServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewArtistAssignRepository(s.T())
	s.service = assign.NewArtistAssignService(s.repo)
	s.objMother = &ArtistAssignObjectMother{}
}

// --- AssignArtistToTrack ---

func (s *ArtistAssignServiceSuite) TestAssignArtistToTrack_Success() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	trackID := s.objMother.DefaultTrackID()

	s.repo.On("AssignArtistToTrack", s.ctx, artistID, trackID).Return(nil)

	err := s.service.AssignArtistToTrack(s.ctx, admin, artistID, trackID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestAssignArtistToTrack_Forbidden() {
	user := s.objMother.UserClaims()
	artistID := s.objMother.DefaultArtistID()
	trackID := s.objMother.DefaultTrackID()

	err := s.service.AssignArtistToTrack(s.ctx, user, artistID, trackID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "AssignArtistToTrack", s.ctx, artistID, trackID)
}

func (s *ArtistAssignServiceSuite) TestAssignArtistToTrack_RepoError() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	trackID := s.objMother.DefaultTrackID()

	s.repo.On("AssignArtistToTrack", s.ctx, artistID, trackID).Return(errors.New("db error"))

	err := s.service.AssignArtistToTrack(s.ctx, admin, artistID, trackID)

	s.ErrorIs(err, usecase.ErrAssignArtistOnTrack)
	s.repo.AssertExpectations(s.T())
}

// --- AssignArtistToAlbum ---

func (s *ArtistAssignServiceSuite) TestAssignArtistToAlbum_Success() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("AssignArtistToAlbum", s.ctx, artistID, albumID).Return(nil)

	err := s.service.AssignArtistToAlbum(s.ctx, admin, artistID, albumID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestAssignArtistToAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	artistID := s.objMother.DefaultArtistID()
	albumID := s.objMother.DefaultAlbumID()

	err := s.service.AssignArtistToAlbum(s.ctx, user, artistID, albumID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "AssignArtistToAlbum", s.ctx, artistID, albumID)
}

func (s *ArtistAssignServiceSuite) TestAssignArtistToAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("AssignArtistToAlbum", s.ctx, artistID, albumID).Return(errors.New("repo error"))

	err := s.service.AssignArtistToAlbum(s.ctx, admin, artistID, albumID)

	s.ErrorIs(err, usecase.ErrAssignArtistOnAlbum)
	s.repo.AssertExpectations(s.T())
}

// --- GetArtistAlbums ---

func (s *ArtistAssignServiceSuite) TestGetArtistAlbums_Success() {
	artistID := s.objMother.DefaultArtistID()
	albums := s.objMother.DefaultAlbums()

	s.repo.On("GetArtistAlbums", s.ctx, artistID).Return(albums, nil)

	res, err := s.service.GetArtistAlbums(s.ctx, artistID)

	s.NoError(err)
	s.Equal(albums, res)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestGetArtistAlbums_RepoError() {
	artistID := s.objMother.DefaultArtistID()

	s.repo.On("GetArtistAlbums", s.ctx, artistID).Return(nil, errors.New("db error"))

	res, err := s.service.GetArtistAlbums(s.ctx, artistID)

	s.ErrorIs(err, usecase.ErrGetArtistAlbums)
	s.Nil(res)
	s.repo.AssertExpectations(s.T())
}

// --- GetArtistTracks ---

func (s *ArtistAssignServiceSuite) TestGetArtistTracks_Success() {
	artistID := s.objMother.DefaultArtistID()
	tracks := s.objMother.DefaultTracks()

	s.repo.On("GetArtistTracks", s.ctx, artistID).Return(tracks, nil)

	res, err := s.service.GetArtistTracks(s.ctx, artistID)

	s.NoError(err)
	s.Equal(tracks, res)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestGetArtistTracks_RepoError() {
	artistID := s.objMother.DefaultArtistID()

	s.repo.On("GetArtistTracks", s.ctx, artistID).Return(nil, errors.New("db error"))

	res, err := s.service.GetArtistTracks(s.ctx, artistID)

	s.ErrorIs(err, usecase.ErrGetArtistTracks)
	s.Nil(res)
	s.repo.AssertExpectations(s.T())
}

// --- GetArtistByAlbum ---

func (s *ArtistAssignServiceSuite) TestGetArtistByAlbum_Success() {
	albumID := s.objMother.DefaultAlbumID()
	artists := s.objMother.DefaultArtists()

	s.repo.On("GetArtistByAlbum", s.ctx, albumID).Return(artists, nil)

	res, err := s.service.GetArtistByAlbum(s.ctx, albumID)

	s.NoError(err)
	s.Equal(artists, res)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestGetArtistByAlbum_RepoError() {
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("GetArtistByAlbum", s.ctx, albumID).Return(nil, errors.New("repo error"))

	res, err := s.service.GetArtistByAlbum(s.ctx, albumID)

	s.ErrorIs(err, usecase.ErrGetArtistByAlbum)
	s.Nil(res)
	s.repo.AssertExpectations(s.T())
}

// --- GetArtistByTrack ---

func (s *ArtistAssignServiceSuite) TestGetArtistByTrack_Success() {
	trackID := s.objMother.DefaultTrackID()
	artists := s.objMother.DefaultArtists()

	s.repo.On("GetArtistByTrack", s.ctx, trackID).Return(artists, nil)

	res, err := s.service.GetArtistByTrack(s.ctx, trackID)

	s.NoError(err)
	s.Equal(artists, res)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestGetArtistByTrack_RepoError() {
	trackID := s.objMother.DefaultTrackID()

	s.repo.On("GetArtistByTrack", s.ctx, trackID).Return(nil, errors.New("repo error"))

	res, err := s.service.GetArtistByTrack(s.ctx, trackID)

	s.ErrorIs(err, usecase.ErrGetArtistByTrack)
	s.Nil(res)
	s.repo.AssertExpectations(s.T())
}

// --- UnassignArtistFromTrack ---

func (s *ArtistAssignServiceSuite) TestUnassignArtistFromTrack_Success() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	trackID := s.objMother.DefaultTrackID()

	s.repo.On("UnassignArtistFromTrack", s.ctx, artistID, trackID).Return(nil)

	err := s.service.UnassignArtistFromTrack(s.ctx, admin, artistID, trackID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestUnassignArtistFromTrack_Forbidden() {
	user := s.objMother.UserClaims()
	artistID := s.objMother.DefaultArtistID()
	trackID := s.objMother.DefaultTrackID()

	err := s.service.UnassignArtistFromTrack(s.ctx, user, artistID, trackID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "UnassignArtistFromTrack", s.ctx, artistID, trackID)
}

func (s *ArtistAssignServiceSuite) TestUnassignArtistFromTrack_RepoError() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	trackID := s.objMother.DefaultTrackID()

	s.repo.On("UnassignArtistFromTrack", s.ctx, artistID, trackID).Return(errors.New("db error"))

	err := s.service.UnassignArtistFromTrack(s.ctx, admin, artistID, trackID)

	s.ErrorIs(err, usecase.ErrUnassignArtistFromTrack)
	s.repo.AssertExpectations(s.T())
}

// --- UnassignArtistFromAlbum ---

func (s *ArtistAssignServiceSuite) TestUnassignArtistFromAlbum_Success() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("UnassignArtistFromAlbum", s.ctx, artistID, albumID).Return(nil)

	err := s.service.UnassignArtistFromAlbum(s.ctx, admin, artistID, albumID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistAssignServiceSuite) TestUnassignArtistFromAlbum_Forbidden() {
	user := s.objMother.UserClaims()
	artistID := s.objMother.DefaultArtistID()
	albumID := s.objMother.DefaultAlbumID()

	err := s.service.UnassignArtistFromAlbum(s.ctx, user, artistID, albumID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "UnassignArtistFromAlbum", s.ctx, artistID, albumID)
}

func (s *ArtistAssignServiceSuite) TestUnassignArtistFromAlbum_RepoError() {
	admin := s.objMother.AdminClaims()
	artistID := s.objMother.DefaultArtistID()
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("UnassignArtistFromAlbum", s.ctx, artistID, albumID).Return(errors.New("db error"))

	err := s.service.UnassignArtistFromAlbum(s.ctx, admin, artistID, albumID)

	s.ErrorIs(err, usecase.ErrUnassignArtistFromAlbum)
	s.repo.AssertExpectations(s.T())
}
