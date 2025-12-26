package assign_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistAssignRepository struct {
	pool *pgxpool.Pool
}

func NewArtistAssignRepository(pool *pgxpool.Pool) *ArtistAssignRepository {
	return &ArtistAssignRepository{pool: pool}
}

func (r *ArtistAssignRepository) AssignArtistToTrack(ctx context.Context, artistID uuid.UUID, trackID uuid.UUID) error {
	query := `
		INSERT INTO artist_tracks (artist_id, track_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := r.pool.Exec(ctx, query, artistID, trackID)
	if err != nil {
		return fmt.Errorf("assign artist to track: %w", err)
	}

	return nil
}

func (r *ArtistAssignRepository) AssignArtistToAlbum(ctx context.Context, artistID uuid.UUID, albumID uuid.UUID) error {
	query := `
		INSERT INTO artist_albums (artist_id, album_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := r.pool.Exec(ctx, query, artistID, albumID)
	if err != nil {
		return fmt.Errorf("assign artist to album: %w", err)
	}

	return nil
}

func (r *ArtistAssignRepository) GetArtistAlbums(ctx context.Context, artistID uuid.UUID) ([]*entity.AlbumMeta, error) {
	query := `
		SELECT a.id, a.title, a.label, a.license_id, a.release_date
		FROM albums a
		JOIN artist_albums aa ON a.id = aa.album_id
		WHERE aa.artist_id = $1
	`
	rows, err := r.pool.Query(ctx, query, artistID)
	if err != nil {
		return nil, fmt.Errorf("get artist albums: %w", err)
	}
	defer rows.Close()

	var albums []*entity.AlbumMeta
	for rows.Next() {
		var album entity.AlbumMeta
		err := rows.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("get artist albums: %w", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get artist albums: %w", err)
	}

	return albums, nil
}

func (r *ArtistAssignRepository) GetArtistTracks(ctx context.Context, artistID uuid.UUID) ([]*entity.TrackMeta, error) {
	query := `
		SELECT t.id, t.name, t.album_id, t.duration, t.explicit, t.license_id, t.genre_id, t.total_streams, t.track_number
		FROM tracks t
		JOIN artist_tracks at ON t.id = at.track_id
		WHERE at.artist_id = $1 ORDER BY t.total_streams DESC
	`
	rows, err := r.pool.Query(ctx, query, artistID)
	if err != nil {
		return nil, fmt.Errorf("get artist tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*entity.TrackMeta
	for rows.Next() {
		var track entity.TrackMeta
		err := rows.Scan(&track.ID, &track.Name, &track.AlbumID, &track.Duration, &track.Explicit,
			&track.LicenseID, &track.GenreID, &track.TotalStreams, &track.TrackNumber)
		if err != nil {
			return nil, fmt.Errorf("get artist tracks: %w", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get artist tracks: %w", err)
	}

	return tracks, nil
}

func (r *ArtistAssignRepository) GetArtistByAlbum(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error) {
	query := `
		SELECT a.id, a.name, a.description, a.country
		FROM artists a
		JOIN artist_albums aa ON a.id = aa.artist_id
		WHERE aa.album_id = $1
	`
	rows, err := r.pool.Query(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("get artist by album: %w", err)
	}
	defer rows.Close()

	var artists []*entity.ArtistMeta
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
		if err != nil {
			return nil, fmt.Errorf("get artist by album: %w", err)
		}
		artists = append(artists, &artist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get artist by album: %w", err)
	}

	return artists, nil
}

func (r *ArtistAssignRepository) GetArtistByTrack(ctx context.Context, trackID uuid.UUID) ([]*entity.ArtistMeta, error) {
	query := `
		SELECT a.id, a.name, a.description, a.country
		FROM artists a
		JOIN artist_tracks at ON a.id = at.artist_id
		WHERE at.track_id = $1
	`
	rows, err := r.pool.Query(ctx, query, trackID)
	if err != nil {
		return nil, fmt.Errorf("get artist by track: %w", err)
	}
	defer rows.Close()

	var artists []*entity.ArtistMeta
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
		if err != nil {
			return nil, fmt.Errorf("get artist by track: %w", err)
		}
		artists = append(artists, &artist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get artist by track: %w", err)
	}

	return artists, nil
}

func (r *ArtistAssignRepository) UnassignArtistFromTrack(ctx context.Context, artistID uuid.UUID, trackID uuid.UUID) error {
	query := `
		DELETE FROM artist_tracks 
		WHERE artist_id = $1 AND track_id = $2
	`
	_, err := r.pool.Exec(ctx, query, artistID, trackID)
	if err != nil {
		return fmt.Errorf("unassign artist from track: %w", err)
	}
	return nil
}

func (r *ArtistAssignRepository) UnassignArtistFromAlbum(ctx context.Context, artistID uuid.UUID, albumID uuid.UUID) error {
	query := `
		DELETE FROM artist_albums 
		WHERE artist_id = $1 AND album_id = $2
	`
	_, err := r.pool.Exec(ctx, query, artistID, albumID)
	if err != nil {
		return fmt.Errorf("unassign artist from album: %w", err)
	}
	return nil
}
