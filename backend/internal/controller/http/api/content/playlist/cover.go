package playlist_ctrl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type PlaylistCoverController struct {
	service playlist.PlaylistCoverService
}

func NewPlaylistCoverController(service playlist.PlaylistCoverService) *PlaylistCoverController {
	return &PlaylistCoverController{
		service: service,
	}
}

func (c *PlaylistCoverController) UploadCover(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cover, err := c.parsePlaylistCover(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UploadCover(ctx.Request.Context(), claims, cover); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cover uploaded successfully"})
}

func (c *PlaylistCoverController) GetCover(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)

	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	cover, err := c.service.GetCover(ctx.Request.Context(), claims, playlistID)
	if err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to get cover"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Header("Content-Type", "image/jpeg")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(cover.Data)))
	//ctx.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	ctx.Data(http.StatusOK, "image/jpeg", cover.Data)
}

func (c *PlaylistCoverController) DeleteCover(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	if err := c.service.DeleteCover(ctx.Request.Context(), claims, playlistID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cover deleted successfully"})
}

func (PlaylistCoverController) parsePlaylistCover(ctx *gin.Context) (*entity.Cover, error) {
	playlistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, errors.New("invalid playlist ID")
	}

	file, err := ctx.FormFile("cover")
	if err != nil {
		return nil, errors.New("cover file not found")
	}

	if file.Size > 10*1024*1024 { // 10MB limit
		return nil, errors.New("file size exceeds 10MB limit")
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, errors.New("only JPEG and PNG files are allowed")
	}

	open, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open cover file")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, open); err != nil {
		return nil, errors.New("failed to copy cover file")
	}

	return &entity.Cover{
		ObjectID: playlistID,
		Data:     buf.Bytes(),
	}, nil
}
