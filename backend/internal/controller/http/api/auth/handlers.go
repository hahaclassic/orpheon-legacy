package auth_ctrl

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/cookie"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
)

type AuthController struct {
	service            auth.AuthService
	cookieTokensSetter *cookie.CookieTokensSetter
	authMiddleware     gin.HandlerFunc
}

func NewAuthController(service auth.AuthService,
	cookieTokensSetter *cookie.CookieTokensSetter,
	authMiddleware gin.HandlerFunc) *AuthController {

	return &AuthController{
		service:            service,
		cookieTokensSetter: cookieTokensSetter,
		authMiddleware:     authMiddleware,
	}
}

func (ac *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.POST("/register", ac.register)
	authGroup.POST("/login", ac.login)
	authGroup.POST("/refresh", ac.refresh)
	authGroup.POST("/logout", ac.logout)

	passwordGroup := authGroup.Group("/password").Use(ac.authMiddleware)
	passwordGroup.POST("/update", ac.updatePassword)
}

func (ac *AuthController) register(c *gin.Context) {
	var creds entity.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	tokens, err := ac.service.RegisterUser(c.Request.Context(), &creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ac.cookieTokensSetter.SetAll(c, tokens)
}

func (ac *AuthController) login(c *gin.Context) {
	var creds entity.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	tokens, err := ac.service.Login(c.Request.Context(), &creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	ac.cookieTokensSetter.SetAll(c, tokens)
}

func (ac *AuthController) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(cookie.RefreshCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	tokens, err := ac.service.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	slog.Info("Tokens refreshed")
	ac.cookieTokensSetter.SetAll(c, tokens)
}

func (ac *AuthController) logout(c *gin.Context) {
	refreshToken, err := c.Cookie(cookie.RefreshCookieName)
	if err == nil {
		if err := ac.service.Logout(c.Request.Context(), refreshToken); err != nil {
			slog.Error("failed to logout", "error", err)
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		}

		ac.cookieTokensSetter.ResetRefresh(c)
	}

	ac.cookieTokensSetter.ResetAccess(c)

	c.Status(http.StatusOK)
}

func (ac *AuthController) updatePassword(c *gin.Context) {
	claims := ctxclaims.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var passwords entity.UserPasswords
	if err := c.ShouldBindJSON(&passwords); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := ac.service.UpdatePassword(c.Request.Context(), claims.UserID, &passwords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
