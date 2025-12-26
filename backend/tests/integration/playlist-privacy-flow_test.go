package integration_test

import (
	"context"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	access_cache_local "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-cache/local"
	access_cache_redis "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-cache/redis"
	access_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/default/postgres"
	access_meta "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/with-cache"
	playlist_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/meta/postgres"
	user_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/user/postgres"
	"github.com/stretchr/testify/require"
)

// ==========================================
//        playlist privacy flow
//
// requirements:
//    1. access cache (local)
//    2. access cache (redis)
//    3. access repo
//    4. access repo with cache
//    5. user repo
//    6. playlist meta repo
//
// flow:
//    1. create user
//    2. create public playlist
//    3. get access meta (is_private = false)
//    4. update playlist privacy
//    5. get access meta (is_private = true)
// ==========================================

func TestPlaylistAccessFlow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic occurred: %v\n%s", r, debug.Stack())
		}
	}()

	require.NoError(t, runMigrationsUp(pgxPool))
	defer func() {
		require.NoError(t, runMigrationsDown(pgxPool))
	}()
	defer func() {
		require.NoError(t, clearRedis(redisClient))
	}()

	ctx := context.Background()
	localCache, err := access_cache_local.NewAccessCache(128)
	require.NoError(t, err)

	redisCache := access_cache_redis.NewAccessCache(redisClient,
		&access_cache_redis.TTLConfig{TTL: 10 * time.Minute, Jitter: 1 * time.Minute})

	accessRepo := access_meta_postgres.NewPlaylistAccessRepository(pgxPool)
	accessRepoWithCache := access_meta.New(accessRepo, access_meta.WithL1Cache(localCache), access_meta.WithL2Cache(redisCache))

	userRepo := user_postgres.NewUserRepository(pgxPool)
	playlistMetaRepo := playlist_meta_postgres.NewPlaylistMetaRepository(pgxPool)

	// 1.
	user1 := &entity.UserInfo{
		ID:               uuid.New(),
		Name:             "user1",
		AccessLvl:        entity.User,
		BirthDate:        time.Date(2000, 10, 10, 0, 0, 0, 0, time.Local),
		RegistrationDate: time.Now(),
	}
	require.NoError(t, userRepo.CreateUser(ctx, user1))

	// 2.
	playlist := &entity.PlaylistMeta{
		ID:        uuid.New(),
		OwnerID:   user1.ID,
		Name:      "playlist",
		IsPrivate: false,
	}
	require.NoError(t, playlistMetaRepo.Create(ctx, playlist))

	// 3.
	meta, err := accessRepoWithCache.GetAccessMeta(ctx, playlist.ID)
	require.NoError(t, err)
	require.False(t, meta.IsPrivate, "playlist should be public")

	// 4.
	require.NoError(t, accessRepoWithCache.UpdatePrivacy(ctx, playlist.ID, true))

	// 5.
	meta, err = accessRepoWithCache.GetAccessMeta(ctx, playlist.ID)
	require.NoError(t, err)
	require.True(t, meta.IsPrivate, "playlist should be private")
}
