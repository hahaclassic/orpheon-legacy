package assign

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type ArtistAssignRepository interface {
	AssignArtistToTrack(ctx context.Context, artistID uuid.UUID, trackID uuid.UUID) error
	AssignArtistToAlbum(ctx context.Context, artistID uuid.UUID, albumID uuid.UUID) error
	UnassignArtistFromTrack(ctx context.Context, artistID uuid.UUID, trackID uuid.UUID) error
	UnassignArtistFromAlbum(ctx context.Context, artistID uuid.UUID, albumID uuid.UUID) error

	GetArtistAlbums(ctx context.Context, artistID uuid.UUID) ([]*entity.AlbumMeta, error)
	GetArtistTracks(ctx context.Context, artistID uuid.UUID) ([]*entity.TrackMeta, error)
	GetArtistByAlbum(ctx context.Context, albumID uuid.UUID) ([]*entity.ArtistMeta, error)
	GetArtistByTrack(ctx context.Context, trackID uuid.UUID) ([]*entity.ArtistMeta, error)
}

type ArtistAssignService struct {
	repo ArtistAssignRepository
}

func NewArtistAssignService(repo ArtistAssignRepository) *ArtistAssignService {
	return &ArtistAssignService{
		repo: repo,
	}
}

func (a *ArtistAssignService) AssignArtistToTrack(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAssignArtistOnTrack, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.AssignArtistToTrack(ctx, artistID, trackID)
}

func (a *ArtistAssignService) AssignArtistToAlbum(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAssignArtistOnAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.AssignArtistToAlbum(ctx, artistID, albumID)
}

func (a *ArtistAssignService) GetArtistAlbums(ctx context.Context, artistID uuid.UUID) (albums []*entity.AlbumMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetArtistAlbums, err)
	}()

	albums, err = a.repo.GetArtistAlbums(ctx, artistID)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (a *ArtistAssignService) GetArtistTracks(ctx context.Context, artistID uuid.UUID) (tracks []*entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetArtistTracks, err)
	}()

	tracks, err = a.repo.GetArtistTracks(ctx, artistID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (a *ArtistAssignService) GetArtistByAlbum(ctx context.Context, albumID uuid.UUID) (artists []*entity.ArtistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetArtistByAlbum, err)
	}()

	artists, err = a.repo.GetArtistByAlbum(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (a *ArtistAssignService) GetArtistByTrack(ctx context.Context, trackID uuid.UUID) (artists []*entity.ArtistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetArtistByTrack, err)
	}()

	artists, err = a.repo.GetArtistByTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (a *ArtistAssignService) UnassignArtistFromTrack(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUnassignArtistFromTrack, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.UnassignArtistFromTrack(ctx, artistID, trackID)
}

func (a *ArtistAssignService) UnassignArtistFromAlbum(ctx context.Context, claims *entity.Claims, artistID uuid.UUID, albumID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUnassignArtistFromAlbum, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	return a.repo.UnassignArtistFromAlbum(ctx, artistID, albumID)
}
