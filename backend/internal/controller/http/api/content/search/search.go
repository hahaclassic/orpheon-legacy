package search_ctrl

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
)

type SearchController struct {
	searchService      search.SearchService
	contentAggregator  aggregator.ContentAggregator
	playlistAggregator playlist.PlaylistAggregator
	authMiddleware     gin.HandlerFunc
}

func NewSearchController(searchService search.SearchService,
	contentAggregator aggregator.ContentAggregator,
	playlistAggregator playlist.PlaylistAggregator,
	authMiddleware gin.HandlerFunc) *SearchController {
	return &SearchController{
		searchService:      searchService,
		contentAggregator:  contentAggregator,
		playlistAggregator: playlistAggregator,
		authMiddleware:     authMiddleware,
	}
}

func (c *SearchController) RegisterRoutes(router *gin.RouterGroup) {
	search := router.Group("/search")
	search.Use(c.authMiddleware) // optional auth middleware
	{
		search.GET("", c.Search)
	}
}

func (c *SearchController) parseSearchRequest(ctx *gin.Context) (*entity.SearchRequest, error) {
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid limit parameter")
	}
	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid offset parameter")
	}

	query := ctx.Query("query")
	country := ctx.Query("country")
	genreID, err := uuid.Parse(ctx.Query("genre_id"))
	if err != nil {
		genreID = uuid.Nil
	}

	searchRequest := &entity.SearchRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
		Filters: entity.Filters{
			GenreID: genreID,
			Country: country,
		},
	}
	return searchRequest, nil
}

func (c *SearchController) Search(ctx *gin.Context) {
	searchRequest, err := c.parseSearchRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := map[string]func(ctx *gin.Context, searchRequest *entity.SearchRequest){
		"track":    c.searchTracks,
		"album":    c.searchAlbums,
		"artist":   c.searchArtists,
		"playlist": c.searchPlaylists,
	}

	searchFunc, ok := search[ctx.Query("type")]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	searchFunc(ctx, searchRequest)
}

func (c *SearchController) searchTracks(ctx *gin.Context, searchRequest *entity.SearchRequest) {
	result, err := c.searchService.SearchTracks(ctx.Request.Context(), searchRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}
	aggregated, err := c.contentAggregator.GetTracks(ctx.Request.Context(), result...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate search results"})
		return
	}
	ctx.JSON(http.StatusOK, aggregated)
}

func (c *SearchController) searchAlbums(ctx *gin.Context, searchRequest *entity.SearchRequest) {
	result, err := c.searchService.SearchAlbums(ctx.Request.Context(), searchRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}
	aggregated, err := c.contentAggregator.GetAlbums(ctx.Request.Context(), result...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate search results"})
		return
	}
	ctx.JSON(http.StatusOK, aggregated)
}

func (c *SearchController) searchArtists(ctx *gin.Context, searchRequest *entity.SearchRequest) {
	result, err := c.searchService.SearchArtists(ctx.Request.Context(), searchRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *SearchController) searchPlaylists(ctx *gin.Context, searchRequest *entity.SearchRequest) {
	claims := ctxclaims.GetClaims(ctx)

	result, err := c.searchService.SearchPlaylists(ctx.Request.Context(), claims, searchRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	aggregated, err := c.playlistAggregator.GetPlaylists(ctx.Request.Context(), claims, result...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate search results"})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}
