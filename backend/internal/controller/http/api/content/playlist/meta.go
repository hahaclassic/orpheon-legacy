package playlist_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/dto"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type PlaylistMetaController struct {
	playlistService playlist.PlaylistMetaService
	deleter         playlist.PlaylistDeletionService
	privacyService  playlist.PlaylistPrivacyChanger
	aggregator      playlist.PlaylistAggregator
}

func NewPlaylistMetaController(playlistService playlist.PlaylistMetaService,
	deleter playlist.PlaylistDeletionService,
	privacyService playlist.PlaylistPrivacyChanger,
	aggregator playlist.PlaylistAggregator) *PlaylistMetaController {
	return &PlaylistMetaController{
		playlistService: playlistService,
		deleter:         deleter,
		privacyService:  privacyService,
		aggregator:      aggregator,
	}
}

// GetPlaylist godoc
// @Summary Get playlist by ID
// @Description Get playlist details by its ID
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {object} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id} [get]
func (c *PlaylistMetaController) GetPlaylist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	playlist, err := c.playlistService.GetMeta(ctx.Request.Context(), claims, id)
	if err != nil {
		if errors.Is(err, commonerr.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		} else if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to get playlist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get playlist"})
		}
		return
	}

	aggregated, err := c.aggregator.GetPlaylists(ctx.Request.Context(), claims, playlist)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get playlist"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated[0])
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Create a new playlist with the provided details
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlist body entity.Playlist true "Playlist object"
// @Success 201 {object} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists [post]
func (c *PlaylistMetaController) CreatePlaylist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var playlist entity.PlaylistMeta
	if err := ctx.ShouldBindJSON(&playlist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.playlistService.CreateMeta(ctx.Request.Context(), claims, &playlist)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Playlist created successfully"})
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Update an existing playlist with new details
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Param playlist body entity.Playlist true "Updated playlist object"
// @Success 200 {object} entity.Playlist
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id} [put]
func (c *PlaylistMetaController) UpdatePlaylist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	var playlist entity.PlaylistMeta
	if err := ctx.ShouldBindJSON(&playlist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlist.ID = id
	err = c.playlistService.UpdateMeta(ctx.Request.Context(), claims, &playlist)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to update playlist"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist updated successfully"})
}

// DeletePlaylist godoc
// @Summary Delete a playlist
// @Description Delete a playlist by its ID
// @Tags playlists
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/playlists/{id} [delete]
func (c *PlaylistMetaController) DeletePlaylist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	err = c.deleter.DeletePlaylist(ctx.Request.Context(), claims, id)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to delete playlist"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PlaylistMetaController) GetMyPlaylists(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlists, err := c.playlistService.GetUserAllPlaylistsMeta(ctx.Request.Context(), claims, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user playlists"})
		return
	}

	aggregated, err := c.aggregator.GetPlaylists(ctx.Request.Context(), claims, playlists...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user playlists"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *PlaylistMetaController) GetUserPlaylists(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	playlists, err := c.playlistService.GetUserAllPlaylistsMeta(ctx.Request.Context(), claims, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user playlists"})
		return
	}

	aggregated, err := c.aggregator.GetPlaylists(ctx.Request.Context(), claims, playlists...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user playlists"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *PlaylistMetaController) UpdatePlaylistPrivacy(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	var privacy dto.PlaylistPrivacy
	if err := ctx.ShouldBindJSON(&privacy); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = c.privacyService.ChangePrivacy(ctx.Request.Context(), claims, id, privacy.IsPrivate)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to update playlist privacy"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update playlist privacy"})
		}
		return
	}

	ctx.Status(http.StatusOK)
}
