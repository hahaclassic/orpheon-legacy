package playlist_cli_ctrl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type PlaylistTrackController struct {
	playlistTrackService playlist.PlaylistTrackService
}

func NewPlaylistTrackController(playlistTrackService playlist.PlaylistTrackService) *PlaylistTrackController {
	return &PlaylistTrackController{
		playlistTrackService: playlistTrackService,
	}
}

func (c *PlaylistTrackController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Get playlist tracks",
			Run:  c.getTracks,
		},
		{
			Name: "Add track to playlist",
			Run:  c.addTrack,
		},
		{
			Name: "Remove track from playlist",
			Run:  c.removeTrack,
		},
		{
			Name: "Change track position",
			Run:  c.changeTrackPosition,
		},
	}
}

func (c *PlaylistTrackController) addTrack(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to add track to playlist.")
		return nil
	}

	var playlistID, trackID string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&playlistID); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	fmt.Print("Enter track ID: ")
	if _, err := fmt.Scan(&trackID); err != nil {
		return fmt.Errorf("failed to read track ID: %w", err)
	}

	playlistUUID, err := uuid.Parse(playlistID)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	trackUUID, err := uuid.Parse(trackID)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	err = c.playlistTrackService.AddTrack(ctx, session.Claims(), &entity.PlaylistTrack{
		PlaylistID: playlistUUID,
		TrackID:    trackUUID,
	})
	if err != nil {
		return fmt.Errorf("failed to add track to playlist: %w", err)
	}

	fmt.Println("Track added to playlist successfully")
	return nil
}

func (c *PlaylistTrackController) removeTrack(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("Login to remove track from playlist.")
		return nil
	}

	var playlistID, trackID string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&playlistID); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	fmt.Print("Enter track ID: ")
	if _, err := fmt.Scan(&trackID); err != nil {
		return fmt.Errorf("failed to read track ID: %w", err)
	}

	playlistUUID, err := uuid.Parse(playlistID)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	trackUUID, err := uuid.Parse(trackID)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	err = c.playlistTrackService.DeleteTrack(ctx, session.Claims(), &entity.PlaylistTrack{
		PlaylistID: playlistUUID,
		TrackID:    trackUUID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove track from playlist: %w", err)
	}

	fmt.Println("Track removed from playlist successfully")
	return nil
}

func (c *PlaylistTrackController) getTracks(ctx context.Context) error {
	var playlistID string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&playlistID); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	playlistUUID, err := uuid.Parse(playlistID)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	tracks, err := c.playlistTrackService.GetAllTracks(ctx, session.Claims(), playlistUUID)
	if err != nil {
		return fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	output.PrintTracks(tracks)
	return nil
}

func (c *PlaylistTrackController) changeTrackPosition(ctx context.Context) error {
	var playlistID, trackID, newPosition string
	fmt.Print("Enter playlist ID: ")
	if _, err := fmt.Scan(&playlistID); err != nil {
		return fmt.Errorf("failed to read playlist ID: %w", err)
	}

	playlistUUID, err := uuid.Parse(playlistID)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	fmt.Print("Enter track ID: ")
	if _, err := fmt.Scan(&trackID); err != nil {
		return fmt.Errorf("failed to read track ID: %w", err)
	}

	trackUUID, err := uuid.Parse(trackID)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	fmt.Print("Enter new position: ")
	if _, err := fmt.Scan(&newPosition); err != nil {
		return fmt.Errorf("failed to read new position: %w", err)
	}

	newPositionInt, err := strconv.Atoi(newPosition)
	if err != nil {
		return fmt.Errorf("failed to parse new position: %w", err)
	}

	err = c.playlistTrackService.ChangeTrackPosition(ctx, session.Claims(), &entity.PlaylistTrack{
		PlaylistID: playlistUUID,
		TrackID:    trackUUID,
		Position:   newPositionInt,
	})
	if err != nil {
		return fmt.Errorf("failed to change track position: %w", err)
	}

	fmt.Println("[OK] Track position changed successfully")
	return nil
}
