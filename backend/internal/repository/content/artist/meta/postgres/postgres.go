package artist_meta_postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type artistMetaRepository struct {
	db *pgxpool.Pool
}

func NewArtistMetaRepository(db *pgxpool.Pool) *artistMetaRepository {
	return &artistMetaRepository{db: db}
}

func (r *artistMetaRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ArtistMeta, error) {
	query := `SELECT id, name, description, country FROM artists WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var artist entity.ArtistMeta
	err := row.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func (r *artistMetaRepository) GetAll(ctx context.Context) ([]*entity.ArtistMeta, error) {
	query := `SELECT id, name, description, country FROM artists`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := make([]*entity.ArtistMeta, 0)
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
		if err != nil {
			return nil, err
		}
		artists = append(artists, &artist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (r *artistMetaRepository) GetByAlbum(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error) {
	query := `SELECT a.id, a.name, a.description, a.country FROM artists a JOIN album_artist aa ON a.id = aa.artist_id WHERE aa.album_id = $1`

	rows, err := r.db.Query(ctx, query, albumID)
	if err != nil {
		return nil, err
	}

	artists := make([]*entity.ArtistMeta, 0)
	for rows.Next() {
		var artist entity.ArtistMeta
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Country)
		if err != nil {
			return nil, err
		}
		artists = append(artists, &artist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (r *artistMetaRepository) Create(ctx context.Context, artist *entity.ArtistMeta) error {
	query := `INSERT INTO artists (id, name, description, country) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, artist.ID, artist.Name, artist.Description, artist.Country)
	return err
}

func (r *artistMetaRepository) Update(ctx context.Context, artist *entity.ArtistMeta) error {
	query := `UPDATE artists SET name = $1, description = $2, country = $3 WHERE id = $4`

	cmdTag, err := r.db.Exec(ctx, query, artist.Name, artist.Description, artist.Country, artist.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *artistMetaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM artists WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
