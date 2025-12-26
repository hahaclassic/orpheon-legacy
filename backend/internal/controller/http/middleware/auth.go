package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	jwttokens "github.com/hahaclassic/orpheon/backend/internal/adapters/tokens/jwt"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
)

func AuthMiddlewareRequired(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Middleware triggered")

		token, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := authService.GetClaims(c, token)
		if err != nil {
			var h gin.H
			if errors.Is(err, jwttokens.ErrExpired) {
				h = gin.H{"error": "token expired"}
			} else {
				h = gin.H{"error": "invalid token"}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, h)
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func AuthMiddlewareOptional(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Middleware triggered")

		token, err := c.Cookie("access_token")
		if err != nil {
			c.Next()
		}

		claims, err := authService.GetClaims(c, token)
		if err != nil {
			c.Next()
		}

		c.Set("claims", claims)
		c.Next()
	}
}
