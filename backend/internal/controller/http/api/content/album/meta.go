package album_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/dto"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type AlbumController struct {
	albumService album.AlbumMetaService
	aggregator   aggregator.ContentAggregator
}

func NewAlbumMetaController(albumService album.AlbumMetaService, aggregator aggregator.ContentAggregator) *AlbumController {
	return &AlbumController{
		albumService: albumService,
		aggregator:   aggregator,
	}
}

func (c *AlbumController) GetAlbum(ctx *gin.Context) {
	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	album, err := c.albumService.GetAlbum(ctx.Request.Context(), albumID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get album"})
		return
	}

	if album == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	aggregated, err := c.aggregator.GetAlbums(ctx.Request.Context(), album)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get album"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated[0])
}

func (c *AlbumController) GetAllAlbums(ctx *gin.Context) {
	albums, err := c.albumService.GetAllAlbums(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all albums"})
		return
	}

	aggregated, err := c.aggregator.GetAlbums(ctx.Request.Context(), albums...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate all albums"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *AlbumController) CreateAlbum(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var album entity.AlbumMeta
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.albumService.CreateAlbum(ctx.Request.Context(), claims, &album)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, dto.ID{ID: id})
}

func (c *AlbumController) UpdateAlbum(ctx *gin.Context) {
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

	var album entity.AlbumMeta
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	album.ID = albumID

	if err := c.albumService.UpdateAlbum(ctx.Request.Context(), claims, &album); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
}

func (c *AlbumController) DeleteAlbum(ctx *gin.Context) {
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

	if err := c.albumService.DeleteAlbum(ctx.Request.Context(), claims, albumID); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}
