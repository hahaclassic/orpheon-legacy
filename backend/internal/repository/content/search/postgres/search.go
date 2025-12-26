package search_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchRepository struct {
	db *pgxpool.Pool
}

func NewSearchRepository(db *pgxpool.Pool) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) SearchTracks(ctx context.Context, req *entity.SearchRequest) ([]*entity.TrackMeta, error) {
	query := `
		SELECT t.id, t.genre_id, t.name, t.duration, t.explicit, t.license_id, t.album_id, t.track_number, t.total_streams
		FROM tracks t
		LEFT JOIN artist_tracks at ON t.id = at.track_id
		LEFT JOIN artists ar ON at.artist_id = ar.id
		WHERE true
	`
	args := []any{}
	argIdx := 1

	if req.Query != "" {
		query += fmt.Sprintf(" AND LOWER(t.name) LIKE LOWER($%d)", argIdx)
		args = append(args, fmt.Sprintf("%%%s%%", req.Query))
		argIdx++
	}
	if req.Filters.GenreID != uuid.Nil {
		query += fmt.Sprintf(" AND t.genre_id = $%d", argIdx)
		args = append(args, req.Filters.GenreID)
		argIdx++
	}
	if req.Filters.Country != "" {
		query += fmt.Sprintf(" AND ar.country = $%d", argIdx)
		args = append(args, req.Filters.Country)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY t.name LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*entity.TrackMeta
	for rows.Next() {
		var track entity.TrackMeta
		err := rows.Scan(&track.ID, &track.GenreID, &track.Name, &track.Duration, &track.Explicit, &track.LicenseID, &track.AlbumID, &track.TrackNumber, &track.TotalStreams)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, rows.Err()
}

func (r *SearchRepository) SearchAlbums(ctx context.Context, req *entity.SearchRequest) ([]*entity.AlbumMeta, error) {
	query := `
		SELECT DISTINCT a.id, a.title, a.label, a.license_id, a.release_date
		FROM albums a
		LEFT JOIN tracks t ON a.id = t.album_id
		LEFT JOIN artist_tracks at ON t.id = at.track_id
		LEFT JOIN artists ar ON at.artist_id = ar.id
		WHERE true
	`
	args := []any{}
	argIdx := 1

	if req.Query != "" {
		query += fmt.Sprintf(" AND LOWER(a.title) LIKE LOWER($%d)", argIdx)
		args = append(args, fmt.Sprintf("%%%s%%", req.Query))
		argIdx++
	}
	if req.Filters.GenreID != uuid.Nil {
		query += fmt.Sprintf(" AND t.genre_id = $%d", argIdx)
		args = append(args, req.Filters.GenreID)
		argIdx++
	}
	if req.Filters.Country != "" {
		query += fmt.Sprintf(" AND ar.country = $%d", argIdx)
		args = append(args, req.Filters.Country)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY a.release_date DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search albums: %w", err)
	}
	defer rows.Close()

	var albums []*entity.AlbumMeta
	for rows.Next() {
		var album entity.AlbumMeta
		err := rows.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, &album)
	}

	return albums, rows.Err()
}

func (r *SearchRepository) SearchArtists(ctx context.Context, req *entity.SearchRequest) ([]*entity.ArtistMeta, error) {
	query := `
		SELECT DISTINCT a.id, a.name, a.country, a.description
		FROM artists a
		LEFT JOIN artist_albums aa ON a.id = aa.artist_id
		LEFT JOIN albums al ON aa.album_id = al.id
		LEFT JOIN tracks t ON al.id = t.album_id
		WHERE true
	`
	args := []any{}
	argIdx := 1

	if req.Query != "" {
		query += fmt.Sprintf(" AND LOWER(a.name) LIKE LOWER($%d)", argIdx)
		args = append(args, fmt.Sprintf("%%%s%%", req.Query))
		argIdx++
	}
	if req.Filters.Country != "" {
		query += fmt.Sprintf(" AND a.country = $%d", argIdx)
		args = append(args, req.Filters.Country)
		argIdx++
	}
	if req.Filters.GenreID != uuid.Nil {
		query += fmt.Sprintf(" AND t.genre_id = $%d", argIdx)
		args = append(args, req.Filters.GenreID)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY a.name LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search artists: %w", err)
	}
	defer rows.Close()

	var artists []*entity.ArtistMeta
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Country, &artist.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artist: %w", err)
		}
		artists = append(artists, &artist)
	}

	return artists, rows.Err()
}

func (r *SearchRepository) SearchPlaylists(ctx context.Context, req *entity.SearchRequest) ([]*entity.PlaylistMeta, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.is_private, p.owner_id, p.created_at, p.updated_at, rating
		FROM playlists p
		LEFT JOIN playlist_tracks pt ON p.id = pt.playlist_id
		LEFT JOIN tracks t ON pt.track_id = t.id
		WHERE true
	`
	args := []any{}
	argIdx := 1

	if req.Query != "" {
		query += fmt.Sprintf(" AND LOWER(p.name) LIKE LOWER($%d)", argIdx)
		args = append(args, fmt.Sprintf("%%%s%%", req.Query))
		argIdx++
	}
	if req.Filters.GenreID != uuid.Nil {
		query += fmt.Sprintf(" AND t.genre_id = $%d", argIdx)
		args = append(args, req.Filters.GenreID)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY p.rating DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlists: %w", err)
	}
	defer rows.Close()

	var playlists []*entity.PlaylistMeta
	for rows.Next() {
		var playlist entity.PlaylistMeta
		err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.Description, &playlist.IsPrivate, &playlist.OwnerID,
			&playlist.CreatedAt, &playlist.UpdatedAt, &playlist.Rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlists = append(playlists, &playlist)
	}

	return playlists, rows.Err()
}
