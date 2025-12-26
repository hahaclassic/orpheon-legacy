package favorites_postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistFavoriteRepository struct {
	pool *pgxpool.Pool
}

func NewPlaylistFavoriteRepository(pool *pgxpool.Pool) *PlaylistFavoriteRepository {
	return &PlaylistFavoriteRepository{pool: pool}
}

func (r *PlaylistFavoriteRepository) AddToFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error {
	const query = `
		INSERT INTO favorite_playlists (user_id, playlist_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, userID, playlistID)
	if err != nil {
		return fmt.Errorf("add to favorites: %w", err)
	}
	return nil
}

func (r *PlaylistFavoriteRepository) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*entity.PlaylistMeta, error) {
	const query = `
		SELECT p.id, p.owner_id, p.name, p.description, p.is_private, p.created_at, p.updated_at, p.rating
		FROM favorite_playlists f
		JOIN playlists p ON p.id = f.playlist_id
		WHERE f.user_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user favorites: %w", err)
	}
	defer rows.Close()

	var result []*entity.PlaylistMeta
	for rows.Next() {
		var meta entity.PlaylistMeta
		if err := rows.Scan(&meta.ID, &meta.OwnerID, &meta.Name, &meta.Description,
			&meta.IsPrivate, &meta.CreatedAt, &meta.UpdatedAt, &meta.Rating); err != nil {
			return nil, fmt.Errorf("scan playlist meta: %w", err)
		}
		result = append(result, &meta)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return result, nil
}

func (r *PlaylistFavoriteRepository) DeleteFromUserFavorites(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) error {
	const query = `
		DELETE FROM favorite_playlists
		WHERE user_id = $1 AND playlist_id = $2
	`

	_, err := r.pool.Exec(ctx, query, userID, playlistID)
	if err != nil {
		return fmt.Errorf("delete from user favorites: %w", err)
	}
	return nil
}

func (r *PlaylistFavoriteRepository) GetUsersWithFavoritePlaylist(ctx context.Context, playlistID uuid.UUID, withOwner bool) ([]uuid.UUID, error) {
	var query string
	if withOwner {
		query = `
			SELECT user_id
			FROM favorite_playlists
			WHERE playlist_id = $1
	    `
	} else {
		query = `
			SELECT user_id
			FROM favorite_playlists
			WHERE playlist_id = $1 AND user_id != (select owner_id from playlists where id = $1)
		`
	}
	rows, err := r.pool.Query(ctx, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("get users with favorite: %w", err)
	}
	defer rows.Close()

	var users []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan user id: %w", err)
		}
		users = append(users, id)
	}

	return users, nil
}

func (r *PlaylistFavoriteRepository) DeleteFromAllFavorites(ctx context.Context, playlistID uuid.UUID, withOwner bool) error {
	var (
		query string
		err   error
	)
	if withOwner {
		query = `
			DELETE FROM favorite_playlists
			WHERE playlist_id = $1
	    `
	} else {
		query = `
			DELETE FROM favorite_playlists
			WHERE playlist_id = $1 AND user_id != (select owner_id from playlists where id = $1) 
		`
	}

	if _, err = r.pool.Exec(ctx, query, playlistID); err != nil {
		return fmt.Errorf("delete from all favorites: %w", err)
	}

	return nil
}

func (r *PlaylistFavoriteRepository) RestoreAllFavorites(ctx context.Context, userIDs []uuid.UUID, playlistID uuid.UUID) error {
	const query = `
		INSERT INTO favorite_playlists (user_id, playlist_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	var batch pgx.Batch
	for _, userID := range userIDs {
		batch.Queue(query, userID, playlistID)
	}

	br := r.pool.SendBatch(ctx, &batch)
	defer func() {
		err := br.Close()
		if err != nil {
			slog.Error("err", "batch results close error", err)
		}
	}()

	for range userIDs {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("restore favorite: %w", err)
		}
	}

	return nil
}

func (r *PlaylistFavoriteRepository) IsFavorite(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM favorite_playlists
			WHERE user_id = $1 AND playlist_id = $2
		)
	`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, userID, playlistID).Scan(&exists); err != nil {
		return false, fmt.Errorf("is favorite: %w", err)
	}

	return exists, nil
}
