package album_ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
)

type AlbumTrackController struct {
	albumTrackService album.AlbumTrackService
	aggregator        aggregator.ContentAggregator
}

func NewAlbumTrackController(albumTrackService album.AlbumTrackService, aggregator aggregator.ContentAggregator) *AlbumTrackController {
	return &AlbumTrackController{
		albumTrackService: albumTrackService,
		aggregator:        aggregator,
	}
}

func (c *AlbumTrackController) GetAlbumTracks(ctx *gin.Context) {
	albumID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	tracks, err := c.albumTrackService.GetAllTracks(ctx.Request.Context(), albumID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	aggregated, err := c.aggregator.GetTracks(ctx.Request.Context(), tracks...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}
