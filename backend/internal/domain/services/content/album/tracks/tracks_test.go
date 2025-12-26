package album_tracks_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestAlbumTracksServiceSuite(t *testing.T) {
	suite.Run(t, new(AlbumTracksServiceSuite))
}

// --- Object Mother ---

type AlbumTracksObjectMother struct{}

func (AlbumTracksObjectMother) DefaultAlbumID() uuid.UUID {
	return uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
}

func (AlbumTracksObjectMother) DefaultTracks() []*entity.TrackMeta {
	return []*entity.TrackMeta{
		{ID: uuid.New(), Name: "Track 1", Duration: 180},
		{ID: uuid.New(), Name: "Track 2", Duration: 200},
	}
}

// --- Suite ---

type AlbumTracksServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.AlbumTrackRepository
	service   *AlbumTrackService
	objMother *AlbumTracksObjectMother
}

func (s *AlbumTracksServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewAlbumTrackRepository(s.T())
	s.service = NewAlbumTrackService(s.repo)
	s.objMother = &AlbumTracksObjectMother{}
}

// --- Tests ---

func (s *AlbumTracksServiceSuite) TestGetAllTracks_Success() {
	albumID := s.objMother.DefaultAlbumID()
	expectedTracks := s.objMother.DefaultTracks()

	s.repo.On("GetAllTracks", s.ctx, albumID).Return(expectedTracks, nil)

	tracks, err := s.service.GetAllTracks(s.ctx, albumID)

	s.NoError(err)
	s.Equal(expectedTracks, tracks)
	s.repo.AssertExpectations(s.T())
}

func (s *AlbumTracksServiceSuite) TestGetAllTracks_RepoError() {
	albumID := s.objMother.DefaultAlbumID()

	s.repo.On("GetAllTracks", s.ctx, albumID).Return(nil, errors.New("db error"))

	tracks, err := s.service.GetAllTracks(s.ctx, albumID)

	s.Error(err)
	s.ErrorIs(err, usecase.ErrGetAllTracks)
	s.Nil(tracks)
	s.repo.AssertExpectations(s.T())
}
