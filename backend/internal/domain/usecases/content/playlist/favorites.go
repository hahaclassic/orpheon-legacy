package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrAddToUserFavorites           = errors.New("failed to add playlist to user favorites")
	ErrGetUserFavorites             = errors.New("failed to get user favorites")
	ErrGetUsersWithFavoritePlaylist = errors.New("failed to get users with favorite playlist")
	ErrDeleteFromUserFavorites      = errors.New("failed to delete playlist from user favorites")
	ErrDeleteFromAllFavorites       = errors.New("failed to delete favorite playlist for all users")
	ErrAddPlaylistToAllFavorites    = errors.New("failed to add playlist to all favorites")
	ErrIsFavorite                   = errors.New("failed to check if playlist is favorite")
)

type PlaylistFavoriteService interface {
	GetUserFavorites(ctx context.Context, claims *entity.Claims) ([]*entity.PlaylistMeta, error)
	AddToUserFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	DeleteFromUserFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	IsFavorite(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
}
