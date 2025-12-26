package refresh_redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

// type TTLConfig struct {
// 	TTL    time.Duration `env:"REDIS_TTL" env-required:"true"`
// 	Jitter time.Duration `env:"RESID_JITTER" env-required:"true"`
// }

type TTLConfig = config.RefreshTokenConfig

type RefreshTokenRepository struct {
	client *redis.Client
	conf   *TTLConfig
}

func NewRefreshTokenRepository(client *redis.Client, conf *TTLConfig) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		client: client,
		conf:   conf,
	}
}

func (r *RefreshTokenRepository) key(token string) string {
	return "refresh_token:" + token
}

func (r *RefreshTokenRepository) Set(ctx context.Context, token string, claims *entity.Claims) error {
	data, err := json.Marshal(claims)
	if err != nil {
		return fmt.Errorf("marshal claims: %w", err)
	}

	if err := r.client.Set(ctx, r.key(token), data, r.getTTL()).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) Get(ctx context.Context, token string) (*entity.Claims, error) {
	data, err := r.client.Get(ctx, r.key(token)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("refresh token not found")
		}
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var claims entity.Claims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	return &claims, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	if err := r.client.Del(ctx, r.key(token)).Err(); err != nil {
		return fmt.Errorf("redis delete: %w", err)
	}
	return nil
}

func (a *RefreshTokenRepository) getTTL() time.Duration {
	jitter := time.Duration(rand.Int63n(a.conf.Jitter.Nanoseconds()))
	return a.conf.TTL + jitter
}
