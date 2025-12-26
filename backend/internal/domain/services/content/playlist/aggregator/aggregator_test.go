package playlist_aggregator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

type PlaylistObjectMother struct{}

func (PlaylistObjectMother) DefaultClaims() *entity.Claims {
	return &entity.Claims{UserID: uuid.New(), AccessLvl: entity.User}
}

func (PlaylistObjectMother) DefaultUser() *entity.UserInfo {
	return &entity.UserInfo{ID: uuid.New(), Name: "user"}
}

func (PlaylistObjectMother) DefaultPlaylistMeta() *entity.PlaylistMeta {
	return &entity.PlaylistMeta{
		ID:          uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "Playlist 1",
		Description: "desc",
		IsPrivate:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Rating:      5,
	}
}

func (PlaylistObjectMother) DefaultTrack() *entity.TrackMeta {
	return &entity.TrackMeta{ID: uuid.New(), Name: "track1"}
}

type PlaylistAggregatorSuite struct {
	suite.Suite

	ctx          context.Context
	metaService  *mocks.PlaylistMetaService
	favService   *mocks.PlaylistFavoriteService
	trackService *mocks.PlaylistTrackService
	userService  *mocks.UserService
	svc          *PlaylistAggregator
	mother       PlaylistObjectMother
}

func TestPlaylistAggregatorSuite(t *testing.T) {
	suite.Run(t, new(PlaylistAggregatorSuite))
}

func (s *PlaylistAggregatorSuite) SetupTest() {
	s.ctx = context.Background()
	s.metaService = mocks.NewPlaylistMetaService(s.T())
	s.favService = mocks.NewPlaylistFavoriteService(s.T())
	s.trackService = mocks.NewPlaylistTrackService(s.T())
	s.userService = mocks.NewUserService(s.T())

	s.svc = NewPlaylistAggregator(
		s.metaService,
		s.favService,
		s.trackService,
		s.userService,
	)
	s.mother = PlaylistObjectMother{}
}

func (s *PlaylistAggregatorSuite) TearDownTest() {
	s.metaService.AssertExpectations(s.T())
	s.favService.AssertExpectations(s.T())
	s.trackService.AssertExpectations(s.T())
	s.userService.AssertExpectations(s.T())
}

// --- GetPlaylistsByIDs ---

func (s *PlaylistAggregatorSuite) TestGetPlaylistsByIDs() {
	claims := s.mother.DefaultClaims()
	meta := s.mother.DefaultPlaylistMeta()
	user := s.mother.DefaultUser()
	tracks := []*entity.TrackMeta{s.mother.DefaultTrack()}

	s.Run("success", func() {
		s.SetupTest()

		s.metaService.On("GetMeta", s.ctx, claims, meta.ID).Return(meta, nil)
		s.userService.On("GetUser", s.ctx, meta.OwnerID).Return(user, nil)
		s.favService.On("IsFavorite", s.ctx, claims, meta.ID).Return(true, nil)
		s.trackService.On("GetAllTracks", s.ctx, claims, meta.ID).Return(tracks, nil)

		result, err := s.svc.GetPlaylistsByIDs(s.ctx, claims, meta.ID)
		s.NoError(err)
		s.Len(result, 1)
		s.Equal(user, result[0].Owner)
		s.Equal(len(tracks), result[0].TracksCount)
		s.True(result[0].IsFavorite)
	})

	s.Run("meta error", func() {
		s.SetupTest()
		s.metaService.On("GetMeta", s.ctx, claims, meta.ID).Return(nil, errors.New("meta error"))
		result, err := s.svc.GetPlaylistsByIDs(s.ctx, claims, meta.ID)
		s.Error(err)
		s.Nil(result)
	})

	s.Run("user error", func() {
		s.SetupTest()
		s.metaService.On("GetMeta", s.ctx, claims, meta.ID).Return(meta, nil)
		s.userService.On("GetUser", s.ctx, meta.OwnerID).Return(nil, errors.New("user error"))

		result, err := s.svc.GetPlaylistsByIDs(s.ctx, claims, meta.ID)
		s.Error(err)
		s.Nil(result)
	})

	s.Run("favorite error", func() {
		s.SetupTest()
		s.metaService.On("GetMeta", s.ctx, claims, meta.ID).Return(meta, nil)
		s.userService.On("GetUser", s.ctx, meta.OwnerID).Return(user, nil)
		s.favService.On("IsFavorite", s.ctx, claims, meta.ID).Return(false, errors.New("fav error"))

		result, err := s.svc.GetPlaylistsByIDs(s.ctx, claims, meta.ID)
		s.Error(err)
		s.Nil(result)
	})

	s.Run("tracks error", func() {
		s.SetupTest()
		s.metaService.On("GetMeta", s.ctx, claims, meta.ID).Return(meta, nil)
		s.userService.On("GetUser", s.ctx, meta.OwnerID).Return(user, nil)
		s.favService.On("IsFavorite", s.ctx, claims, meta.ID).Return(false, nil)
		s.trackService.On("GetAllTracks", s.ctx, claims, meta.ID).Return(nil, errors.New("tracks error"))

		result, err := s.svc.GetPlaylistsByIDs(s.ctx, claims, meta.ID)
		s.Error(err)
		s.Nil(result)
	})
}

// --- GetPlaylists ---

func (s *PlaylistAggregatorSuite) TestGetPlaylists() {
	claims := s.mother.DefaultClaims()
	meta := s.mother.DefaultPlaylistMeta()
	user := s.mother.DefaultUser()
	tracks := []*entity.TrackMeta{s.mother.DefaultTrack()}

	s.Run("success", func() {
		s.SetupTest()

		s.userService.On("GetUser", s.ctx, meta.OwnerID).Return(user, nil)
		s.favService.On("IsFavorite", s.ctx, claims, meta.ID).Return(false, nil)
		s.trackService.On("GetAllTracks", s.ctx, claims, meta.ID).Return(tracks, nil)

		result, err := s.svc.GetPlaylists(s.ctx, claims, meta)
		s.NoError(err)
		s.Len(result, 1)
		s.Equal(user, result[0].Owner)
		s.Equal(len(tracks), result[0].TracksCount)
	})
}
