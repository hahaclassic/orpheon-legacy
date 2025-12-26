package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrCreateMeta           = errors.New("failed to create playlist meta")
	ErrGetMeta              = errors.New("failed to get playlist meta")
	ErrGetUserPlaylistsMeta = errors.New("failed to get user playlists meta")
	ErrUpdateMeta           = errors.New("failed to update playlist meta")
	ErrDeleteMeta           = errors.New("failed to delete playlist meta")
)

type PlaylistMetaService interface {
	CreateMeta(ctx context.Context, claims *entity.Claims, playlist *entity.PlaylistMeta) error
	GetMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (*entity.PlaylistMeta, error)
	GetUserAllPlaylistsMeta(ctx context.Context, claims *entity.Claims, userID uuid.UUID) ([]*entity.PlaylistMeta, error)
	UpdateMeta(ctx context.Context, claims *entity.Claims, playlist *entity.PlaylistMeta) error
	DeleteMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}
