package album_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type AlbumMainController struct {
	albumMetaService  usecase.AlbumMetaService
	albumCoverService usecase.AlbumCoverService
}

func NewAlbumMainController(albumMetaService usecase.AlbumMetaService, albumCoverService usecase.AlbumCoverService) *AlbumMainController {
	return &AlbumMainController{
		albumMetaService:  albumMetaService,
		albumCoverService: albumCoverService,
	}
}

func (c *AlbumMainController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Create",
			Run:  c.createAlbum,
		},
		{
			Name: "Get by ID",
			Run:  c.getAlbumByID,
		},
		{
			Name: "List albums",
			Run:  c.getAllAlbums,
		},
		// {
		// 	Name: "Update",
		// 	Run:  c.updateAlbum,
		// },
		// {
		// 	Name: "Delete",
		// 	Run:  c.deleteAlbum,
		// },
	}
}

func (c *AlbumMainController) createAlbum(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter album label: ")
	scanner.Scan()
	label := scanner.Text()

	fmt.Print("Enter license ID: ")
	scanner.Scan()
	licenseIDStr := scanner.Text()
	licenseID, err := uuid.Parse(licenseIDStr)
	if err != nil {
		licenseID = uuid.Nil
	}

	fmt.Print("Enter release date (YYYY-MM-DD): ")
	scanner.Scan()
	releaseDateStr := scanner.Text()
	releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
	if err != nil {
		return fmt.Errorf("failed to parse release date: %w", err)
	}

	album := &entity.AlbumMeta{
		Title:       title,
		Label:       label,
		LicenseID:   licenseID,
		ReleaseDate: releaseDate,
	}

	id, err := c.albumMetaService.CreateAlbum(ctx, session.Claims(), album)
	if err != nil {
		return fmt.Errorf("failed to create album: %w", err)
	}

	fmt.Printf("Album created successfully with ID: %s\n", id)
	return nil
}

func (c *AlbumMainController) getAlbumByID(ctx context.Context) error {
	var id string

	fmt.Print("Enter album ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read album ID: %w", err)
	}

	albumID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	album, err := c.albumMetaService.GetAlbum(ctx, albumID)
	if err != nil {
		return fmt.Errorf("failed to get album: %w", err)
	}

	output.PrintAlbum(album)
	return nil
}

func (c *AlbumMainController) getAllAlbums(ctx context.Context) error {
	albums, err := c.albumMetaService.GetAllAlbums(ctx)
	if err != nil {
		return fmt.Errorf("failed to list albums: %w", err)
	}

	output.PrintAlbums(albums)
	return nil
}

// func (c *AlbumController) updateAlbum(ctx context.Context) error {
// 	scanner := bufio.NewScanner(os.Stdin)

// 	fmt.Print("Enter album ID: ")
// 	scanner.Scan()
// 	id := scanner.Text()

// 	fmt.Print("Enter new album title: ")
// 	scanner.Scan()
// 	title := scanner.Text()

// 	fmt.Print("Enter new album label: ")
// 	scanner.Scan()
// 	label := scanner.Text()

// 	fmt.Print("Enter new license ID: ")
// 	scanner.Scan()
// 	licenseIDStr := scanner.Text()
// 	licenseID, err := uuid.Parse(licenseIDStr)
// 	if err != nil {

// 		fmt.Print("Enter new album cover path (leave empty to keep current): ")
// 		scanner.Scan()
// 		coverPath := scanner.Text()

// 		albumID, err := uuid.Parse(id)
// 		if err != nil {
// 			return fmt.Errorf("failed to parse album ID: %w", err)
// 		}

// 		album := &entity.AlbumMeta{
// 			ID:          albumID,
// 			Title:       title,
// 			Label:       label,
// 			LicenseID:   licenseID,
// 			ReleaseDate: releaseDate,
// 		}

// 		claims := session.Claims()

// 		err = c.albumMetaService.UpdateAlbumMeta(ctx, claims, album)
// 		if err != nil {
// 			return fmt.Errorf("failed to update album: %w", err)
// 		}

// 		if coverPath != "" {
// 			coverFile, err := os.Open(coverPath)
// 			if err != nil {
// 				return fmt.Errorf("failed to open cover file: %w", err)
// 			}
// 			defer coverFile.Close()

// 			err = c.albumCoverService.SaveCover(ctx, claims, album.ID, coverFile)
// 			if err != nil {
// 				return fmt.Errorf("failed to save album cover: %w", err)
// 			}
// 		}

// 		fmt.Println("Album updated successfully")
// 		return nil
// 	}

// func (c *AlbumController) deleteAlbum(ctx context.Context) error {
// 	var id string

// 	fmt.Print("Enter album ID: ")
// 	if _, err := fmt.Scan(&id); err != nil {
// 		return fmt.Errorf("failed to read album ID: %w", err)
// 	}

// 	albumID, err := uuid.Parse(id)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse album ID: %w", err)
// 	}

// 	claims := session.Claims()

// 	err = c.albumMetaService.DeleteAlbum(ctx, claims, albumID)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete album: %w", err)
// 	}

// 	fmt.Println("Album deleted successfully")
// return nil
// }
