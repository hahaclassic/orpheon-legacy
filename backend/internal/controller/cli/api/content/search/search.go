package search_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type SearchController struct {
	searchService search.SearchService
}

func NewSearchController(searchService search.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

func (c *SearchController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Search tracks",
			Run:  c.searchTracks,
		},
		{
			Name: "Search albums",
			Run:  c.searchAlbums,
		},
		{
			Name: "Search artists",
			Run:  c.searchArtists,
		},
		{
			Name: "Search playlists",
			Run:  c.searchPlaylists,
		},
	}
}

func (c *SearchController) searchTracks(ctx context.Context) (err error) {
	req, err := c.getRequest()
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	tracks, err := c.searchService.SearchTracks(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to search tracks: %w", err)
	}

	output.PrintTracks(tracks)
	return nil
}

func (c *SearchController) searchAlbums(ctx context.Context) (err error) {
	req, err := c.getRequest()
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	albums, err := c.searchService.SearchAlbums(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to search albums: %w", err)
	}

	output.PrintAlbums(albums)
	return nil
}

func (c *SearchController) searchArtists(ctx context.Context) error {
	req, err := c.getRequest()
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	artists, err := c.searchService.SearchArtists(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to search artists: %w", err)
	}

	output.PrintArtists(artists)
	return nil
}

func (c *SearchController) searchPlaylists(ctx context.Context) error {
	req, err := c.getRequest()
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	playlists, err := c.searchService.SearchPlaylists(ctx, session.Claims(), req)
	if err != nil {
		return fmt.Errorf("failed to search playlists: %w", err)
	}

	output.PrintPlaylists(playlists)
	return nil
}

func (c *SearchController) getRequest() (req *entity.SearchRequest, err error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter search query (leave empty for all): ")
	scanner.Scan()
	query := scanner.Text()

	var country string
	fmt.Print("Enter country (leave empty for all): ")
	scanner.Scan()
	country = scanner.Text()

	var genreIDStr string
	fmt.Print("Enter genre (leave empty for all): ")
	scanner.Scan()
	genreIDStr = scanner.Text()
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse genre id: %w", err)
	}

	var limit int
	fmt.Print("Enter limit (leave empty for 10): ")
	scanner.Scan()
	limitStr := scanner.Text()
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse limit: %w", err)
		}
	} else {
		limit = 10
	}

	var offset int
	fmt.Print("Enter offset (leave empty for 0): ")
	scanner.Scan()
	offsetStr := scanner.Text()
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse offset: %w", err)
		}
	} else {
		offset = 0
	}

	searchReq := &entity.SearchRequest{
		Query: query,
		Filters: entity.Filters{
			Country: country,
			GenreID: genreID,
		},
		Limit:  limit,
		Offset: offset,
	}

	return searchReq, nil
}
