package integration_test

import (
	"context"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	album_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/album/meta/postgres"
	genre_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/genre/meta/postgres"
	license_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/license/postgres"
	playlist_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/meta/postgres"
	playlist_tracks_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/tracks/postgres"
	track_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/meta/postgres"
	user_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/user/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================
//             user playlist flow
//
// setup:
//    1. create genre
//    2. create license
//    3. create album
//    4. create track
//    5. create user
//
// flow:
//    1. create playlist meta
//        1.1 check number of playlists
//    2. add track to playlist
//        2.1 check number of tracks in playlist
//    3. delete track from playlist
//        3.1 check number of tracks in playlist
//    4. delete playlist meta
//        4.1 check number of playlists
// =============================================

func TestUserPlaylistFlow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic occurred: %v\n%s", r, debug.Stack())
		}
	}()

	genreRepo := genre_meta_postgres.NewGenreRepository(pgxPool)
	licenseRepo := license_postgres.NewLicenseRepository(pgxPool)
	albumRepo := album_meta_postgres.NewAlbumRepository(pgxPool)
	trackRepo := track_meta_postgres.NewTrackMetaRepository(pgxPool)
	userRepo := user_postgres.NewUserRepository(pgxPool)
	playlistMetaRepo := playlist_meta_postgres.NewPlaylistMetaRepository(pgxPool)
	playlistTrackRepo := playlist_tracks_postgres.NewPlaylistTracksRepository(pgxPool)

	require.NoError(t, runMigrationsUp(pgxPool))
	defer func() {
		require.NoError(t, runMigrationsDown(pgxPool))
	}()

	ctx := context.Background()

	genre := &entity.Genre{
		ID:    uuid.New(),
		Title: "Rock",
	}

	// setup
	require.NoError(t, genreRepo.Create(ctx, genre)) // 1. create genre

	license := &entity.License{
		ID:          uuid.New(),
		Title:       "CC-BY",
		Description: "Attribution",
	}

	require.NoError(t, licenseRepo.Create(ctx, license)) // 2. create license

	album := &entity.AlbumMeta{
		ID:          uuid.New(),
		Title:       "best album",
		Label:       "label some label",
		LicenseID:   license.ID,
		ReleaseDate: time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
	}
	require.NoError(t, albumRepo.CreateAlbum(ctx, album)) // 3. create album

	track := &entity.TrackMeta{
		ID:          uuid.New(),
		Name:        "new track",
		AlbumID:     album.ID,
		TrackNumber: 0,
		Duration:    120,
		LicenseID:   license.ID,
		GenreID:     genre.ID,
	}
	require.NoError(t, trackRepo.Create(ctx, track)) // 4. create track

	user := &entity.UserInfo{
		ID:               uuid.New(),
		Name:             "test_user",
		AccessLvl:        entity.User,
		BirthDate:        time.Date(2000, 10, 10, 0, 0, 0, 0, time.Local),
		RegistrationDate: time.Now(),
	}
	require.NoError(t, userRepo.CreateUser(ctx, user)) // 5. create user

	// flow
	playlist := &entity.PlaylistMeta{
		ID:      uuid.New(),
		Name:    "My Playlist",
		OwnerID: user.ID,
	}
	require.NoError(t, playlistMetaRepo.Create(ctx, playlist)) // 1. create playlist

	// 1.1 check number of playlists
	playlists, err := playlistMetaRepo.GetByUser(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, playlists, 1, "Only one playlist should be in user's collection!")
	assert.Equal(t, playlist.ID, playlists[0].ID, "IDs should be equal")

	// 2. add track to playlist
	playlistTrack := &entity.PlaylistTrack{
		PlaylistID: playlist.ID,
		TrackID:    track.ID,
	}
	require.NoError(t, playlistTrackRepo.AddTrackToPlaylist(ctx, playlistTrack))

	// 2.1 check number of tracks in playlist
	tracks, err := playlistTrackRepo.GetAllPlaylistTracks(ctx, playlist.ID)
	require.NoError(t, err)

	assert.Len(t, tracks, 1, "Only one track should be in playlist!")
	assert.Equal(t, track.ID, tracks[0].ID, "Track IDs should match")

	// 3. delete track from playlist
	require.NoError(t, playlistTrackRepo.DeleteTrackFromPlaylist(ctx, playlistTrack))

	// 3.1 check number of tracks in playlist
	tracks, err = playlistTrackRepo.GetAllPlaylistTracks(ctx, playlist.ID)
	require.NoError(t, err)
	assert.Len(t, tracks, 0, "Playlist should be empty!")

	// 4 delete playlist
	require.NoError(t, playlistMetaRepo.Delete(ctx, playlist.ID))

	// 4.1 check number of playlists
	playlists, err = playlistMetaRepo.GetByUser(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, playlists, 0, "There should be 0 playlists left")
}
