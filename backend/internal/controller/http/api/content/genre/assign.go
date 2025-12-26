package genre_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type GenreAssignController struct {
	genreAssignService genre.GenreAssignService
}

func NewGenreAssignController(genreAssignService genre.GenreAssignService) *GenreAssignController {
	return &GenreAssignController{genreAssignService: genreAssignService}
}

func (c *GenreAssignController) AssignGenreToAlbum(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	genreID, err := uuid.Parse(ctx.Param("genre_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	err = c.genreAssignService.AssignGenreToAlbum(ctx, claims, genreID, albumID)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Genre assigned to album"})
}

func (c *GenreAssignController) UnassignGenreFromAlbum(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	genreID, err := uuid.Parse(ctx.Param("genre_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	err = c.genreAssignService.UnassignGenreFromAlbum(ctx, claims, genreID, albumID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Genre unassigned from album"})
}
