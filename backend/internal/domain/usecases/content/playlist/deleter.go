package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrDeletePlaylist = errors.New("failed to delete playlist")
)

type PlaylistDeletionService interface {
	DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
