package deleter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/deleter"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

type PlaylistDeleterObjectMother struct{}

func (PlaylistDeleterObjectMother) Claims() *entity.Claims { return &entity.Claims{UserID: uuid.New()} }
func (PlaylistDeleterObjectMother) PlaylistID() uuid.UUID  { return uuid.New() }
func (PlaylistDeleterObjectMother) TrackMeta() []*entity.TrackMeta {
	return []*entity.TrackMeta{{ID: uuid.New()}}
}
func (PlaylistDeleterObjectMother) UserIDs() []uuid.UUID { return []uuid.UUID{uuid.New()} }
func (PlaylistDeleterObjectMother) Cover() *entity.Cover {
	return &entity.Cover{ObjectID: uuid.New(), Data: []byte("image")}
}

type PlaylistDeleterSuite struct {
	suite.Suite
	ctx    context.Context
	mother PlaylistDeleterObjectMother

	meta     *mocks.MetaDeletionService
	tracks   *mocks.TrackDeletionService
	fav      *mocks.FavoritesDeletionService
	coverSvc *mocks.PlaylistCoverDeletionService
}

func TestPlaylistDeleterSuite(t *testing.T) {
	suite.Run(t, new(PlaylistDeleterSuite))
}

func (s *PlaylistDeleterSuite) SetupTest() {
	s.ctx = context.Background()

	s.meta = mocks.NewMetaDeletionService(s.T())
	s.tracks = mocks.NewTrackDeletionService(s.T())
	s.fav = mocks.NewFavoritesDeletionService(s.T())
	s.coverSvc = mocks.NewPlaylistCoverDeletionService(s.T())
}

// --- тесты ---
func (s *PlaylistDeleterSuite) TestDeletePlaylistSuccess() {
	s.SetupTest()
	claims := s.mother.Claims()
	playlistID := s.mother.PlaylistID()
	tracks := s.mother.TrackMeta()
	userIDs := s.mother.UserIDs()
	cover := s.mother.Cover()

	// Настройка моков для успешного выполнения
	s.meta.On("DeleteMeta", s.ctx, claims, playlistID).Return(nil)
	s.tracks.On("GetAllTracks", s.ctx, claims, playlistID).Return(tracks, nil)
	s.tracks.On("DeleteAllTracks", s.ctx, claims, playlistID).Return(nil)
	s.fav.On("GetUsersWithFavoritePlaylist", s.ctx, claims, playlistID, true).Return(userIDs, nil)
	s.fav.On("DeleteFromAllFavorites", s.ctx, claims, playlistID, true).Return(nil)
	s.coverSvc.On("GetCover", s.ctx, claims, playlistID).Return(cover, nil)
	s.coverSvc.On("DeleteCover", s.ctx, claims, playlistID).Return(nil)

	svc := deleter.New(
		deleter.WithMetaDeletion(s.meta),
		deleter.WithTracksDeletion(s.tracks),
		deleter.WithFavoritesDeletion(s.fav),
		deleter.WithCoverDeletion(s.coverSvc),
	)

	err := svc.DeletePlaylist(s.ctx, claims, playlistID)
	assert.NoError(s.T(), err)
}

func (s *PlaylistDeleterSuite) TestDeletePlaylistFavoritesError() {
	s.SetupTest()
	claims := s.mother.Claims()
	playlistID := s.mother.PlaylistID()

	// Ошибка на этапе удаления фаворитов
	s.fav.On("GetUsersWithFavoritePlaylist", s.ctx, claims, playlistID, true).Return(nil, errors.New("fail"))

	svc := deleter.New(
		deleter.WithMetaDeletion(s.meta),
		deleter.WithTracksDeletion(s.tracks),
		deleter.WithFavoritesDeletion(s.fav),
		deleter.WithCoverDeletion(s.coverSvc),
	)

	err := svc.DeletePlaylist(s.ctx, claims, playlistID)
	assert.ErrorIs(s.T(), err, usecase.ErrDeletePlaylist)
}

func (s *PlaylistDeleterSuite) TestDeletePlaylistCoverError() {
	s.SetupTest()
	claims := s.mother.Claims()
	playlistID := s.mother.PlaylistID()
	userIDs := s.mother.UserIDs()

	s.fav.On("GetUsersWithFavoritePlaylist", s.ctx, claims, playlistID, true).Return(userIDs, nil)
	s.fav.On("DeleteFromAllFavorites", s.ctx, claims, playlistID, true).Return(nil)
	s.fav.On("AddPlaylistToAllFavorites", s.ctx, claims, userIDs, playlistID).Return(nil)
	s.coverSvc.On("GetCover", s.ctx, claims, playlistID).Return(nil, errors.New("cover fail"))

	svc := deleter.New(
		deleter.WithMetaDeletion(s.meta),
		deleter.WithFavoritesDeletion(s.fav),
		deleter.WithCoverDeletion(s.coverSvc),
	)

	err := svc.DeletePlaylist(s.ctx, claims, playlistID)
	assert.ErrorIs(s.T(), err, usecase.ErrDeletePlaylist)
}

func (s *PlaylistDeleterSuite) TestDeletePlaylistTrackError() {
	s.SetupTest()
	claims := s.mother.Claims()
	playlistID := s.mother.PlaylistID()
	userIDs := s.mother.UserIDs()
	cover := s.mother.Cover()

	s.fav.On("GetUsersWithFavoritePlaylist", s.ctx, claims, playlistID, true).Return(userIDs, nil)
	s.fav.On("DeleteFromAllFavorites", s.ctx, claims, playlistID, true).Return(nil)
	s.fav.On("AddPlaylistToAllFavorites", s.ctx, claims, userIDs, playlistID).Return(nil)
	s.coverSvc.On("GetCover", s.ctx, claims, playlistID).Return(cover, nil)
	s.coverSvc.On("DeleteCover", s.ctx, claims, playlistID).Return(nil)
	s.coverSvc.On("UploadCover", s.ctx, claims, cover).Return(nil)

	s.tracks.On("GetAllTracks", s.ctx, claims, playlistID).Return(nil, errors.New("track fail"))

	svc := deleter.New(
		deleter.WithTracksDeletion(s.tracks),
		deleter.WithFavoritesDeletion(s.fav),
		deleter.WithCoverDeletion(s.coverSvc),
	)

	err := svc.DeletePlaylist(s.ctx, claims, playlistID)
	assert.ErrorIs(s.T(), err, usecase.ErrDeletePlaylist)
}

func (s *PlaylistDeleterSuite) TestDeletePlaylistMetaError() {
	s.SetupTest()
	claims := s.mother.Claims()
	playlistID := s.mother.PlaylistID()
	tracks := s.mother.TrackMeta()
	userIDs := s.mother.UserIDs()
	cover := s.mother.Cover()

	s.fav.On("GetUsersWithFavoritePlaylist", s.ctx, claims, playlistID, true).Return(userIDs, nil)
	s.fav.On("DeleteFromAllFavorites", s.ctx, claims, playlistID, true).Return(nil)
	s.fav.On("AddPlaylistToAllFavorites", s.ctx, claims, userIDs, playlistID).Return(nil)

	s.coverSvc.On("GetCover", s.ctx, claims, playlistID).Return(cover, nil)
	s.coverSvc.On("DeleteCover", s.ctx, claims, playlistID).Return(nil)
	s.coverSvc.On("UploadCover", s.ctx, claims, cover).Return(nil)

	s.tracks.On("GetAllTracks", s.ctx, claims, playlistID).Return(tracks, nil)
	s.tracks.On("DeleteAllTracks", s.ctx, claims, playlistID).Return(nil)
	s.tracks.On("RestoreAllTracks", s.ctx, claims, playlistID, []uuid.UUID{tracks[0].ID}).Return(nil)

	s.meta.On("DeleteMeta", s.ctx, claims, playlistID).Return(errors.New("meta fail"))

	svc := deleter.New(
		deleter.WithMetaDeletion(s.meta),
		deleter.WithTracksDeletion(s.tracks),
		deleter.WithFavoritesDeletion(s.fav),
		deleter.WithCoverDeletion(s.coverSvc),
	)

	err := svc.DeletePlaylist(s.ctx, claims, playlistID)
	assert.ErrorIs(s.T(), err, usecase.ErrDeletePlaylist)
}
