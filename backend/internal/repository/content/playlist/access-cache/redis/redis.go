package access_cache_redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	access_meta "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/with-cache"
	"github.com/redis/go-redis/v9"
)

type TTLConfig = config.RedisAccessMetaConfig

// type TTLConfig struct {
// 	TTL    time.Duration `env:"REDIS_TTL" env-required:"true"`
// 	Jitter time.Duration `env:"RESID_JITTER" env-required:"true"`
// }

type AccessCache struct {
	client *redis.Client
	conf   *TTLConfig
}

func NewAccessCache(client *redis.Client, conf *TTLConfig) *AccessCache {
	return &AccessCache{
		client: client,
		conf:   conf,
	}
}

func (a *AccessCache) key(id uuid.UUID) string {
	return "playlist_access:" + id.String()
}

func (a *AccessCache) Set(ctx context.Context, playlistID uuid.UUID, meta *entity.PlaylistAccessMeta) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal playlist access meta: %w", err)
	}

	if err := a.client.Set(ctx, a.key(playlistID), data, a.getTTL()).Err(); err != nil {
		return fmt.Errorf("failed to set playlist access meta in redis: %w", err)
	}

	return nil
}

func (a *AccessCache) Get(ctx context.Context, playlistID uuid.UUID) (*entity.PlaylistAccessMeta, error) {
	data, err := a.client.Get(ctx, a.key(playlistID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, access_meta.ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get playlist access meta from redis: %w", err)
	}

	var meta entity.PlaylistAccessMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist access meta: %w", err)
	}

	return &meta, nil
}

func (a *AccessCache) Delete(ctx context.Context, playlistID uuid.UUID) error {
	if err := a.client.Del(ctx, a.key(playlistID)).Err(); err != nil {
		return fmt.Errorf("failed to delete playlist access meta from redis: %w", err)
	}
	return nil
}

func (a *AccessCache) getTTL() time.Duration {
	jitter := time.Duration(rand.Int63n(a.conf.Jitter.Nanoseconds()))
	return a.conf.TTL + jitter
}
