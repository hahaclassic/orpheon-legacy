package playlist_cli_ctrl

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type PlaylistFavoriteController struct {
	playlistFavoriteService playlist.PlaylistFavoriteService
}

func NewPlaylistFavoriteController(playlistFavoriteService playlist.PlaylistFavoriteService) *PlaylistFavoriteController {
	return &PlaylistFavoriteController{
		playlistFavoriteService: playlistFavoriteService,
	}
}

func (c *PlaylistFavoriteController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Get favorite playlists",
			Run:  c.getFavoritePlaylists,
		},
		{
			Name: "Add playlist to favorites",
			Run:  c.addToFavorites,
		},
		{
			Name: "Remove playlist from favorites",
			Run:  c.removeFromFavorites,
		},
	}
}

func (c *PlaylistFavoriteController) addToFavorites(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to add playlist to favorites.")
		return nil
	}

	var id string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	playlistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	err = c.playlistFavoriteService.AddToUserFavorites(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to add playlist to favorites: %w", err)
	}

	fmt.Println("Playlist added to favorites successfully")
	return nil
}

func (c *PlaylistFavoriteController) removeFromFavorites(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to remove playlist from favorites.")
		return nil
	}

	var id string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	playlistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	err = c.playlistFavoriteService.DeleteFromUserFavorites(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to remove playlist from favorites: %w", err)
	}

	fmt.Println("Playlist removed from favorites successfully")
	return nil
}

func (c *PlaylistFavoriteController) getFavoritePlaylists(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to view favorite playlists.")
		return nil
	}

	playlists, err := c.playlistFavoriteService.GetUserFavorites(ctx, session.Claims())
	if err != nil {
		return fmt.Errorf("failed to get favorite playlists: %w", err)
	}

	output.PrintPlaylists(playlists)
	return nil
}
