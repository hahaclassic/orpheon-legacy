package track_ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
)

type TrackSegmentController struct {
	service track.TrackSegmentService
}

func NewTrackSegmentController(service track.TrackSegmentService) *TrackSegmentController {
	return &TrackSegmentController{service: service}
}

func (c *TrackSegmentController) GetSegments(ctx *gin.Context) {
	trackID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	segments, err := c.service.GetSegments(ctx.Request.Context(), trackID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, segments)
}
