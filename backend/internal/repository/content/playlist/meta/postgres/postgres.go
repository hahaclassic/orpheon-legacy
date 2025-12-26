package playlist_meta_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistMetaRepository struct {
	pool *pgxpool.Pool
}

func NewPlaylistMetaRepository(pool *pgxpool.Pool) *PlaylistMetaRepository {
	return &PlaylistMetaRepository{pool: pool}
}

func (r *PlaylistMetaRepository) Create(ctx context.Context, playlist *entity.PlaylistMeta) error {
	const query = `
		INSERT INTO playlists (id, name, description, is_private, owner_id, created_at, updated_at, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query, playlist.ID, playlist.Name, playlist.Description,
		playlist.IsPrivate, playlist.OwnerID, playlist.CreatedAt, playlist.UpdatedAt, 0)
	if err != nil {
		return fmt.Errorf("create playlist: %w", err)
	}

	return nil
}

func (r *PlaylistMetaRepository) GetByID(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistMeta, error) {
	const query = `
		SELECT id, owner_id, name, description, is_private, rating, created_at, updated_at
		FROM playlists
		WHERE id = $1
	`

	var playlist entity.PlaylistMeta
	err := r.pool.QueryRow(ctx, query, playlistID).Scan(&playlist.ID, &playlist.OwnerID,
		&playlist.Name, &playlist.Description, &playlist.IsPrivate, &playlist.Rating,
		&playlist.CreatedAt, &playlist.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get playlist by id: %w", err)
	}

	return &playlist, nil
}

func (r *PlaylistMetaRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.PlaylistMeta, error) {
	const query = `
		SELECT id, owner_id, name, description, is_private, rating, created_at, updated_at
		FROM playlists
		WHERE owner_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get playlists by user: %w", err)
	}
	defer rows.Close()

	var playlists []*entity.PlaylistMeta
	for rows.Next() {
		var playlist entity.PlaylistMeta
		if err := rows.Scan(&playlist.ID, &playlist.OwnerID, &playlist.Name, &playlist.Description,
			&playlist.IsPrivate, &playlist.Rating, &playlist.CreatedAt, &playlist.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan playlist: %w", err)
		}
		playlists = append(playlists, &playlist)
	}

	return playlists, nil
}

func (r *PlaylistMetaRepository) Update(ctx context.Context, playlist *entity.PlaylistMeta) error {
	const query = `
		UPDATE playlists
		SET name = $1, description = $2, is_private = $3, updated_at = $4
		WHERE id = $5
	`

	ct, err := r.pool.Exec(ctx, query, playlist.Name, playlist.Description, playlist.IsPrivate, playlist.UpdatedAt,
		playlist.ID)
	if err != nil {
		return fmt.Errorf("update playlist: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("playlist %s not found", playlist.ID)
	}

	return nil
}

func (r *PlaylistMetaRepository) Delete(ctx context.Context, playlistID uuid.UUID) error {
	const query = `
		DELETE FROM playlists
		WHERE id = $1
	`

	ct, err := r.pool.Exec(ctx, query, playlistID)
	if err != nil {
		return fmt.Errorf("delete playlist: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("playlist %s not found", playlistID)
	}

	return nil
}
