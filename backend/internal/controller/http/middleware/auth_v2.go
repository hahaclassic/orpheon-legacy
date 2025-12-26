package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	jwttokens "github.com/hahaclassic/orpheon/backend/internal/adapters/tokens/jwt"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/cookie"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
)

var (
	ErrNoAccessToken = errors.New("no access token")
	ErrNoTokens      = errors.New("no tokens")
	ErrExpiredToken  = errors.New("token expired")
	ErrInvalidToken  = errors.New("invalid token")
)

type AuthMiddleware struct {
	authService        auth.AuthService
	cookieTokensSetter *cookie.CookieTokensSetter
}

func NewAuthMiddleware(authService auth.AuthService, cookieTokensSetter *cookie.CookieTokensSetter) *AuthMiddleware {
	return &AuthMiddleware{
		authService:        authService,
		cookieTokensSetter: cookieTokensSetter,
	}
}

func (a *AuthMiddleware) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := a.setClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

			return
		}

		c.Next()
	}
}

func (a *AuthMiddleware) Optional() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = a.setClaims(c)
		c.Next()
	}
}

func (a *AuthMiddleware) setClaims(c *gin.Context) error {
	accessToken, err := c.Cookie(cookie.AccessCookieName)
	if err != nil || accessToken == "" {
		refreshToken, err := c.Cookie(cookie.RefreshCookieName)
		if err != nil || refreshToken == "" {
			return ErrNoTokens
		}

		tokens, err := a.authService.RefreshTokens(c.Request.Context(), refreshToken)
		if err != nil {
			return err
		}

		a.cookieTokensSetter.SetAll(c, tokens)
		accessToken = tokens.Access
	}

	claims, err := a.authService.GetClaims(c, accessToken)
	if errors.Is(err, jwttokens.ErrExpired) {
		return ErrExpiredToken
	} else if err != nil {
		return ErrInvalidToken
	}

	c.Set("claims", claims)

	return nil
}
