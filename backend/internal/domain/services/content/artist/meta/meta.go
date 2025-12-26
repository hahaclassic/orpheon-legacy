package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGenerateID = errors.New("id generation error")
)

type ArtistMetaRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ArtistMeta, error)
	GetAll(ctx context.Context) ([]*entity.ArtistMeta, error)
	Create(ctx context.Context, artist *entity.ArtistMeta) error
	Update(ctx context.Context, artist *entity.ArtistMeta) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ArtistMetaService struct {
	repo ArtistMetaRepository
}

func New(repo ArtistMetaRepository) *ArtistMetaService {
	return &ArtistMetaService{repo: repo}
}

func (s *ArtistMetaService) GetArtistMeta(ctx context.Context, artistID uuid.UUID) (_ *entity.ArtistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetArtistMeta, err)
	}()

	return s.repo.GetByID(ctx, artistID)
}

func (s *ArtistMetaService) GetAllArtistMeta(ctx context.Context) (_ []*entity.ArtistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllArtistMeta, err)
	}()

	return s.repo.GetAll(ctx)
}

func (s *ArtistMetaService) CreateArtistMeta(ctx context.Context, claims *entity.Claims, artist *entity.ArtistMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateArtistMeta, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	if artist.ID, err = uuid.NewRandom(); err != nil {
		return ErrGenerateID
	}

	return s.repo.Create(ctx, artist)
}

func (s *ArtistMetaService) UpdateArtistMeta(ctx context.Context, claims *entity.Claims, artist *entity.ArtistMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateArtistMeta, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.Update(ctx, artist)
}

func (s *ArtistMetaService) DeleteArtistMeta(ctx context.Context, claims *entity.Claims, artistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteArtistMeta, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.Delete(ctx, artistID)
}
