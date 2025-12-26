package artist_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type ArtistAvatarController struct {
	artistAvatarService artist.ArtistAvatarService
}

func NewArtistAvatarController(artistAvatarService artist.ArtistAvatarService) *ArtistAvatarController {
	return &ArtistAvatarController{
		artistAvatarService: artistAvatarService,
	}
}

func (c *ArtistAvatarController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Upload avatar",
			Run:  c.uploadAvatar,
		},
		{
			Name: "Get avatar",
			Run:  c.getAvatar,
		},
		{
			Name: "Delete avatar",
			Run:  c.deleteAvatar,
		},
	}
}

func (c *ArtistAvatarController) uploadAvatar(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	fmt.Print("Enter path to avatar file: ")
	scanner.Scan()
	filePath := scanner.Text()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	cover := &entity.Cover{
		ObjectID: artistID,
		Data:     fileData,
	}

	err = c.artistAvatarService.UploadCover(ctx, session.Claims(), cover)
	if err != nil {
		return fmt.Errorf("failed to upload avatar: %w", err)
	}

	fmt.Println("Avatar uploaded successfully")
	return nil
}

func (c *ArtistAvatarController) getAvatar(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	fmt.Print("Enter path to save avatar: ")
	scanner.Scan()
	savePath := scanner.Text()

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	cover, err := c.artistAvatarService.GetCover(ctx, artistID)
	if err != nil {
		return fmt.Errorf("failed to get avatar: %w", err)
	}

	_, err = file.Write(cover.Data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("Avatar downloaded successfully")
	return nil
}

func (c *ArtistAvatarController) deleteAvatar(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	artistIDStr := scanner.Text()
	artistID, err := uuid.Parse(artistIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	err = c.artistAvatarService.DeleteCover(ctx, session.Claims(), artistID)
	if err != nil {
		return fmt.Errorf("failed to delete avatar: %w", err)
	}

	fmt.Println("Avatar deleted successfully")
	return nil
}
