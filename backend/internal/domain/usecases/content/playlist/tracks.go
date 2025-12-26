package playlist

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrAddTrack            = errors.New("track addition error")
	ErrGetAllTracks        = errors.New("get all tracks error")
	ErrDeleteTrack         = errors.New("delete track error")
	ErrDeleteAllTracks     = errors.New("delete all tracks error")
	ErrChangeTrackPosition = errors.New("change track position error")
	ErrRestoreAllTracks    = errors.New("restore all tracks error")
)

type PlaylistTrackService interface {
	AddTrack(ctx context.Context, claims *entity.Claims, playlistTrack *entity.PlaylistTrack) error
	GetAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (tracks []*entity.TrackMeta, err error)
	DeleteTrack(ctx context.Context, claims *entity.Claims, playlistTrack *entity.PlaylistTrack) error
	DeleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	ChangeTrackPosition(ctx context.Context, claims *entity.Claims, playlistTrack *entity.PlaylistTrack) error
}
