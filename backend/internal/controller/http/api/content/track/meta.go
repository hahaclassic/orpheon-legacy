package track_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/dto"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type TrackMetaController struct {
	trackService track.TrackMetaService
	aggregator   aggregator.ContentAggregator
}

func NewTrackMetaController(trackService track.TrackMetaService, aggregator aggregator.ContentAggregator) *TrackMetaController {
	return &TrackMetaController{
		trackService: trackService,
		aggregator:   aggregator,
	}
}

// GetTrack godoc
// @Summary Get track by ID
// @Description Get track details by its ID
// @Tags tracks
// @Produce json
// @Param id path string true "Track ID"
// @Success 200 {object} entity.Track
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks/{id} [get]
func (c *TrackMetaController) GetTrack(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	track, err := c.trackService.GetTrackMeta(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	aggregated, err := c.aggregator.GetTracks(ctx.Request.Context(), track)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aggregated[0])
}

// CreateTrack godoc
// @Summary Create a new track
// @Description Create a new track with the provided details
// @Tags tracks
// @Accept json
// @Produce json
// @Param track body entity.Track true "Track object"
// @Success 201 {object} entity.Track
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks [post]
func (c *TrackMetaController) CreateTrack(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var track entity.TrackMeta
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.trackService.CreateTrackMeta(ctx.Request.Context(), claims, &track)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create track"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, dto.ID{ID: id})
}

// UpdateTrack godoc
// @Summary Update a track
// @Description Update an existing track with new details
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Param track body entity.Track true "Updated track object"
// @Success 200 {object} entity.Track
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks/{id} [put]
func (c *TrackMetaController) UpdateTrack(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	var track entity.TrackMeta
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track.ID = id
	err = c.trackService.UpdateTrackMeta(ctx.Request.Context(), claims, &track)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update track"})
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteTrack godoc
// @Summary Delete a track
// @Description Delete a track by its ID
// @Tags tracks
// @Produce json
// @Param id path string true "Track ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks/{id} [delete]
func (c *TrackMetaController) DeleteTrack(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	err = c.trackService.DeleteTrackMeta(ctx.Request.Context(), claims, id)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete track"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
