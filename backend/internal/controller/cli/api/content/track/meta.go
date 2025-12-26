package track_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type TrackMetaController struct {
	trackMetaService track.TrackMetaService
}

func NewTrackMetaController(trackMetaService track.TrackMetaService) *TrackMetaController {
	return &TrackMetaController{
		trackMetaService: trackMetaService,
	}
}

func (c *TrackMetaController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Create",
			Run:  c.createTrackMeta,
		},
		{
			Name: "Get by ID",
			Run:  c.getTrackMetaByID,
		},
		{
			Name: "Update",
			Run:  c.updateTrackMeta,
		},
		{
			Name: "Delete",
			Run:  c.deleteTrackMeta,
		},
	}
}

func (c *TrackMetaController) createTrackMeta(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter track duration (in seconds): ")
	var duration int
	if _, err := fmt.Scan(&duration); err != nil {
		return fmt.Errorf("failed to read duration: %w", err)
	}

	fmt.Print("Enter track genre ID: ")
	scanner.Scan()
	genreIDStr := scanner.Text()
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse genre ID: %w", err)
	}

	fmt.Print("Is track explicit? (y/n): ")
	var explicit string
	if _, err := fmt.Scan(&explicit); err != nil {
		return fmt.Errorf("failed to read explicit flag: %w", err)
	}

	fmt.Print("Enter license ID (leave empty for none): ")
	scanner.Scan()
	licenseIDStr := scanner.Text()
	licenseID, err := uuid.Parse(licenseIDStr)
	if err != nil {
		licenseID = uuid.Nil
	}

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		albumID = uuid.Nil
	}

	track := &entity.TrackMeta{
		Name:      name,
		GenreID:   genreID,
		Duration:  duration,
		Explicit:  explicit == "y",
		LicenseID: licenseID,
		AlbumID:   albumID,
	}

	trackID, err := c.trackMetaService.CreateTrackMeta(ctx, session.Claims(), track)
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	fmt.Printf("[OK] Track created successfully with ID: %s\n", trackID)
	return nil
}

func (c *TrackMetaController) getTrackMetaByID(ctx context.Context) error {
	var id string

	fmt.Print("Enter track ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read track ID: %w", err)
	}

	trackID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	track, err := c.trackMetaService.GetTrackMeta(ctx, trackID)
	if err != nil {
		return fmt.Errorf("failed to get track: %w", err)
	}

	output.PrintTrack(track)
	return nil
}

func (c *TrackMetaController) updateTrackMeta(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	id := scanner.Text()

	fmt.Print("Enter new track name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter new track duration (in seconds): ")
	var duration int
	if _, err := fmt.Scan(&duration); err != nil {
		return fmt.Errorf("failed to read duration: %w", err)
	}

	fmt.Print("Is track explicit? (true/false): ")
	var explicit bool
	if _, err := fmt.Scan(&explicit); err != nil {
		return fmt.Errorf("failed to read explicit flag: %w", err)
	}

	fmt.Print("Enter new license ID (leave empty for none): ")
	scanner.Scan()
	licenseIDStr := scanner.Text()
	licenseID, err := uuid.Parse(licenseIDStr)
	if err != nil {
		licenseID = uuid.Nil
	}

	fmt.Print("Enter new album ID (leave empty for none): ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		albumID = uuid.Nil
	}

	trackID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	track := &entity.TrackMeta{
		ID:        trackID,
		Name:      name,
		Duration:  duration,
		Explicit:  explicit,
		LicenseID: licenseID,
		AlbumID:   albumID,
	}

	claims := session.Claims()

	err = c.trackMetaService.UpdateTrackMeta(ctx, claims, track)
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	fmt.Println("Track updated successfully")
	return nil
}

func (c *TrackMetaController) deleteTrackMeta(ctx context.Context) error {
	var id string

	fmt.Print("Enter track ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read track ID: %w", err)
	}

	trackID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	claims := session.Claims()

	err = c.trackMetaService.DeleteTrackMeta(ctx, claims, trackID)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	fmt.Println("Track deleted successfully")
	return nil
}
