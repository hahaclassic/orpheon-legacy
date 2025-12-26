package genre

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *entity.Genre) error
	GetByID(ctx context.Context, genreID uuid.UUID) (*entity.Genre, error)
	GetAll(ctx context.Context) ([]*entity.Genre, error)
	Update(ctx context.Context, genre *entity.Genre) error
	Delete(ctx context.Context, genreID uuid.UUID) error
	GetByAlbum(ctx context.Context, albumID uuid.UUID) ([]*entity.Genre, error)
}

var (
	ErrInvalidGenreID = errors.New("invalid genre ID")
	ErrGenerateID     = errors.New("generate ID error")
)

type GenreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) CreateGenre(ctx context.Context, claims *entity.Claims, genre *entity.Genre) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateGenre, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	genre.ID, err = uuid.NewRandom()
	if err != nil {
		return ErrGenerateID
	}

	return s.repo.Create(ctx, genre)
}

func (s *GenreService) GetGenreByID(ctx context.Context, genreID uuid.UUID) (_ *entity.Genre, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetGenre, err)
	}()

	if genreID == uuid.Nil {
		return nil, ErrInvalidGenreID
	}

	return s.repo.GetByID(ctx, genreID)
}

func (s *GenreService) GetAllGenres(ctx context.Context) (_ []*entity.Genre, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllGenres, err)
	}()

	return s.repo.GetAll(ctx)
}

func (s *GenreService) UpdateGenre(ctx context.Context, claims *entity.Claims, genre *entity.Genre) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateGenre, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	if genre.ID == uuid.Nil {
		return ErrInvalidGenreID
	}

	return s.repo.Update(ctx, genre)
}

func (s *GenreService) DeleteGenre(ctx context.Context, claims *entity.Claims, genreID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteGenre, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	if genreID == uuid.Nil {
		return ErrInvalidGenreID
	}

	return s.repo.Delete(ctx, genreID)
}

func (s *GenreService) GetGenreByAlbum(ctx context.Context, albumID uuid.UUID) (_ []*entity.Genre, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetGenreByAlbum, err)
	}()

	return s.repo.GetByAlbum(ctx, albumID)
}
