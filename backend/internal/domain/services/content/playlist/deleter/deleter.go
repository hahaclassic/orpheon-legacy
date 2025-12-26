package deleter

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type MetaDeletionService interface {
	DeleteMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}

type TrackDeletionService interface {
	GetAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) ([]*entity.TrackMeta, error)
	DeleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
	RestoreAllTracks(ctx context.Context, claims *entity.Claims,
		playlistID uuid.UUID, trackIDs []uuid.UUID) error
}

type FavoritesDeletionService interface {
	GetUsersWithFavoritePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, withOwner bool) ([]uuid.UUID, error)
	DeleteFromAllFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, withOwner bool) error
	AddPlaylistToAllFavorites(ctx context.Context, claims *entity.Claims, userIDs []uuid.UUID, playlistID uuid.UUID) error
}

type PlaylistCoverDeletionService interface {
	GetCover(ctx context.Context, claims *entity.Claims, objectID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, claims *entity.Claims, objectID uuid.UUID) error
	UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) error
}

type PlaylistDeleter struct {
	meta      MetaDeletionService
	tracks    TrackDeletionService
	favorites FavoritesDeletionService
	cover     PlaylistCoverDeletionService
}

type OptionFunc func(*PlaylistDeleter)

func WithMetaDeletion(metaService MetaDeletionService) OptionFunc {
	return func(pd *PlaylistDeleter) {
		pd.meta = metaService
	}
}

func WithTracksDeletion(trackService TrackDeletionService) OptionFunc {
	return func(pd *PlaylistDeleter) {
		pd.tracks = trackService
	}
}

func WithFavoritesDeletion(favoriteService FavoritesDeletionService) OptionFunc {
	return func(pd *PlaylistDeleter) {
		pd.favorites = favoriteService
	}
}

func WithCoverDeletion(coverService PlaylistCoverDeletionService) OptionFunc {
	return func(pd *PlaylistDeleter) {
		pd.cover = coverService
	}
}

func New(options ...OptionFunc) *PlaylistDeleter {
	deleter := &PlaylistDeleter{}
	for _, configure := range options {
		configure(deleter)
	}

	return deleter
}

type rollback func() error

func (p *PlaylistDeleter) DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	var rollbacks []rollback

	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeletePlaylist, err)

			for i := len(rollbacks) - 1; i >= 0; i-- {
				_ = rollbacks[i]()
			}
		}
	}()

	if p.favorites != nil {
		rollback, err := p.deleteFavorites(ctx, claims, playlistID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, rollback)
	}

	if p.cover != nil {
		rollback, err := p.deleteCover(ctx, claims, playlistID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, rollback)
	}

	if p.tracks != nil {
		rollback, err := p.deleteAllTracks(ctx, claims, playlistID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, rollback)
	}

	if p.meta != nil {
		err = p.deleteMeta(ctx, claims, playlistID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PlaylistDeleter) deleteFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (rollback, error) {
	userIDs, err := p.favorites.GetUsersWithFavoritePlaylist(ctx, claims, playlistID, true)
	if err != nil {
		return nil, err
	}

	err = p.favorites.DeleteFromAllFavorites(ctx, claims, playlistID, true)
	if err != nil {
		return nil, err
	}

	return func() error {
		return p.favorites.AddPlaylistToAllFavorites(ctx, claims, userIDs, playlistID)
	}, nil
}

func (p *PlaylistDeleter) deleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (rollback, error) {
	tracks, err := p.tracks.GetAllTracks(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	err = p.tracks.DeleteAllTracks(ctx, claims, playlistID)
	if errors.Is(err, commonerr.ErrNotFound) {
		return func() error {
			return nil
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return func() error {
		trackIDs := make([]uuid.UUID, len(tracks))
		for i := range tracks {
			trackIDs[i] = tracks[i].ID
		}

		return p.tracks.RestoreAllTracks(ctx, claims, playlistID, trackIDs)
	}, nil
}

func (p *PlaylistDeleter) deleteCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (rollback, error) {
	cover, err := p.cover.GetCover(ctx, claims, playlistID)
	if errors.Is(err, commonerr.ErrNotFound) {
		return func() error {
			return nil
		}, nil
	}
	if err != nil {
		return nil, err
	}

	err = p.cover.DeleteCover(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	return func() error {
		return p.cover.UploadCover(ctx, claims, cover)
	}, nil
}

func (p *PlaylistDeleter) deleteMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	return p.meta.DeleteMeta(ctx, claims, playlistID)
}
