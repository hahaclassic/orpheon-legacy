package genre_assign

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type GenreAssignRepository interface {
	AssignGenreToAlbum(ctx context.Context, genreID uuid.UUID, albumID uuid.UUID) error
	UnassignGenreFromAlbum(ctx context.Context, genreID uuid.UUID, albumID uuid.UUID) error
}

type GenreAssignService struct {
	repo GenreAssignRepository
}

func NewGenreAssignService(repo GenreAssignRepository) *GenreAssignService {
	return &GenreAssignService{repo: repo}
}

func (s *GenreAssignService) AssignGenreToAlbum(ctx context.Context, claims *entity.Claims, genreID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAssignGenreToAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.AssignGenreToAlbum(ctx, genreID, albumID)
}

func (s *GenreAssignService) UnassignGenreFromAlbum(ctx context.Context, claims *entity.Claims, genreID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUnassignGenreFromAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return s.repo.UnassignGenreFromAlbum(ctx, genreID, albumID)
}
