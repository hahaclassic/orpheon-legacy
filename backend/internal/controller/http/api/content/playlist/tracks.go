package playlist_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/dto"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type PlaylistTrackController struct {
	tracksService playlist.PlaylistTrackService
	aggregator    aggregator.ContentAggregator
}

func NewPlaylistTrackController(tracksService playlist.PlaylistTrackService,
	aggregator aggregator.ContentAggregator) *PlaylistTrackController {
	return &PlaylistTrackController{
		tracksService: tracksService,
		aggregator:    aggregator,
	}
}

func (c *PlaylistTrackController) GetPlaylistTracks(ctx *gin.Context) {
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

	tracks, err := c.tracksService.GetAllTracks(ctx.Request.Context(), claims, id)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get playlist tracks"})
		}
		return
	}

	aggregated, err := c.aggregator.GetTracks(ctx.Request.Context(), tracks...)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get playlist tracks"})
		}
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *PlaylistTrackController) AddTrackToPlaylist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	var request dto.TrackAdditionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.tracksService.AddTrack(ctx.Request.Context(), claims, &entity.PlaylistTrack{
		PlaylistID: playlistID,
		TrackID:    request.TrackID,
	})
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add track to playlist"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Track added to playlist"})
}

func (c *PlaylistTrackController) DeleteTrackFromPlaylist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	trackID, err := uuid.Parse(ctx.Param("track_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	err = c.tracksService.DeleteTrack(ctx.Request.Context(), claims, &entity.PlaylistTrack{
		PlaylistID: playlistID,
		TrackID:    trackID,
	})
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove track from playlist"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Track removed from playlist"})
}

func (c *PlaylistTrackController) ChangeTrackPosition(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	trackID, err := uuid.Parse(ctx.Param("track_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	var request dto.PlaylistTrackPosition
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.tracksService.ChangeTrackPosition(ctx.Request.Context(), claims, &entity.PlaylistTrack{
		PlaylistID: playlistID,
		TrackID:    trackID,
		Position:   request.Position,
	})
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change track position"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Track position changed"})
}
