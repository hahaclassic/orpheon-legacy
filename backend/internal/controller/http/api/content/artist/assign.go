package artist_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type ArtistAssignController struct {
	artistService artist.ArtistAssignService
	aggregator    aggregator.ContentAggregator
}

func NewArtistAssignController(artistService artist.ArtistAssignService,
	aggregator aggregator.ContentAggregator) *ArtistAssignController {
	return &ArtistAssignController{artistService: artistService, aggregator: aggregator}
}

func (c *ArtistAssignController) AssignArtistToTrack(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	trackID, err := uuid.Parse(ctx.Param("track_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	err = c.artistService.AssignArtistToTrack(ctx, claims, artistID, trackID)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Artist assigned to track successfully"})
}

func (c *ArtistAssignController) AssignArtistToAlbum(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	albumID, err := uuid.Parse(ctx.Param("album_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	err = c.artistService.AssignArtistToAlbum(ctx, claims, artistID, albumID)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Artist assigned to album successfully"})
}

func (c *ArtistAssignController) GetAlbumsByArtist(ctx *gin.Context) {
	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	albums, err := c.artistService.GetArtistAlbums(ctx, artistID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	aggregated, err := c.aggregator.GetAlbums(ctx, albums...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *ArtistAssignController) GetTracksByArtist(ctx *gin.Context) {
	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	tracks, err := c.artistService.GetArtistTracks(ctx, artistID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	aggregated, err := c.aggregator.GetTracks(ctx, tracks...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aggregated)
}

func (c *ArtistAssignController) UnassignArtistFromTrack(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	trackID, err := uuid.Parse(ctx.Param("track_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = c.artistService.UnassignArtistFromTrack(ctx, claims, artistID, trackID)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Artist removed from track successfully"})
}

func (c *ArtistAssignController) UnassignArtistFromAlbum(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	albumID, err := uuid.Parse(ctx.Param("album_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	err = c.artistService.UnassignArtistFromAlbum(ctx, claims, artistID, albumID)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Artist removed from album successfully"})
}
