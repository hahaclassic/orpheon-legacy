package playlist_ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
)

type PlaylistFavoritesController struct {
	favoritesService playlist.PlaylistFavoriteService
	aggregator       playlist.PlaylistAggregator
}

func NewPlaylistFavoritesController(favoritesService playlist.PlaylistFavoriteService,
	aggregator playlist.PlaylistAggregator) *PlaylistFavoritesController {
	return &PlaylistFavoritesController{
		favoritesService: favoritesService,
		aggregator:       aggregator,
	}
}

// GetFavoritePlaylists godoc
// @Summary Get favorite playlists
// @Description Get all favorite playlists for the current user
// @Tags playlists
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Playlist
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/favorites [get]
func (c *PlaylistFavoritesController) GetFavoritePlaylists(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlists, err := c.favoritesService.GetUserFavorites(ctx.Request.Context(), claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorite playlists"})
		return
	}

	aggregated, err := c.aggregator.GetPlaylists(ctx.Request.Context(), claims, playlists...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorite playlists"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *PlaylistFavoritesController) AddToFavorites(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(ctx.Param("playlist_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	err = c.favoritesService.AddToUserFavorites(ctx.Request.Context(), claims, playlistID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to add playlist to favorites"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *PlaylistFavoritesController) RemoveFromFavorites(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(ctx.Param("playlist_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	err = c.favoritesService.DeleteFromUserFavorites(ctx.Request.Context(), claims, playlistID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to remove playlist from favorites"})
		return
	}

	ctx.Status(http.StatusOK)
}
