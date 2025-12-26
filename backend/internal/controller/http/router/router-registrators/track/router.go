package track_router

import (
	"github.com/gin-gonic/gin"
)

type TrackMetaController interface {
	GetTrack(c *gin.Context)
	CreateTrack(c *gin.Context)
	UpdateTrack(c *gin.Context)
	DeleteTrack(c *gin.Context)
}

type TrackSegmentController interface {
	GetSegments(c *gin.Context)
}

type TrackAudioController interface {
	GetAudioChunk(c *gin.Context)
	UploadAudioFile(c *gin.Context)
	DeleteAudioFile(c *gin.Context)
}

type ArtistAssignController interface {
	AssignArtistToTrack(c *gin.Context)
	UnassignArtistFromTrack(c *gin.Context)
}

type StatController interface {
	UpdateStat(c *gin.Context)
}

type TrackRouter struct {
	trackMetaController    TrackMetaController
	segmentService         TrackSegmentController
	audioService           TrackAudioController
	artistAssignController ArtistAssignController
	statController         StatController
	authMiddleware         gin.HandlerFunc
}

func NewTrackRouter(trackMetaController TrackMetaController,
	segmentService TrackSegmentController,
	audioService TrackAudioController,
	statController StatController,
	artistAssignController ArtistAssignController,
	authMiddleware gin.HandlerFunc) *TrackRouter {
	return &TrackRouter{
		trackMetaController:    trackMetaController,
		segmentService:         segmentService,
		audioService:           audioService,
		statController:         statController,
		artistAssignController: artistAssignController,
		authMiddleware:         authMiddleware,
	}
}

func (r *TrackRouter) RegisterRoutes(router *gin.RouterGroup) {
	tracks := router.Group("/tracks")
	{
		tracks.GET("/:id", r.trackMetaController.GetTrack)
		tracks.GET("/:id/segments", r.segmentService.GetSegments)

		tracksProtected := tracks.Group("")
		tracksProtected.Use(r.authMiddleware)
		{
			tracksProtected.POST("", r.trackMetaController.CreateTrack)
			tracksProtected.PUT("/:id", r.trackMetaController.UpdateTrack)
			tracksProtected.DELETE("/:id", r.trackMetaController.DeleteTrack)
			tracksProtected.POST("/:id/stats", r.statController.UpdateStat)
		}

		tracksAudio := tracks.Group("/:id/audio")
		{
			tracksAudio.GET("", r.audioService.GetAudioChunk)
			tracksAudioProtected := tracksAudio.Group("")
			tracksAudioProtected.Use(r.authMiddleware)
			{
				tracksAudioProtected.POST("", r.audioService.UploadAudioFile)
				tracksAudioProtected.DELETE("", r.audioService.DeleteAudioFile)
			}
		}
	}
}
