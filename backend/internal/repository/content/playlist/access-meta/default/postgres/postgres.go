package access_meta_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistAccessRepository struct {
	pool *pgxpool.Pool
}

func NewPlaylistAccessRepository(pool *pgxpool.Pool) *PlaylistAccessRepository {
	return &PlaylistAccessRepository{pool: pool}
}

func (r *PlaylistAccessRepository) GetAccessMeta(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error) {
	const query = `
		SELECT owner_id, is_private
		FROM playlists
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, playlistID)

	var meta entity.PlaylistAccessMeta
	err := row.Scan(&meta.OwnerID, &meta.IsPrivate)
	if err != nil {
		return nil, fmt.Errorf("get access meta: %w", err)
	}

	return &meta, nil
}

func (r *PlaylistAccessRepository) UpdatePrivacy(ctx context.Context, playlistID uuid.UUID, isPrivate bool) error {
	const query = `
		UPDATE playlists
		SET is_private = $1
		WHERE id = $2
	`

	ct, err := r.pool.Exec(ctx, query, isPrivate, playlistID)
	if err != nil {
		return fmt.Errorf("update access meta: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("playlist %s not found", playlistID)
	}

	return nil
}

func (r *PlaylistAccessRepository) DeleteAccessMeta(ctx context.Context, playlistID uuid.UUID) error {
	return nil
}
