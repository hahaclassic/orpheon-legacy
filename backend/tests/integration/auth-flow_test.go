package integration_test

import (
	"context"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	auth_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/auth/auth-repo/postgres"
	refresh_redis "github.com/hahaclassic/orpheon/backend/internal/repository/auth/refresh-token/redis"
	user_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/user/postgres"
	"github.com/stretchr/testify/require"
)

// ==========================================
//                auth flow
//
// flow:
//    1. create user
//    2. save user credentials
//    3. save refresh token + claims
//    4. updated password
//        4.1 check new password
//    5. get refresh token [fail, short ttl]
// ==========================================

func TestAuthFlow(t *testing.T) {
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

	userCreator := user_postgres.NewUserRepository(pgxPool)
	authRepo := auth_postgres.NewAuthRepository(pgxPool)
	refreshTokenRepo := refresh_redis.NewRefreshTokenRepository(redisClient,
		&refresh_redis.TTLConfig{TTL: 1 * time.Second, Jitter: 1 * time.Nanosecond}) // short ttl for test

	ctx := context.Background()

	// 1.
	userInfo := &entity.UserInfo{
		ID:               uuid.New(),
		Name:             "test_user",
		AccessLvl:        entity.User,
		BirthDate:        time.Date(2000, 10, 10, 0, 0, 0, 0, time.Local),
		RegistrationDate: time.Now(),
	}

	require.NoError(t, userCreator.CreateUser(ctx, userInfo))

	// 2.
	creds := &entity.UserCredentials{Login: "test_user", Password: "secret123"}
	require.NoError(t, authRepo.SaveCredentials(ctx, userInfo.ID, creds))

	// 3.
	refreshToken := "refresh-token-123"
	claims := &entity.Claims{UserID: userInfo.ID, AccessLvl: entity.User}
	require.NoError(t, refreshTokenRepo.Set(ctx, refreshToken, claims))

	// 4.
	newPassword := "newsecret123"
	require.NoError(t, authRepo.UpdatePassword(ctx, userInfo.ID, newPassword))

	// 4.1
	savedPwd, err := authRepo.GetPasswordByID(ctx, userInfo.ID)
	require.NoError(t, err)
	require.Equal(t, newPassword, savedPwd)

	time.Sleep(2 * time.Second)
	// 5.
	_, err = refreshTokenRepo.Get(ctx, refreshToken)
	require.Error(t, err) // needs errorsIs(ErrCacheMiss)
}
