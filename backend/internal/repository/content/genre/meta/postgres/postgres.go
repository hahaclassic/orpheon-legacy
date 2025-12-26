package genre_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type GenreRepository struct {
	pool *pgxpool.Pool
}

func NewGenreRepository(pool *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{pool: pool}
}

func (r *GenreRepository) Create(ctx context.Context, genre *entity.Genre) error {
	query := `INSERT INTO genres (id, title) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, genre.ID, genre.Title)
	return err
}

func (r *GenreRepository) GetByID(ctx context.Context, genreID uuid.UUID) (*entity.Genre, error) {
	query := `SELECT id, title FROM genres WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, genreID)

	var g entity.Genre
	err := row.Scan(&g.ID, &g.Title)
	if err != nil {
		return nil, fmt.Errorf("genre not found: %w", err)
	}
	return &g, nil
}

func (r *GenreRepository) GetAll(ctx context.Context) ([]*entity.Genre, error) {
	query := `SELECT id, title FROM genres`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all genres: %w", err)
	}
	defer rows.Close()

	var genres []*entity.Genre
	for rows.Next() {
		var g entity.Genre
		err := rows.Scan(&g.ID, &g.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to scan genre: %w", err)
		}
		genres = append(genres, &g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return genres, nil
}

func (r *GenreRepository) Update(ctx context.Context, genre *entity.Genre) error {
	query := `UPDATE genres SET title = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, genre.Title, genre.ID)
	return err
}

func (r *GenreRepository) Delete(ctx context.Context, genreID uuid.UUID) error {
	query := `DELETE FROM genres WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, genreID)
	return err
}

func (r *GenreRepository) GetByAlbum(ctx context.Context, albumID uuid.UUID) ([]*entity.Genre, error) {
	query := `SELECT g.id, g.title FROM genres g
	JOIN album_genres ag ON g.id = ag.genre_id
	WHERE ag.album_id = $1`

	rows, err := r.pool.Query(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get genres by album: %w", err)
	}
	defer rows.Close()

	var genres []*entity.Genre
	for rows.Next() {
		var g entity.Genre
		err := rows.Scan(&g.ID, &g.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to scan genre: %w", err)
		}
		genres = append(genres, &g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return genres, nil
}
