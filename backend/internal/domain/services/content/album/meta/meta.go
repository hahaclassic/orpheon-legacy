package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGenerateID = errors.New("id generation error")
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	GetAlbum(ctx context.Context, id uuid.UUID) (*entity.AlbumMeta, error)
	UpdateAlbum(ctx context.Context, album *entity.AlbumMeta) error
	DeleteAlbum(ctx context.Context, id uuid.UUID) error
	GetAllAlbums(ctx context.Context) ([]*entity.AlbumMeta, error)
	GetAlbumArtists(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error)
	GetAlbumGenres(ctx context.Context, albumID uuid.UUID) ([]*entity.Genre, error)
}

type AlbumService struct {
	repo AlbumRepository
}

func New(repo AlbumRepository) *AlbumService {
	return &AlbumService{
		repo: repo,
	}
}

func (a *AlbumService) CreateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) (id uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateAlbum, err)
	}()
	if claims == nil || claims.AccessLvl != entity.Admin {
		return uuid.Nil, commonerr.ErrForbidden
	}

	album.ID, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, ErrGenerateID
	}

	err = a.repo.CreateAlbum(ctx, album)
	if err != nil {
		return uuid.Nil, err
	}

	return album.ID, nil
}

func (a *AlbumService) GetAlbum(ctx context.Context, albumID uuid.UUID) (_ *entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAlbum, err)
	}()

	return a.repo.GetAlbum(ctx, albumID)
}

func (a *AlbumService) GetAllAlbums(ctx context.Context) (_ []*entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllAlbums, err)
	}()

	return a.repo.GetAllAlbums(ctx)
}

func (a *AlbumService) UpdateAlbum(ctx context.Context, claims *entity.Claims, album *entity.AlbumMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.UpdateAlbum(ctx, album)
}

func (a *AlbumService) DeleteAlbum(ctx context.Context, claims *entity.Claims, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.DeleteAlbum(ctx, albumID)
}
