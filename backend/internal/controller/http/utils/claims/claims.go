package ctxclaims

import (
	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

func GetClaims(c *gin.Context) *entity.Claims {
	claims, exists := c.Get("claims")
	if !exists {
		return nil
	}

	parsedClaims, ok := claims.(*entity.Claims)
	if !ok {
		return nil
	}

	return parsedClaims
}
