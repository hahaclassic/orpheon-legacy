package meta

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrEmptyPlaylistName = errors.New("playlist name cannot be empty")
)

type PlaylistMetaRepository interface {
	Create(ctx context.Context, playlist *entity.PlaylistMeta) error
	GetByID(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistMeta, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.PlaylistMeta, error)
	Update(ctx context.Context, playlist *entity.PlaylistMeta) error
	Delete(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistAccessMetaDeleter interface {
	DeleteAccessMeta(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistMetaService struct {
	repo       PlaylistMetaRepository
	policy     usecase.PlaylistPolicyService
	accessRepo PlaylistAccessMetaDeleter
}

func NewPlaylistMetaService(repo PlaylistMetaRepository, policy usecase.PlaylistPolicyService, accessRepo PlaylistAccessMetaDeleter) *PlaylistMetaService {
	return &PlaylistMetaService{
		repo:       repo,
		policy:     policy,
		accessRepo: accessRepo,
	}
}

func (p *PlaylistMetaService) CreateMeta(ctx context.Context, claims *entity.Claims, playlist *entity.PlaylistMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateMeta, err)
	}()

	if playlist.Name == "" {
		return ErrEmptyPlaylistName
	}

	playlist.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	playlist.OwnerID = claims.UserID

	playlist.CreatedAt = time.Now()
	playlist.UpdatedAt = playlist.CreatedAt

	return p.repo.Create(ctx, playlist)
}

func (p *PlaylistMetaService) GetMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (_ *entity.PlaylistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetMeta, err)
	}()

	if err = p.policy.CanView(ctx, claims, playlistID); err != nil {
		return nil, err
	}

	return p.repo.GetByID(ctx, playlistID)
}

func (p *PlaylistMetaService) GetUserAllPlaylistsMeta(ctx context.Context, claims *entity.Claims, userID uuid.UUID) (_ []*entity.PlaylistMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetUserPlaylistsMeta, err)
	}()

	playlists, err := p.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// if user is owner, show all playlists (user can see all his playlists)
	if claims != nil && claims.UserID == userID {
		return playlists, nil
	}

	publicPlaylists := playlists[:]
	currIdx := 0
	for i := range playlists {
		if !playlists[i].IsPrivate {
			publicPlaylists[currIdx] = playlists[i]
			currIdx++
		}
	}

	return publicPlaylists[:currIdx], nil
}

func (p *PlaylistMetaService) UpdateMeta(ctx context.Context, claims *entity.Claims, playlist *entity.PlaylistMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateMeta, err)
	}()

	if err = p.policy.CanEdit(ctx, claims, playlist.ID); err != nil {
		return err
	}

	playlist.UpdatedAt = time.Now()

	return p.repo.Update(ctx, playlist)
}

func (p *PlaylistMetaService) DeleteMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteMeta, err)
	}()

	if err = p.policy.CanDelete(ctx, claims, playlistID); err != nil {
		return err
	}

	if err = p.accessRepo.DeleteAccessMeta(ctx, playlistID); err != nil {
		return err
	}

	return p.repo.Delete(ctx, playlistID)
}
