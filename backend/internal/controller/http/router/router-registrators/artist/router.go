package artist_router

import (
	"github.com/gin-gonic/gin"
)

type ArtistController interface {
	GetArtist(c *gin.Context)
	GetAllArtists(c *gin.Context)
	CreateArtist(c *gin.Context)
	UpdateArtist(c *gin.Context)
	DeleteArtist(c *gin.Context)
}

type ArtistAvatarController interface {
	GetAvatar(c *gin.Context)
	UploadAvatar(c *gin.Context)
	DeleteAvatar(c *gin.Context)
}

type ArtistAssignController interface {
	GetAlbumsByArtist(c *gin.Context)
	GetTracksByArtist(c *gin.Context)
	AssignArtistToTrack(c *gin.Context)
	UnassignArtistFromTrack(c *gin.Context)
	AssignArtistToAlbum(c *gin.Context)
	UnassignArtistFromAlbum(c *gin.Context)
}

type ArtistRouter struct {
	artistController       ArtistController
	artistAvatarController ArtistAvatarController
	artistAssignController ArtistAssignController
	authMiddleware         gin.HandlerFunc
}

func NewArtistRouter(
	artistController ArtistController,
	artistAvatarController ArtistAvatarController,
	artistAssignController ArtistAssignController,
	authMiddleware gin.HandlerFunc,
) *ArtistRouter {
	return &ArtistRouter{
		artistController:       artistController,
		artistAvatarController: artistAvatarController,
		artistAssignController: artistAssignController,
		authMiddleware:         authMiddleware,
	}
}

func (r *ArtistRouter) RegisterRoutes(router *gin.RouterGroup) {
	artistGroup := router.Group("/artists")

	artistGroup.GET("", r.artistController.GetAllArtists)
	artistGroup.GET("/:id", r.artistController.GetArtist)
	artistGroup.GET("/:id/albums", r.artistAssignController.GetAlbumsByArtist)
	artistGroup.GET("/:id/tracks", r.artistAssignController.GetTracksByArtist)

	artistProtected := artistGroup.Group("")
	artistProtected.Use(r.authMiddleware)
	{
		artistProtected.POST("", r.artistController.CreateArtist)
		artistProtected.PUT("/:id", r.artistController.UpdateArtist)
		artistProtected.DELETE("/:id", r.artistController.DeleteArtist)

		artistProtected.POST("/:id/tracks/:track_id", r.artistAssignController.AssignArtistToTrack)
		artistProtected.DELETE("/:id/tracks/:track_id", r.artistAssignController.UnassignArtistFromTrack)
		artistProtected.POST("/:id/albums/:album_id", r.artistAssignController.AssignArtistToAlbum)
		artistProtected.DELETE("/:id/albums/:album_id", r.artistAssignController.UnassignArtistFromAlbum)
	}

	avatarGroup := artistGroup.Group("/:id/avatar")
	avatarGroup.GET("", r.artistAvatarController.GetAvatar)

	avatarProtected := avatarGroup.Group("")
	avatarProtected.Use(r.authMiddleware)
	{
		avatarProtected.POST("", r.artistAvatarController.UploadAvatar)
		avatarProtected.DELETE("", r.artistAvatarController.DeleteAvatar)
	}
}
