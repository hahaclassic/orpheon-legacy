package playlist_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type PlaylistMetaController struct {
	playlistService playlist.PlaylistMetaService
	privacyChanger  playlist.PlaylistPrivacyChanger
	deleter         playlist.PlaylistDeletionService
}

func NewPlaylistMetaController(playlistService playlist.PlaylistMetaService,
	privacyChanger playlist.PlaylistPrivacyChanger,
	deleter playlist.PlaylistDeletionService) *PlaylistMetaController {
	return &PlaylistMetaController{
		playlistService: playlistService,
		privacyChanger:  privacyChanger,
		deleter:         deleter,
	}
}

func (c *PlaylistMetaController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Get my playlists",
			Run:  c.getMyPlaylists,
		},
		{
			Name: "Get user playlists",
			Run:  c.getUserPlaylists,
		},
		{
			Name: "Get playlist by ID",
			Run:  c.getPlaylist,
		},
		{
			Name: "Create playlist",
			Run:  c.createPlaylist,
		},
		{
			Name: "Update playlist info",
			Run:  c.updatePlaylist,
		},
		{
			Name: "Change playlist privacy",
			Run:  c.changePlaylistPrivacy,
		},
		{
			Name: "Delete playlist",
			Run:  c.deletePlaylist,
		},
	}
}

func (c *PlaylistMetaController) createPlaylist(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to create playlist.")
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter playlist name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter playlist description: ")
	scanner.Scan()
	description := scanner.Text()

	var isPrivate bool
	fmt.Print("Is playlist private? (y/n): ")
	scanner.Scan()
	if scanner.Text() == "y" {
		isPrivate = true
	}
	fmt.Println(isPrivate)

	playlist := &entity.PlaylistMeta{
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
	}

	err := c.playlistService.CreateMeta(ctx, session.Claims(), playlist)
	if err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}

	fmt.Println("Playlist created successfully")
	return nil
}

func (c *PlaylistMetaController) getMyPlaylists(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to get your playlists.")
		return nil
	}

	playlists, err := c.playlistService.GetUserAllPlaylistsMeta(ctx, session.Claims(), session.Claims().UserID)
	if err != nil {
		return fmt.Errorf("failed to get user playlists: %w", err)
	}
	output.PrintPlaylists(playlists)
	return nil
}

func (c *PlaylistMetaController) getUserPlaylists(ctx context.Context) error {
	var id string

	fmt.Print("Enter user ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read user ID: %w", err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	playlists, err := c.playlistService.GetUserAllPlaylistsMeta(ctx, session.Claims(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user playlists: %w", err)
	}

	output.PrintPlaylists(playlists)

	return nil
}

func (c *PlaylistMetaController) getPlaylist(ctx context.Context) error {
	var id string

	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	playlistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	playlist, err := c.playlistService.GetMeta(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	output.PrintPlaylist(playlist)
	return nil
}

func (c *PlaylistMetaController) updatePlaylist(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to update playlist.")
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

	playlist, err := c.playlistService.GetMeta(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter new playlist name (leave empty to skip): ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter new playlist description (leave empty to skip): ")
	scanner.Scan()
	description := scanner.Text()

	if name != "" {
		playlist.Name = name
	}

	if description != "" {
		playlist.Description = description
	}

	err = c.playlistService.UpdateMeta(ctx, session.Claims(), playlist)
	if err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}

	fmt.Println("Playlist updated successfully")
	return nil
}

func (c *PlaylistMetaController) changePlaylistPrivacy(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to change playlist privacy.")
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

	scanner := bufio.NewScanner(os.Stdin)
	var isPrivate bool
	fmt.Print("Is playlist private? (y/n): ")
	scanner.Scan()
	if scanner.Text() == "y" {
		isPrivate = true
	}

	err = c.privacyChanger.ChangePrivacy(ctx, session.Claims(), playlistID, isPrivate)
	if err != nil {
		return fmt.Errorf("failed to change playlist privacy: %w", err)
	}

	fmt.Println("Playlist privacy changed successfully")
	return nil
}

func (c *PlaylistMetaController) deletePlaylist(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to delete playlist.")
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

	err = c.deleter.DeletePlaylist(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	fmt.Println("Playlist deleted successfully")
	return nil
}
