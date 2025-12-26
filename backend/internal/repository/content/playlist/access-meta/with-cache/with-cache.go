package access_meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGetAccessMeta    = errors.New("failed to get access meta")
	ErrUpdateAccessMeta = errors.New("failed to update access meta")
	ErrDeleteAccessMeta = errors.New("failed to delete access meta")

	ErrCacheMiss = errors.New("cache miss")
)

type AccessCache interface {
	Get(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error)
	Set(ctx context.Context, playlistID uuid.UUID, meta *entity.PlaylistAccessMeta) error
	Delete(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistAccessRepository interface {
	GetAccessMeta(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error)
	UpdatePrivacy(ctx context.Context, playlistID uuid.UUID, isPrivate bool) error
}

type PlaylistAccessRepoWithCache struct {
	l1Cache AccessCache
	l2Cache AccessCache
	repo    PlaylistAccessRepository
}

type OptionFunc func(*PlaylistAccessRepoWithCache)

func WithL1Cache(l1 AccessCache) OptionFunc {
	return func(p *PlaylistAccessRepoWithCache) {
		p.l1Cache = l1
	}
}

func WithL2Cache(l2 AccessCache) OptionFunc {
	return func(p *PlaylistAccessRepoWithCache) {
		p.l2Cache = l2
	}
}

func New(repo PlaylistAccessRepository, options ...OptionFunc) *PlaylistAccessRepoWithCache {
	repoWithCache := &PlaylistAccessRepoWithCache{
		repo: repo,
	}
	for _, configure := range options {
		configure(repoWithCache)
	}

	return repoWithCache
}

func (p *PlaylistAccessRepoWithCache) GetAccessMeta(ctx context.Context, playlistID uuid.UUID) (_ *entity.PlaylistAccessMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrGetAccessMeta, err)
	}()

	var meta *entity.PlaylistAccessMeta

	if p.l1Cache != nil {
		meta, err = p.l1Cache.Get(ctx, playlistID)
		if err == nil {
			return meta, nil
		}
		if !errors.Is(err, ErrCacheMiss) {
			return nil, err
		}
	}

	if p.l2Cache != nil {
		meta, err = p.l2Cache.Get(ctx, playlistID)
		if err == nil {
			if p.l1Cache != nil {
				_ = p.l1Cache.Set(ctx, playlistID, meta) // the error is not checked
			}

			return meta, nil
		}
		if !errors.Is(err, ErrCacheMiss) {
			return nil, err
		}
	}

	meta, err = p.repo.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	if p.l2Cache != nil {
		_ = p.l2Cache.Set(ctx, playlistID, meta) // the error is not checked
	}
	if p.l1Cache != nil {
		_ = p.l1Cache.Set(ctx, playlistID, meta) // the error is not checked
	}

	return meta, nil
}

func (p *PlaylistAccessRepoWithCache) UpdatePrivacy(ctx context.Context, playlistID uuid.UUID, isPrivate bool) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrUpdateAccessMeta, err)
	}()

	err = p.repo.UpdatePrivacy(ctx, playlistID, isPrivate)
	if err != nil {
		return err
	}

	meta, err := p.GetAccessMeta(ctx, playlistID)
	if err != nil {
		return err
	}

	meta.IsPrivate = isPrivate

	if p.l2Cache != nil {
		err = p.l2Cache.Set(ctx, playlistID, meta)
		if err != nil {
			return err
		}
	}

	if p.l1Cache != nil {
		return p.l1Cache.Set(ctx, playlistID, meta)
	}

	return nil
}

func (p *PlaylistAccessRepoWithCache) DeleteAccessMeta(ctx context.Context, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrDeleteAccessMeta, err)
	}()

	if p.l2Cache != nil {
		err = p.l2Cache.Delete(ctx, playlistID)
		if err != nil {
			return err
		}
	}

	if p.l1Cache != nil {
		return p.l1Cache.Delete(ctx, playlistID)
	}

	return nil
}
