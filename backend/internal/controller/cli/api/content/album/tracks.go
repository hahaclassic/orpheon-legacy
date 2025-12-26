package album_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type AlbumTrackController struct {
	albumTrackService album.AlbumTrackService
}

func NewAlbumTrackController(albumTrackService album.AlbumTrackService) *AlbumTrackController {
	return &AlbumTrackController{albumTrackService: albumTrackService}
}

func (c *AlbumTrackController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Tracks",
			Run:  c.GetAllTracks,
		},
	}
}

func (c *AlbumTrackController) GetAllTracks(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	tracks, err := c.albumTrackService.GetAllTracks(ctx, albumID)
	if err != nil {
		return fmt.Errorf("failed to get all tracks: %w", err)
	}

	output.PrintTracks(tracks)

	return nil
}
