package integration_test

import (
	"context"
	"math/rand/v2"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	album_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/album/meta/postgres"
	genre_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/genre/meta/postgres"
	license_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/license/postgres"
	audio_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/audio/minio"
	track_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/meta/postgres"
	"github.com/stretchr/testify/require"
)

// ========================
//       track flow
//
// setup:
//    1. create genre
//    2. create license
//
// flow:
//    1. create album meta
//    2. create track meta
//    3. load audio file
//    4. get track meta
//    5. get audio chunk
// ========================

func TestTrackUploadFlow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic occurred: %v\n%s", r, debug.Stack())
		}
	}()

	ctx := context.Background()

	genreRepo := genre_meta_postgres.NewGenreRepository(pgxPool)
	licenseRepo := license_postgres.NewLicenseRepository(pgxPool)
	albumRepo := album_meta_postgres.NewAlbumRepository(pgxPool)
	trackRepo := track_meta_postgres.NewTrackMetaRepository(pgxPool)
	audioRepo, err := audio_minio.NewAudioFileRepository(ctx, minioClient, minioAudioBucketName)
	require.NoError(t, err)

	require.NoError(t, runMigrationsUp(pgxPool))
	defer func() {
		require.NoError(t, runMigrationsDown(pgxPool))
	}()
	defer func() {
		require.NoError(t, clearMinioBucket(ctx, minioClient, minioAudioBucketName))
	}()

	// setup
	// 1.
	genre := &entity.Genre{ID: uuid.New(), Title: "Jazz"}
	require.NoError(t, genreRepo.Create(ctx, genre))

	// 2.
	license := &entity.License{
		Title:       "CC BY 4.0",
		Description: "Creative Commons Attribution 4.0",
	}
	require.NoError(t, licenseRepo.Create(ctx, license))

	// flow
	// 1.
	album := &entity.AlbumMeta{
		ID:          uuid.New(),
		Title:       "best album",
		Label:       "label some label",
		LicenseID:   license.ID,
		ReleaseDate: time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
	}
	require.NoError(t, albumRepo.CreateAlbum(ctx, album))

	// 2.
	track := &entity.TrackMeta{
		ID:          uuid.New(),
		Name:        "new track",
		AlbumID:     album.ID,
		TrackNumber: 1,
		Duration:    120,
		LicenseID:   license.ID,
		GenreID:     genre.ID,
	}
	require.NoError(t, trackRepo.Create(ctx, track))

	// 5.
	audioBytes := make([]byte, 1048576) // 1Mb file
	for i := range audioBytes {
		audioBytes[i] = byte(rand.N(256))
	}
	require.NoError(t, audioRepo.UploadAudioFile(ctx, &entity.AudioChunk{
		TrackID: track.ID,
		Start:   0,
		End:     int64(len(audioBytes) - 1),
		Data:    audioBytes,
	}))

	// 6.
	trackMeta, err := trackRepo.GetByID(ctx, track.ID)
	require.NoError(t, err)
	require.Equal(t, *track, *trackMeta)

	// 7. Получение аудио чанка
	chunk, err := audioRepo.GetAudioChunk(ctx, &entity.AudioChunk{
		TrackID: track.ID,
		Start:   1000000,
		End:     1010240, // 10kb
	})

	require.NoError(t, err)
	require.Equal(t, audioBytes[1000000:1010240], chunk.Data)
}
