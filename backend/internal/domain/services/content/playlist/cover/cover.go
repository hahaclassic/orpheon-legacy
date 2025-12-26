package cover

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type PlaylistCoverRepository interface {
	SaveCover(ctx context.Context, cover *entity.Cover) error
	GetCover(ctx context.Context, playlistID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistCoverService struct {
	policy usecase.PlaylistPolicyService
	repo   PlaylistCoverRepository
}

func New(repo PlaylistCoverRepository, policy usecase.PlaylistPolicyService) *PlaylistCoverService {
	return &PlaylistCoverService{
		policy: policy,
		repo:   repo,
	}
}

func (c *PlaylistCoverService) GetCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (_ *entity.Cover, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetCover, err)
	}()

	err = c.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	return c.repo.GetCover(ctx, playlistID)
}

func (c *PlaylistCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUploadCover, err)
	}()

	err = c.policy.CanEdit(ctx, claims, cover.ObjectID)
	if err != nil {
		return err
	}

	return c.repo.SaveCover(ctx, cover)
}

func (c *PlaylistCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteCover, err)
	}()

	err = c.policy.CanDelete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return c.repo.DeleteCover(ctx, playlistID)
}
