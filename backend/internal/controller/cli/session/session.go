package session

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/ilyakaznacheev/cleanenv"
)

var (
	once        sync.Once
	s           *Session
	refreshFile = ".refresh.env"
)

type RefreshToken struct {
	Token string `env:"ORPHEON_REFRESH_TOKEN"`
}

type Session struct {
	claims *entity.Claims
	tokens *entity.AuthTokens
	mu     *sync.RWMutex
}

func Instance() *Session {
	once.Do(func() {
		refresh := &RefreshToken{}
		_ = cleanenv.ReadConfig(refreshFile, refresh)

		s = &Session{
			mu: &sync.RWMutex{},
			claims: &entity.Claims{
				UserID:    uuid.Nil,
				AccessLvl: entity.Unauthorized,
			},
			tokens: &entity.AuthTokens{
				Refresh: refresh.Token,
			},
		}
	})
	return s
}

func StartSession(claims *entity.Claims, tokens *entity.AuthTokens) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tokens.Refresh != "" {
		envContent := fmt.Sprintf("ORPHEON_REFRESH_TOKEN=%s", tokens.Refresh)
		err := os.WriteFile(refreshFile, []byte(envContent), 0644)
		if err != nil {
			fmt.Printf("Failed to save refresh token to .env: %v\n", err)
		}
	}

	s.claims = claims
	s.tokens = tokens
}

func EndSession() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.claims.UserID = uuid.Nil
	s.claims.AccessLvl = entity.Unauthorized
	s.tokens.Access = ""
	s.tokens.Refresh = ""
}

func Claims() *entity.Claims {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.claims
}

func Tokens() *entity.AuthTokens {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.tokens
}

func UpdateTokens(tokens *entity.AuthTokens) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens = tokens
}

func IsAuthenticated() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.claims.UserID != uuid.Nil
}
