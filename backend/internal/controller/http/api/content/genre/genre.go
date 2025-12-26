package genre_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type GenreController struct {
	genreService   genre.GenreService
	authMiddleware gin.HandlerFunc
}

func NewGenreController(genreService genre.GenreService, authMiddleware gin.HandlerFunc) *GenreController {
	return &GenreController{
		genreService:   genreService,
		authMiddleware: authMiddleware,
	}
}

func (c *GenreController) RegisterRoutes(router *gin.RouterGroup) {
	genres := router.Group("/genres")
	{
		genres.GET("/:id", c.GetGenre)
		genres.GET("", c.GetAllGenres)

		protected := genres.Group("")
		protected.Use(c.authMiddleware)
		{
			protected.POST("", c.CreateGenre)
			protected.PUT("/:id", c.UpdateGenre)
			protected.DELETE("/:id", c.DeleteGenre)
		}
	}
}

// GetGenre godoc
// @Summary Get genre by ID
// @Description Get genre details by its ID
// @Tags genres
// @Produce json
// @Param id path string true "Genre ID"
// @Success 200 {object} entity.Genre
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/genres/{id} [get]
func (c *GenreController) GetGenre(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	genre, err := c.genreService.GetGenreByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	ctx.JSON(http.StatusOK, genre)
}

// GetAllGenres godoc
// @Summary Get all genres
// @Description Get all genres
// @Tags genres
// @Produce json
// @Success 200 {array} entity.Genre
// @Failure 500 {object} gin.H
// @Router /api/v1/genres [get]
func (c *GenreController) GetAllGenres(ctx *gin.Context) {
	genres, err := c.genreService.GetAllGenres(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all genres"})
		return
	}

	ctx.JSON(http.StatusOK, genres)
}

// CreateGenre godoc
// @Summary Create a new genre
// @Description Create a new genre with the provided details
// @Tags genres
// @Accept json
// @Produce json
// @Param genre body entity.Genre true "Genre object"
// @Success 201 {object} entity.Genre
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/genres [post]
func (c *GenreController) CreateGenre(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var genre entity.Genre
	if err := ctx.ShouldBindJSON(&genre); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.genreService.CreateGenre(ctx.Request.Context(), claims, &genre)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to create genre"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
		return
	}

	ctx.Status(http.StatusCreated)
}

// UpdateGenre godoc
// @Summary Update a genre
// @Description Update an existing genre with new details
// @Tags genres
// @Accept json
// @Produce json
// @Param id path string true "Genre ID"
// @Param genre body entity.Genre true "Updated genre object"
// @Success 200 {object} entity.Genre
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/genres/{id} [put]
func (c *GenreController) UpdateGenre(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	var genre entity.Genre
	if err := ctx.ShouldBindJSON(&genre); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre.ID = id
	err = c.genreService.UpdateGenre(ctx.Request.Context(), claims, &genre)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to update genre"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update genre"})
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteGenre godoc
// @Summary Delete a genre
// @Description Delete a genre by its ID
// @Tags genres
// @Produce json
// @Param id path string true "Genre ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/genres/{id} [delete]
func (c *GenreController) DeleteGenre(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	err = c.genreService.DeleteGenre(ctx.Request.Context(), claims, id)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to delete genre"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete genre"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
