package album_tracks_postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumTrackRepository struct {
	pool *pgxpool.Pool
}

func NewAlbumTrackRepository(pool *pgxpool.Pool) *AlbumTrackRepository {
	return &AlbumTrackRepository{pool: pool}
}

func (r *AlbumTrackRepository) GetAllTracks(ctx context.Context, albumID uuid.UUID) ([]*entity.TrackMeta, error) {
	query := `
		SELECT t.id, t.name, t.duration, t.explicit, t.license_id, t.album_id,
			   t.track_number, t.total_streams, t.genre_id
		FROM tracks t WHERE t.album_id = $1 ORDER BY t.track_number ASC
	`
	rows, err := r.pool.Query(ctx, query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make([]*entity.TrackMeta, 0)
	for rows.Next() {
		var track entity.TrackMeta
		if err := rows.Scan(&track.ID, &track.Name, &track.Duration,
			&track.Explicit, &track.LicenseID, &track.AlbumID,
			&track.TrackNumber, &track.TotalStreams, &track.GenreID); err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}
