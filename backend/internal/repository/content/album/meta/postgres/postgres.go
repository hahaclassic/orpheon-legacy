package album_meta_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumRepository struct {
	pool *pgxpool.Pool
}

func NewAlbumRepository(pool *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{pool: pool}
}

func (r *AlbumRepository) CreateAlbum(ctx context.Context, album *entity.AlbumMeta) error {
	query := `
		INSERT INTO albums (id, title, label, license_id, release_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, album.ID, album.Title, album.Label, album.LicenseID, album.ReleaseDate)
	if err != nil {
		return fmt.Errorf("create album: %w", err)
	}
	return nil
}

func (r *AlbumRepository) GetAlbum(ctx context.Context, id uuid.UUID) (*entity.AlbumMeta, error) {
	query := `
		SELECT id, title, label, license_id, release_date
		FROM albums
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var album entity.AlbumMeta
	err := row.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("get album: %w", err)
	}
	return &album, nil
}

func (r *AlbumRepository) GetAllAlbums(ctx context.Context) ([]*entity.AlbumMeta, error) {
	query := `
		SELECT id, title, label, license_id, release_date
		FROM albums
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all albums: %w", err)
	}
	defer rows.Close()

	albums := make([]*entity.AlbumMeta, 0)
	for rows.Next() {
		var album entity.AlbumMeta
		err := rows.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("get all albums: %w", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all albums: %w", err)
	}

	return albums, nil
}

func (r *AlbumRepository) GetAlbumArtists(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error) {
	query := `
		SELECT a.id, a.name, a.description, a.country
		FROM artists a
		JOIN artist_albums aa ON a.id = aa.artist_id
		WHERE aa.album_id = $1
	`
	rows, err := r.pool.Query(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("get album artists: %w", err)
	}
	defer rows.Close()

	artists := make([]*entity.ArtistMeta, 0)
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
		if err != nil {
			return nil, fmt.Errorf("get album artists: %w", err)
		}
		artists = append(artists, &artist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get album artists: %w", err)
	}

	return artists, nil
}

func (r *AlbumRepository) GetAlbumGenres(ctx context.Context, albumID uuid.UUID) ([]*entity.Genre, error) {
	query := `
		SELECT g.id, g.title
		FROM genres g
		JOIN album_genres ag ON g.id = ag.genre_id
		WHERE ag.album_id = $1
	`
	rows, err := r.pool.Query(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("get album genres: %w", err)
	}
	defer rows.Close()

	genres := make([]*entity.Genre, 0)
	for rows.Next() {
		var genre entity.Genre
		err := rows.Scan(&genre.ID, &genre.Title)
		if err != nil {
			return nil, fmt.Errorf("get album genres: %w", err)
		}
		genres = append(genres, &genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get album genres: %w", err)
	}

	return genres, nil
}

// func (r *AlbumRepository) GetAlbumByArtist(ctx context.Context, artistID uuid.UUID) ([]*entity.AlbumMeta, error) {
// 	query := `
// 		SELECT id, title, label, license_id, release_date
// 		FROM albums
// 		WHERE artist_id = $1
// 	`
// 	rows, err := r.pool.Query(ctx, query, artistID)
// 	if err != nil {
// 		return nil, fmt.Errorf("get album by artist: %w", err)
// 	}
// 	defer rows.Close()

// 	albums := make([]*entity.AlbumMeta, 0)
// 	for rows.Next() {
// 		var album entity.AlbumMeta
// 		err := rows.Scan(&album.ID, &album.Title, &album.Label, &album.LicenseID, &album.ReleaseDate)
// 		if err != nil {
// 			return nil, fmt.Errorf("get album by artist: %w", err)
// 		}
// 		albums = append(albums, &album)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("get album by artist: %w", err)
// 	}

// 	return albums, nil
// }

func (r *AlbumRepository) UpdateAlbum(ctx context.Context, album *entity.AlbumMeta) error {
	query := `
		UPDATE albums
		SET title = $1, label = $2, release_date = $3
		WHERE id = $4
	`
	cmd, err := r.pool.Exec(ctx, query, album.Title, album.Label, album.ReleaseDate, album.ID)
	if err != nil {
		return fmt.Errorf("update album: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("update album: no rows affected")
	}
	return nil
}

func (r *AlbumRepository) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM albums
		WHERE id = $1
	`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete album: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("delete album: no rows affected")
	}
	return nil
}
