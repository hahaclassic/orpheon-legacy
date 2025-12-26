package avatar

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type ArtistAvatarRepository interface {
	SaveCover(ctx context.Context, cover *entity.Cover) error
	GetCover(ctx context.Context, artistID uuid.UUID) (*entity.Cover, error)
	DeleteCover(ctx context.Context, artistID uuid.UUID) error
}

type ArtistCoverService struct {
	repo ArtistAvatarRepository
}

func NewArtistCoverService(repo ArtistAvatarRepository) *ArtistCoverService {
	return &ArtistCoverService{
		repo: repo,
	}
}

func (s *ArtistCoverService) GetCover(ctx context.Context, artistID uuid.UUID) (_ *entity.Cover, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAvatar, err)
	}()

	return s.repo.GetCover(ctx, artistID)
}

func (s *ArtistCoverService) UploadCover(ctx context.Context, claims *entity.Claims, cover *entity.Cover) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUploadAvatar, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.SaveCover(ctx, cover)
}

func (s *ArtistCoverService) DeleteCover(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAvatar, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.DeleteCover(ctx, artistID)
}
