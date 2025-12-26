package artist_ctrl

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
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type ArtistAvatarController struct {
	service artist.ArtistAvatarService
}

func NewArtistAvatarController(service artist.ArtistAvatarService) *ArtistAvatarController {
	return &ArtistAvatarController{service: service}
}

func (c *ArtistAvatarController) UploadAvatar(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cover, err := c.parseArtistAvatar(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UploadCover(ctx.Request.Context(), claims, cover); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully"})
}

func (c *ArtistAvatarController) GetAvatar(ctx *gin.Context) {
	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	avatar, err := c.service.GetCover(ctx.Request.Context(), artistID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "image/jpeg") // or "image/png"
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(avatar.Data)))
	//ctx.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	ctx.Data(http.StatusOK, "image/jpeg", avatar.Data)
}

// DeleteAvatar deletes the artist's avatar
func (c *ArtistAvatarController) DeleteAvatar(ctx *gin.Context) {
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

	if err := c.service.DeleteCover(ctx.Request.Context(), claims, artistID); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete avatar"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Avatar deleted successfully"})
}

func (ArtistAvatarController) parseArtistAvatar(ctx *gin.Context) (*entity.Cover, error) {
	artistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, errors.New("invalid artist ID")
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		return nil, errors.New("avatar file not found")
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
		return nil, errors.New("failed to open avatar file")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, open); err != nil {
		return nil, errors.New("failed to copy avatar file")
	}

	return &entity.Cover{
		ObjectID: artistID,
		Data:     buf.Bytes(),
	}, nil
}
