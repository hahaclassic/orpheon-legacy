package genre_assign_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GenreAssignRepository struct {
	pool *pgxpool.Pool
}

func NewGenreAssignRepository(pool *pgxpool.Pool) *GenreAssignRepository {
	return &GenreAssignRepository{pool: pool}
}

func (r *GenreAssignRepository) AssignGenreToAlbum(ctx context.Context, genreID uuid.UUID, albumID uuid.UUID) error {
	query := `
		INSERT INTO album_genres (album_id, genre_id)
		VALUES ($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, albumID, genreID)
	if err != nil {
		return fmt.Errorf("assign genre to album: %w", err)
	}
	return nil
}

func (r *GenreAssignRepository) UnassignGenreFromAlbum(ctx context.Context, genreID uuid.UUID, albumID uuid.UUID) error {
	query := `
		DELETE FROM album_genres
		WHERE album_id = $1 AND genre_id = $2
	`
	_, err := r.pool.Exec(ctx, query, albumID, genreID)
	if err != nil {
		return fmt.Errorf("unassign genre from album: %w", err)
	}
	return nil
}
