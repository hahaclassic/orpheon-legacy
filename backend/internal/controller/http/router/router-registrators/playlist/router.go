package playlist_router

import "github.com/gin-gonic/gin"

type PlaylistMetaController interface {
	GetPlaylist(c *gin.Context)
	CreatePlaylist(c *gin.Context)
	UpdatePlaylist(c *gin.Context)
	DeletePlaylist(c *gin.Context)
	UpdatePlaylistPrivacy(c *gin.Context)
}

type PlaylistCoverController interface {
	GetCover(c *gin.Context)
	UploadCover(c *gin.Context)
	DeleteCover(c *gin.Context)
}

type PlaylistTrackController interface {
	GetPlaylistTracks(c *gin.Context)
	AddTrackToPlaylist(c *gin.Context)
	DeleteTrackFromPlaylist(c *gin.Context)
	ChangeTrackPosition(c *gin.Context)
}

type PlaylistRouter struct {
	playlistMetaController  PlaylistMetaController
	playlistTrackController PlaylistTrackController
	playlistCoverController PlaylistCoverController
	authMiddleware          gin.HandlerFunc
}

func NewPlaylistRouter(
	playlistMetaController PlaylistMetaController,
	playlistTrackController PlaylistTrackController,
	playlistCoverController PlaylistCoverController,
	authMiddleware gin.HandlerFunc,
) *PlaylistRouter {
	return &PlaylistRouter{
		playlistMetaController:  playlistMetaController,
		playlistTrackController: playlistTrackController,
		playlistCoverController: playlistCoverController,
		authMiddleware:          authMiddleware,
	}
}

func (r *PlaylistRouter) RegisterRoutes(router *gin.RouterGroup) {
	playlistGroup := router.Group("/playlists")
	playlistGroup.Use(r.authMiddleware)
	{
		playlistGroup.GET("/:id", r.playlistMetaController.GetPlaylist)
		playlistGroup.POST("", r.playlistMetaController.CreatePlaylist)
		playlistGroup.PUT("/:id", r.playlistMetaController.UpdatePlaylist)
		playlistGroup.DELETE("/:id", r.playlistMetaController.DeletePlaylist)
		playlistGroup.PATCH("/:id/privacy", r.playlistMetaController.UpdatePlaylistPrivacy)
	}

	coverGroup := playlistGroup.Group("/:id/cover")
	{
		coverGroup.GET("", r.playlistCoverController.GetCover)
		coverGroup.POST("", r.playlistCoverController.UploadCover)
		coverGroup.DELETE("", r.playlistCoverController.DeleteCover)
	}

	tracksGroup := playlistGroup.Group("/:id/tracks")
	{
		tracksGroup.GET("", r.playlistTrackController.GetPlaylistTracks)
		tracksGroup.POST("", r.playlistTrackController.AddTrackToPlaylist)
		tracksGroup.DELETE("/:track_id", r.playlistTrackController.DeleteTrackFromPlaylist)
		tracksGroup.PATCH("/:track_id/position", r.playlistTrackController.ChangeTrackPosition)
	}
}
