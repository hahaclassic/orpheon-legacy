package genre

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetGenre        = errors.New("failed to get genre")
	ErrGetAllGenres    = errors.New("failed to get all genres")
	ErrCreateGenre     = errors.New("failed to create genre")
	ErrUpdateGenre     = errors.New("failed to update genre")
	ErrDeleteGenre     = errors.New("failed to delete genre")
	ErrGetGenreByAlbum = errors.New("failed to get genre by album")
)

type GenreService interface {
	GetGenreByID(ctx context.Context, genreID uuid.UUID) (_ *entity.Genre, err error)
	GetAllGenres(ctx context.Context) (_ []*entity.Genre, err error)

	CreateGenre(ctx context.Context, claims *entity.Claims, genre *entity.Genre) (err error)
	UpdateGenre(ctx context.Context, claims *entity.Claims, genre *entity.Genre) (err error)
	DeleteGenre(ctx context.Context, claims *entity.Claims, genreID uuid.UUID) (err error)

	GetGenreByAlbum(ctx context.Context, albumID uuid.UUID) (_ []*entity.Genre, err error)
}
