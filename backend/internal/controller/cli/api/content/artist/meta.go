package artist_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type ArtistMetaController struct {
	artistService usecase.ArtistMetaService
}

func NewArtistMetaController(artistService usecase.ArtistMetaService) *ArtistMetaController {
	return &ArtistMetaController{
		artistService: artistService,
	}
}

func (c *ArtistMetaController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Create",
			Run:  c.createArtistMeta,
		},
		{
			Name: "Get by ID",
			Run:  c.getArtistMetaByID,
		},
		{
			Name: "List artists",
			Run:  c.getAllArtistMeta,
		},
		{
			Name: "Update",
			Run:  c.updateArtistMeta,
		},
		{
			Name: "Delete",
			Run:  c.deleteArtistMeta,
		},
	}
}

func (c *ArtistMetaController) createArtistMeta(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter artist description: ")
	scanner.Scan()
	description := scanner.Text()

	fmt.Print("Enter artist country: ")
	scanner.Scan()
	country := scanner.Text()

	artist := &entity.ArtistMeta{
		Name:        name,
		Description: description,
		Country:     country,
	}

	claims := session.Claims()

	err := c.artistService.CreateArtistMeta(ctx, claims, artist)
	if err != nil {
		return fmt.Errorf("failed to create artist: %w", err)
	}

	fmt.Println("Artist created successfully")
	return nil
}

func (c *ArtistMetaController) getArtistMetaByID(ctx context.Context) error {
	var id string

	fmt.Print("Enter artist ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read artist ID: %w", err)
	}

	artistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	artist, err := c.artistService.GetArtistMeta(ctx, artistID)
	if err != nil {
		return fmt.Errorf("failed to get artist: %w", err)
	}

	output.PrintArtist(artist)
	return nil
}

func (c *ArtistMetaController) getAllArtistMeta(ctx context.Context) error {
	artists, err := c.artistService.GetAllArtistMeta(ctx)
	if err != nil {
		return fmt.Errorf("failed to list artists: %w", err)
	}

	output.PrintArtists(artists)

	return nil
}

func (c *ArtistMetaController) updateArtistMeta(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter artist ID: ")
	scanner.Scan()
	id := scanner.Text()

	fmt.Print("Enter new artist title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter new artist description: ")
	scanner.Scan()
	description := scanner.Text()

	artistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	artist := &entity.ArtistMeta{
		ID:          artistID,
		Name:        title,
		Description: description,
	}

	claims := session.Claims()

	err = c.artistService.UpdateArtistMeta(ctx, claims, artist)
	if err != nil {
		return fmt.Errorf("failed to update artist: %w", err)
	}

	fmt.Println("Artist updated successfully")
	return nil
}

func (c *ArtistMetaController) deleteArtistMeta(ctx context.Context) error {
	var id string

	fmt.Print("Enter artist ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read artist ID: %w", err)
	}

	artistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse artist ID: %w", err)
	}

	claims := session.Claims()

	err = c.artistService.DeleteArtistMeta(ctx, claims, artistID)
	if err != nil {
		return fmt.Errorf("failed to delete artist: %w", err)
	}

	fmt.Println("Artist deleted successfully")
	return nil
}
