package artist_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type ArtistMetaController struct {
	artistService artist.ArtistMetaService
}

func NewArtistMetaController(artistService artist.ArtistMetaService) *ArtistMetaController {
	return &ArtistMetaController{
		artistService: artistService,
	}
}

// GetArtist godoc
// @Summary Get artist by ID
// @Description Get artist details by its ID
// @Tags artists
// @Produce json
// @Param id path string true "Artist ID"
// @Success 200 {object} entity.Artist
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/artists/{id} [get]
func (c *ArtistMetaController) GetArtist(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	artist, err := c.artistService.GetArtistMeta(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	ctx.JSON(http.StatusOK, artist)
}

func (c *ArtistMetaController) GetAllArtists(ctx *gin.Context) {
	artists, err := c.artistService.GetAllArtistMeta(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all artists"})
		return
	}

	ctx.JSON(http.StatusOK, artists)
}

// CreateArtist godoc
// @Summary Create a new artist
// @Description Create a new artist with the provided details
// @Tags artists
// @Accept json
// @Produce json
// @Param artist body entity.Artist true "Artist object"
// @Success 201 {object} entity.Artist
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/artists [post]
func (c *ArtistMetaController) CreateArtist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var artist entity.ArtistMeta
	if err := ctx.ShouldBindJSON(&artist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.artistService.CreateArtistMeta(ctx.Request.Context(), claims, &artist)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create artist"})
		}
		return
	}

	ctx.Status(http.StatusCreated)
}

// UpdateArtist godoc
// @Summary Update an artist
// @Description Update an existing artist with new details
// @Tags artists
// @Accept json
// @Produce json
// @Param id path string true "Artist ID"
// @Param artist body entity.Artist true "Updated artist object"
// @Success 200 {object} entity.Artist
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/artists/{id} [put]
func (c *ArtistMetaController) UpdateArtist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	var artist entity.ArtistMeta
	if err := ctx.ShouldBindJSON(&artist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	artist.ID = id
	err = c.artistService.UpdateArtistMeta(ctx.Request.Context(), claims, &artist)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update artist"})
		}
		return
	}

	ctx.Status(http.StatusCreated)
}

// DeleteArtist godoc
// @Summary Delete an artist
// @Description Delete an artist by its ID
// @Tags artists
// @Produce json
// @Param id path string true "Artist ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/artists/{id} [delete]
func (c *ArtistMetaController) DeleteArtist(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	err = c.artistService.DeleteArtistMeta(ctx.Request.Context(), claims, id)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete artist"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
