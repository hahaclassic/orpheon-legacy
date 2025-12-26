package genre_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type GenreController struct {
	genreService usecase.GenreService
}

func NewGenreController(genreService usecase.GenreService) *GenreController {
	return &GenreController{
		genreService: genreService,
	}
}

func (c *GenreController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Create",
			Run:  c.createGenre,
		},
		{
			Name: "Get by ID",
			Run:  c.getGenreByID,
		},
		{
			Name: "List genres",
			Run:  c.getAllGenres,
		},
		{
			Name: "Update",
			Run:  c.updateGenre,
		},
		{
			Name: "Delete",
			Run:  c.deleteGenre,
		},
	}
}

func (c *GenreController) createGenre(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter genre name: ")
	scanner.Scan()
	name := scanner.Text()

	genre := &entity.Genre{
		Title: name,
	}

	claims := session.Claims()

	err := c.genreService.CreateGenre(ctx, claims, genre)
	if err != nil {
		return fmt.Errorf("failed to create genre: %w", err)
	}

	fmt.Println("Genre created successfully")
	return nil
}

func (c *GenreController) getGenreByID(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter genre ID: ")
	scanner.Scan()
	id := scanner.Text()

	genreID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse genre ID: %w", err)
	}

	genre, err := c.genreService.GetGenreByID(ctx, genreID)
	if err != nil {
		return fmt.Errorf("failed to get genre: %w", err)
	}

	output.PrintGenre(genre)
	return nil
}

func (c *GenreController) getAllGenres(ctx context.Context) error {
	genres, err := c.genreService.GetAllGenres(ctx)
	if err != nil {
		return fmt.Errorf("failed to list genres: %w", err)
	}

	output.PrintGenres(genres)

	return nil
}

func (c *GenreController) updateGenre(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter genre ID: ")
	scanner.Scan()
	id := scanner.Text()

	fmt.Print("Enter new genre title: ")
	scanner.Scan()
	title := scanner.Text()

	genreID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse genre ID: %w", err)
	}

	genre := &entity.Genre{
		ID:    genreID,
		Title: title,
	}

	claims := session.Claims()

	err = c.genreService.UpdateGenre(ctx, claims, genre)
	if err != nil {
		return fmt.Errorf("failed to update genre: %w", err)
	}

	fmt.Println("Genre updated successfully")
	return nil
}

func (c *GenreController) deleteGenre(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter genre ID: ")
	scanner.Scan()
	id := scanner.Text()

	genreID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse genre ID: %w", err)
	}

	claims := session.Claims()

	err = c.genreService.DeleteGenre(ctx, claims, genreID)
	if err != nil {
		return fmt.Errorf("failed to delete genre: %w", err)
	}

	fmt.Println("Genre deleted successfully")
	return nil
}
