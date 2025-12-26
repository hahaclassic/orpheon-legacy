package playlist_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type PlaylistCoverController struct {
	playlistCoverService playlist.PlaylistCoverService
}

func NewPlaylistCoverController(playlistCoverService playlist.PlaylistCoverService) *PlaylistCoverController {
	return &PlaylistCoverController{
		playlistCoverService: playlistCoverService,
	}
}

func (c *PlaylistCoverController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Upload cover",
			Run:  c.uploadCover,
		},
		{
			Name: "Download cover",
			Run:  c.getCover,
		},
		{
			Name: "Delete cover",
			Run:  c.deleteCover,
		},
	}
}

func (c *PlaylistCoverController) uploadCover(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to upload playlist cover.")
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

	fmt.Print("Enter path to cover file: ")
	scanner.Scan()
	coverPath := scanner.Text()

	coverData, err := os.ReadFile(coverPath)
	if err != nil {
		return fmt.Errorf("failed to read cover file: %w", err)
	}

	cover := &entity.Cover{
		ObjectID: playlistID,
		Data:     coverData,
	}

	err = c.playlistCoverService.UploadCover(ctx, session.Claims(), cover)
	if err != nil {
		return fmt.Errorf("failed to upload playlist cover: %w", err)
	}

	fmt.Println("Playlist cover uploaded successfully")
	return nil
}

func (c *PlaylistCoverController) getCover(ctx context.Context) error {
	var id string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	playlistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	cover, err := c.playlistCoverService.GetCover(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist cover: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter path to save cover: ")
	scanner.Scan()
	savePath := scanner.Text()

	err = os.WriteFile(savePath, cover.Data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save cover: %w", err)
	}

	fmt.Println("Cover saved successfully to:", savePath)
	return nil
}

func (c *PlaylistCoverController) deleteCover(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to delete playlist cover.")
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

	err = c.playlistCoverService.DeleteCover(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to delete playlist cover: %w", err)
	}

	fmt.Println("Playlist cover deleted successfully")
	return nil
}
