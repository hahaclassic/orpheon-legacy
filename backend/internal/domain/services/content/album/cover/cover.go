package cover

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type AlbumCoverRepository interface {
	GetCover(ctx context.Context, albumID uuid.UUID) (*entity.Cover, error)
	SaveCover(ctx context.Context, cover *entity.Cover) error
	DeleteCover(ctx context.Context, albumID uuid.UUID) error
}

type AlbumCoverService struct {
	repo AlbumCoverRepository
}

func New(repo AlbumCoverRepository) *AlbumCoverService {
	return &AlbumCoverService{
		repo: repo,
	}
}

func (c *AlbumCoverService) GetCover(ctx context.Context, albumID uuid.UUID) (_ *entity.Cover, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetCover, err)
	}()

	return c.repo.GetCover(ctx, albumID)
}

func (c *AlbumCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUploadCover, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return c.repo.SaveCover(ctx, cover)
}

func (c *AlbumCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteCover, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return c.repo.DeleteCover(ctx, albumID)
}
