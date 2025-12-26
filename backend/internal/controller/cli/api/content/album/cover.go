package album_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type AlbumCoverController struct {
	albumCoverService album.AlbumCoverService
}

func NewAlbumCoverController(albumCoverService album.AlbumCoverService) *AlbumCoverController {
	return &AlbumCoverController{
		albumCoverService: albumCoverService,
	}
}

func (c *AlbumCoverController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Get cover",
			Run:  c.getCover,
		},
		{
			Name: "Upload cover",
			Run:  c.uploadCover,
		},
		{
			Name: "Delete cover",
			Run:  c.deleteCover,
		},
	}
}

func (c *AlbumCoverController) getCover(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	fmt.Print("Enter path to save cover: ")
	scanner.Scan()
	savePath := scanner.Text()

	cover, err := c.albumCoverService.GetCover(ctx, albumID)
	if err != nil {
		return fmt.Errorf("failed to get cover: %w", err)
	}

	err = os.WriteFile(savePath, cover.Data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save cover: %w", err)
	}

	fmt.Println("Cover saved successfully")
	return nil
}

func (c *AlbumCoverController) uploadCover(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	fmt.Print("Enter cover file path: ")
	scanner.Scan()
	coverPath := scanner.Text()

	coverFile, err := os.Open(coverPath)
	if err != nil {
		return fmt.Errorf("failed to open cover file: %w", err)
	}
	defer coverFile.Close()

	coverData, err := io.ReadAll(coverFile)
	if err != nil {
		return fmt.Errorf("failed to read cover file: %w", err)
	}

	cover := &entity.Cover{
		ObjectID: albumID,
		Data:     coverData,
	}

	err = c.albumCoverService.UploadCover(ctx, session.Claims(), cover)
	if err != nil {
		return fmt.Errorf("failed to upload cover: %w", err)
	}

	fmt.Println("Cover uploaded successfully")
	return nil
}

func (c *AlbumCoverController) deleteCover(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	albumIDStr := scanner.Text()
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	err = c.albumCoverService.DeleteCover(ctx, session.Claims(), albumID)
	if err != nil {
		return fmt.Errorf("failed to delete cover: %w", err)
	}

	fmt.Println("Cover deleted successfully")
	return nil
}
