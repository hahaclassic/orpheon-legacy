package privacy

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type PlaylistPrivacyRepository interface {
	UpdatePrivacy(ctx context.Context, playlistID uuid.UUID, isPrivate bool) error
}

type FavoritesDeletionService interface {
	GetUsersWithFavoritePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, withOwner bool) ([]uuid.UUID, error)
	DeleteFromAllFavorites(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, withOwner bool) error
	AddPlaylistToAllFavorites(ctx context.Context, claims *entity.Claims, userIDs []uuid.UUID, playlistID uuid.UUID) error
}

type PlaylistPrivacyChanger struct {
	playlistPolicyService playlist.PlaylistPolicyService
	favoritesService      FavoritesDeletionService
	accessRepo            PlaylistPrivacyRepository
}

func NewPlaylistPrivacyChanger(playlistPolicyService playlist.PlaylistPolicyService,
	favoritesService FavoritesDeletionService, accessRepo PlaylistPrivacyRepository) *PlaylistPrivacyChanger {
	return &PlaylistPrivacyChanger{
		playlistPolicyService: playlistPolicyService,
		favoritesService:      favoritesService,
		accessRepo:            accessRepo,
	}
}

type rollback func() error

func (p *PlaylistPrivacyChanger) ChangePrivacy(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, isPrivate bool) (err error) {
	var rollbacks []rollback

	if p.playlistPolicyService.CanEdit(ctx, claims, playlistID) != nil {
		return commonerr.ErrForbidden
	}

	defer func() {
		if err != nil {
			err = errwrap.Wrap(playlist.ErrChangePrivacy, err)

			for i := len(rollbacks) - 1; i >= 0; i-- {
				_ = rollbacks[i]()
			}
		}
	}()

	if isPrivate {
		userIDs, err := p.favoritesService.GetUsersWithFavoritePlaylist(ctx, claims, playlistID, false)
		if err != nil {
			return err
		}

		err = p.favoritesService.DeleteFromAllFavorites(ctx, claims, playlistID, false)
		if err != nil {
			return err
		}

		rollbacks = append(rollbacks, rollback(func() error {
			return p.favoritesService.AddPlaylistToAllFavorites(ctx, claims, userIDs, playlistID)
		}))
	}

	err = p.accessRepo.UpdatePrivacy(ctx, playlistID, isPrivate)
	if err != nil {
		return err
	}

	return nil
}
