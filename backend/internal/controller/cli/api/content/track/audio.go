package track_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type TrackAudioController struct {
	audioFileService track.AudioFileService
}

func NewTrackAudioController(audioFileService track.AudioFileService) *TrackAudioController {
	return &TrackAudioController{
		audioFileService: audioFileService,
	}
}

func (c *TrackAudioController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Upload",
			Run:  c.uploadAudio,
		},
		{
			Name: "Download",
			Run:  c.getAudio,
		},
		{
			Name: "Delete",
			Run:  c.deleteAudio,
		},
	}
}

func (c *TrackAudioController) uploadAudio(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	trackIDStr := scanner.Text()
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	fmt.Print("Enter path to audio file: ")
	scanner.Scan()
	filePath := scanner.Text()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	chunk := &entity.AudioChunk{
		TrackID: trackID,
		Start:   0,
		End:     int64(len(data)),
		Data:    data,
	}

	err = c.audioFileService.UploadAudioFile(ctx, session.Claims(), chunk)
	if err != nil {
		return fmt.Errorf("failed to upload audio: %w", err)
	}

	fmt.Println("Audio uploaded successfully")
	return nil
}

func (c *TrackAudioController) getAudio(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	trackIDStr := scanner.Text()
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	fmt.Print("Enter path to save audio file (.mp3): ")
	scanner.Scan()
	filePath := scanner.Text()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	chunk, err := c.audioFileService.GetAudioChunk(ctx, &entity.AudioChunk{
		TrackID: trackID,
		Start:   0,
		End:     int64(math.MaxInt64),
	})
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}

	_, err = file.Write(chunk.Data)
	if err != nil {
		return fmt.Errorf("failed to write audio to file: %w", err)
	}

	fmt.Println("Audio downloaded successfully")
	return nil
}

func (c *TrackAudioController) deleteAudio(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	trackIDStr := scanner.Text()
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	err = c.audioFileService.DeleteAudioFile(ctx, session.Claims(), trackID)
	if err != nil {
		return fmt.Errorf("failed to delete audio: %w", err)
	}

	fmt.Println("Audio deleted successfully")
	return nil
}
