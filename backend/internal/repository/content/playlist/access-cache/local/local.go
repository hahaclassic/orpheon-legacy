package access_cache_local

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	access_meta "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/with-cache"
	lru "github.com/hashicorp/golang-lru/v2"
)

type AccessCache struct {
	cache *lru.Cache[uuid.UUID, *entity.PlaylistAccessMeta]
}

func NewAccessCache(length int) (*AccessCache, error) {
	cache, err := lru.New[uuid.UUID, *entity.PlaylistAccessMeta](length)
	if err != nil {
		return nil, fmt.Errorf("failed to create local cache: %w", err)
	}

	return &AccessCache{
		cache: cache,
	}, nil
}

func (a *AccessCache) Set(_ context.Context, playlistID uuid.UUID, meta *entity.PlaylistAccessMeta) error {
	a.cache.Add(playlistID, meta)

	return nil
}

func (a *AccessCache) Get(_ context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error) {
	meta, ok := a.cache.Get(playlistID)
	if !ok {
		return nil, access_meta.ErrCacheMiss
	}

	return meta, nil
}

func (a *AccessCache) Delete(_ context.Context, playlistID uuid.UUID) error {
	a.cache.Remove(playlistID)

	return nil
}
