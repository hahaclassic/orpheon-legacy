package playlist_tracks_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistTracksRepository struct {
	pool *pgxpool.Pool
}

func NewPlaylistTracksRepository(pool *pgxpool.Pool) *PlaylistTracksRepository {
	return &PlaylistTracksRepository{pool: pool}
}

func (r *PlaylistTracksRepository) AddTrackToPlaylist(ctx context.Context, playlistTrack *entity.PlaylistTrack) error {
	const query = `
		WITH max_position AS (
			SELECT COALESCE(MAX(position), 0) + 1 as next_position
			FROM playlist_tracks
			WHERE playlist_id = $1
		)
		INSERT INTO playlist_tracks (playlist_id, track_id, position)
		SELECT $1, $2, next_position
		FROM max_position
		ON CONFLICT (playlist_id, track_id) DO NOTHING
		RETURNING position
	`

	var position int
	err := r.pool.QueryRow(ctx, query, playlistTrack.PlaylistID, playlistTrack.TrackID).Scan(&position)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("track %s already exists in playlist %s", playlistTrack.TrackID, playlistTrack.PlaylistID)
		}
		return fmt.Errorf("add track to playlist: %w", err)
	}

	return nil
}

func (r *PlaylistTracksRepository) DeleteTrackFromPlaylist(ctx context.Context, playlistTrack *entity.PlaylistTrack) error {
	const query = `CALL delete_track_from_playlist($1, $2);`

	_, err := r.pool.Exec(ctx, query, playlistTrack.PlaylistID, playlistTrack.TrackID)
	if err != nil {
		return fmt.Errorf("delete track from playlist: %w", err)
	}

	return nil
}

func (r *PlaylistTracksRepository) DeleteAllTracksFromPlaylist(ctx context.Context, playlistID uuid.UUID) error {
	const query = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = $1
	`

	ct, err := r.pool.Exec(ctx, query, playlistID)
	if err != nil {
		return fmt.Errorf("delete all tracks from playlist: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%w: no tracks found in playlist %s", commonerr.ErrNotFound, playlistID)
	}

	return nil
}

func (r *PlaylistTracksRepository) GetAllPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]*entity.TrackMeta, error) {
	const query = `
		SELECT 
			t.id, t.genre_id, t.name, t.duration, t.explicit,
			t.license_id, t.album_id, t.track_number, t.total_streams
		FROM playlist_tracks pt
		JOIN tracks t ON pt.track_id = t.id
		WHERE pt.playlist_id = $1
		ORDER BY pt.position ASC
	`

	rows, err := r.pool.Query(ctx, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("get all tracks from playlist: %w", err)
	}
	defer rows.Close()

	var tracks []*entity.TrackMeta
	for rows.Next() {
		var track entity.TrackMeta
		if err := rows.Scan(
			&track.ID,
			&track.GenreID,
			&track.Name,
			&track.Duration,
			&track.Explicit,
			&track.LicenseID,
			&track.AlbumID,
			&track.TrackNumber,
			&track.TotalStreams,
		); err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}

func (r *PlaylistTracksRepository) ChangeTrackPosition(ctx context.Context, playlistTrack *entity.PlaylistTrack) error {
	const query = `
		CALL change_track_position($1, $2, $3);
	`

	_, err := r.pool.Exec(ctx, query, playlistTrack.PlaylistID, playlistTrack.TrackID, playlistTrack.Position)
	if err != nil {
		return fmt.Errorf("change track position: %w", err)
	}

	return nil
}
