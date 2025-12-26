package artist_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type ArtistAssignController struct {
	artistAssignService artist.ArtistAssignService
}

func NewArtistAssignController(artistAssignService artist.ArtistAssignService) *ArtistAssignController {
	return &ArtistAssignController{
		artistAssignService: artistAssignService,
	}
}

func (c *ArtistAssignController) TracksMenu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Add track",
			Run:  c.addTrack,
		},
		{
			Name: "Remove track",
			Run:  c.removeTrack,
		},
		{
			Name: "View artist tracks",
			Run:  c.getArtistTracks,
		},
	}
}

func (c *ArtistAssignController) AlbumsMenu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Add album",
			Run:  c.addAlbum,
		},
		{
			Name: "Remove album",
			Run:  c.removeAlbum,
		},
		{
			Name: "View artist albums",
			Run:  c.getArtistAlbums,
		},
	}
}

func (c *ArtistAssignController) addAlbum(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	err = c.artistAssignService.AssignArtistToAlbum(ctx, session.Claims(), artistID, albumID)
	if err != nil {
		return fmt.Errorf("failed to assign artist: %w", err)
	}

	fmt.Println("Artist assigned successfully")
	return nil
}

func (c *ArtistAssignController) getArtistAlbums(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	albums, err := c.artistAssignService.GetArtistAlbums(ctx, artistID)
	if err != nil {
		return fmt.Errorf("failed to get artist albums: %w", err)
	}

	output.PrintAlbums(albums)
	return nil
}

func (c *ArtistAssignController) getArtistTracks(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	tracks, err := c.artistAssignService.GetArtistTracks(ctx, artistID)
	if err != nil {
		return fmt.Errorf("failed to get artist tracks: %w", err)
	}

	output.PrintTracks(tracks)
	return nil
}

func (c *ArtistAssignController) addTrack(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	trackIDStr := scanner.Text()
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	err = c.artistAssignService.AssignArtistToTrack(ctx, session.Claims(), artistID, trackID)
	if err != nil {
		return fmt.Errorf("failed to assign artist: %w", err)
	}

	fmt.Println("Artist assigned successfully")
	return nil
}

func (c *ArtistAssignController) removeTrack(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	trackIDStr := scanner.Text()
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	err = c.artistAssignService.UnassignArtistFromTrack(ctx, session.Claims(), artistID, trackID)
	if err != nil {
		return fmt.Errorf("failed to unassign artist: %w", err)
	}

	fmt.Println("Artist unassigned successfully")
	return nil
}

func (c *ArtistAssignController) removeAlbum(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	err = c.artistAssignService.UnassignArtistFromAlbum(ctx, session.Claims(), artistID, albumID)
	if err != nil {
		return fmt.Errorf("failed to unassign artist: %w", err)
	}

	fmt.Println("Artist unassigned successfully")
	return nil
}
