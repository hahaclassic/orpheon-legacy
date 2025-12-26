package search_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/search"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestSearchServiceSuite(t *testing.T) {
	suite.Run(t, &SearchServiceSuite{})
}

type SearchServiceSuite struct {
	suite.Suite
	ctx     context.Context
	repo    *mocks.SearchRepository
	service *search.SearchService
	userID  uuid.UUID
	claims  *entity.Claims
	req     *entity.SearchRequest
}

func (s *SearchServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewSearchRepository(s.T())
	s.service = search.NewSearchService(s.repo)
	s.userID = uuid.New()
	s.claims = &entity.Claims{UserID: s.userID}
	s.req = &entity.SearchRequest{Query: "test"}
}

func (s *SearchServiceSuite) TestSearchTracksSuccess() {
	tracks := []*entity.TrackMeta{{ID: uuid.New(), Name: "Track1"}}
	s.repo.On("SearchTracks", mock.Anything, s.req).Return(tracks, nil)

	res, err := s.service.SearchTracks(s.ctx, s.req)
	assert.Equal(s.T(), tracks, res)
	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchTracksRepoError() {
	s.repo.On("SearchTracks", mock.Anything, s.req).Return(nil, errors.New("repo error"))

	res, err := s.service.SearchTracks(s.ctx, s.req)
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), err, usecase.ErrSearchTracks)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchAlbumsSuccess() {
	albums := []*entity.AlbumMeta{{ID: uuid.New(), Title: "Album1"}}
	s.repo.On("SearchAlbums", mock.Anything, s.req).Return(albums, nil)

	res, err := s.service.SearchAlbums(s.ctx, s.req)
	assert.Equal(s.T(), albums, res)
	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchAlbumsRepoError() {
	s.repo.On("SearchAlbums", mock.Anything, s.req).Return(nil, errors.New("repo error"))

	res, err := s.service.SearchAlbums(s.ctx, s.req)
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), err, usecase.ErrSearchAlbums)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchArtistsSuccess() {
	artists := []*entity.ArtistMeta{{ID: uuid.New(), Name: "Artist1"}}
	s.repo.On("SearchArtists", mock.Anything, s.req).Return(artists, nil)

	res, err := s.service.SearchArtists(s.ctx, s.req)
	assert.Equal(s.T(), artists, res)
	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchArtistsRepoError() {
	s.repo.On("SearchArtists", mock.Anything, s.req).Return(nil, errors.New("repo error"))

	res, err := s.service.SearchArtists(s.ctx, s.req)
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), err, usecase.ErrSearchArtists)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchPlaylistsSuccess() {
	userID := s.userID
	playlist1 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: false}
	playlist2 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: true, OwnerID: userID}
	playlist3 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: true, OwnerID: uuid.New()}

	s.repo.On("SearchPlaylists", mock.Anything, s.req).Return([]*entity.PlaylistMeta{playlist1, playlist2, playlist3}, nil)

	res, err := s.service.SearchPlaylists(s.ctx, s.claims, s.req)
	assert.Equal(s.T(), []*entity.PlaylistMeta{playlist1, playlist2}, res)
	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchPlaylistsRepoError() {
	s.repo.On("SearchPlaylists", mock.Anything, s.req).Return(nil, errors.New("repo error"))

	res, err := s.service.SearchPlaylists(s.ctx, s.claims, s.req)
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), err, usecase.ErrSearchPlaylists)
	s.repo.AssertExpectations(s.T())
}

func (s *SearchServiceSuite) TestSearchPlaylistsNilClaims() {
	playlist1 := &entity.PlaylistMeta{ID: uuid.New(), IsPrivate: false}
	s.repo.On("SearchPlaylists", mock.Anything, s.req).Return([]*entity.PlaylistMeta{playlist1}, nil)

	res, err := s.service.SearchPlaylists(s.ctx, nil, s.req)
	assert.Equal(s.T(), []*entity.PlaylistMeta{playlist1}, res)
	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}
