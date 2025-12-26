package track_meta_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type TrackMetaRepository struct {
	pool *pgxpool.Pool
}

func NewTrackMetaRepository(pool *pgxpool.Pool) *TrackMetaRepository {
	return &TrackMetaRepository{pool: pool}
}

func (r *TrackMetaRepository) GetByID(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error) {
	query := `
		SELECT id, genre_id, name, duration, explicit, license_id, album_id, track_number, total_streams
		FROM tracks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, trackID)

	var track entity.TrackMeta
	err := row.Scan(
		&track.ID,
		&track.GenreID,
		&track.Name,
		&track.Duration,
		&track.Explicit,
		&track.LicenseID,
		&track.AlbumID,
		&track.TrackNumber,
		&track.TotalStreams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

func (r *TrackMetaRepository) GetTrackArtists(ctx context.Context, trackID uuid.UUID) ([]*entity.ArtistMeta, error) {
	query := `
		SELECT a.id, a.name, a.description, a.country
		FROM artists
		JOIN artist_tracks ON artists.id = artist_tracks.artist_id
		WHERE artist_tracks.track_id = $1
	`

	rows, err := r.pool.Query(ctx, query, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track artists: %w", err)
	}
	defer rows.Close()

	artists := make([]*entity.ArtistMeta, 0)
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artist: %w", err)
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *TrackMetaRepository) Create(ctx context.Context, track *entity.TrackMeta) error {
	query := `
		WITH max_track_number AS (
			SELECT COALESCE(MAX(track_number), 0) + 1 as next_number
			FROM tracks
			WHERE album_id = $7
		)
		INSERT INTO tracks (id, genre_id, name, duration, explicit, license_id, album_id, track_number, total_streams)
		SELECT $1, $2, $3, $4, $5, $6, $7, next_number, $8
		FROM max_track_number
	`

	_, err := r.pool.Exec(ctx, query,
		track.ID,
		track.GenreID,
		track.Name,
		track.Duration,
		track.Explicit,
		track.LicenseID,
		track.AlbumID,
		0,
	)

	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	return nil
}

func (r *TrackMetaRepository) Update(ctx context.Context, track *entity.TrackMeta) error {
	query := `
		UPDATE tracks
		SET genre_id = $2,
			name = $3,
			duration = $4,
			explicit = $5,
			license_id = $6
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		track.ID,
		track.GenreID,
		track.Name,
		track.Duration,
		track.Explicit,
		track.LicenseID,
	)

	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	return nil
}

func (r *TrackMetaRepository) Delete(ctx context.Context, trackID uuid.UUID) error {
	query := `
		WITH deleted_track AS (
			DELETE FROM tracks
			WHERE id = $1
			RETURNING album_id, track_number
		)
		UPDATE tracks
		SET track_number = track_number - 1
		WHERE album_id = (SELECT album_id FROM deleted_track)
		AND track_number > (SELECT track_number FROM deleted_track)
	`

	_, err := r.pool.Exec(ctx, query, trackID)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	return nil
}

func (r *TrackMetaRepository) IncrementTrackTotalStreams(ctx context.Context, trackID uuid.UUID) error {
	query := `
		UPDATE tracks
		SET total_streams = total_streams + 1
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, trackID)
	if err != nil {
		return fmt.Errorf("failed to increment track total streams: %w", err)
	}

	return nil
}
