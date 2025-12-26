package user_me_router

import (
	"github.com/gin-gonic/gin"
)

type PlaylistMetaController interface {
	GetMyPlaylists(c *gin.Context)
	GetUserPlaylists(c *gin.Context)
}

type UserController interface {
	GetMe(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateMe(c *gin.Context)
}

type PlaylistFavoritesController interface {
	GetFavoritePlaylists(c *gin.Context)
	AddToFavorites(c *gin.Context)
	RemoveFromFavorites(c *gin.Context)
}

type UserMeRouter struct {
	playlistMetaController      PlaylistMetaController
	userController              UserController
	playlistFavoritesController PlaylistFavoritesController
	authMiddleware              gin.HandlerFunc
}

func NewMeRouter(playlistMetaController PlaylistMetaController,
	userController UserController,
	playlistFavoritesController PlaylistFavoritesController,
	authMiddleware gin.HandlerFunc) *UserMeRouter {

	return &UserMeRouter{
		playlistMetaController:      playlistMetaController,
		userController:              userController,
		playlistFavoritesController: playlistFavoritesController,
		authMiddleware:              authMiddleware,
	}
}

func (r *UserMeRouter) RegisterRoutes(router *gin.RouterGroup) {
	me := router.Group("/me")
	me.Use(r.authMiddleware)
	{
		me.GET("/playlists", r.playlistMetaController.GetMyPlaylists)
		me.GET("/favorites", r.playlistFavoritesController.GetFavoritePlaylists)
		me.POST("/favorites/:playlist_id", r.playlistFavoritesController.AddToFavorites)
		me.DELETE("/favorites/:playlist_id", r.playlistFavoritesController.RemoveFromFavorites)
		me.GET("", r.userController.GetMe)
		me.PUT("", r.userController.UpdateMe)
	}

	user := router.Group("/users")
	{
		user.GET("/:id", r.userController.GetUser)
		user.GET("/:id/playlists", r.playlistMetaController.GetUserPlaylists)
	}
}
