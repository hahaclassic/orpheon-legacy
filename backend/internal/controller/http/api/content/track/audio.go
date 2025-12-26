package track_ctrl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
)

type TrackAudioController struct {
	service usecase.AudioFileService
}

func NewTrackAudioController(service usecase.AudioFileService) *TrackAudioController {
	return &TrackAudioController{service: service}
}

func (c *TrackAudioController) GetAudioChunk(ctx *gin.Context) {
	trackID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	start, end, err := c.parseRangeHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chunk := &entity.AudioChunk{
		TrackID: trackID,
		Start:   start,
		End:     end,
	}

	result, err := c.service.GetAudioChunk(ctx.Request.Context(), chunk)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "audio/mpeg")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(result.Data)))
	ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", result.Start, result.End-1, result.End))
	ctx.Header("Accept-Ranges", "bytes")

	ctx.Data(http.StatusPartialContent, "audio/mpeg", result.Data)
}

func (c *TrackAudioController) parseRangeHeader(ctx *gin.Context) (int64, int64, error) {
	rangeHeader := ctx.GetHeader("Range")
	if rangeHeader == "" {
		return 0, 0, errors.New("Range header is required")
	}

	// Parse Range header (format: "bytes=start-end")
	rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeStr, "-")
	if len(parts) < 1 || len(parts) > 2 {
		return 0, 0, errors.New("Invalid range format")
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("Invalid start range")
	}

	end := int64(0)
	if len(parts) == 2 && parts[1] != "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, errors.New("Invalid end range")
		}
	} else {
		end = math.MaxInt64
	}

	return start, end, nil
}

func (c *TrackAudioController) UploadAudioFile(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chunk, err := c.parseAudioChunk(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UploadAudioFile(ctx.Request.Context(), claims, chunk); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Audio file uploaded successfully"})
}

func (c *TrackAudioController) DeleteAudioFile(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trackID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	if err := c.service.DeleteAudioFile(ctx.Request.Context(), claims, trackID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Audio file deleted successfully"})
}

func (c *TrackAudioController) parseAudioChunk(ctx *gin.Context) (*entity.AudioChunk, error) {
	trackID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, errors.New("invalid track ID")
	}

	file, err := ctx.FormFile("audio")
	if err != nil {
		return nil, errors.New("audio file not found")
	}

	if file.Size > 30*1024*1024 { // 30MB limit
		return nil, errors.New("file size exceeds 30MB limit")
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "audio/mpeg" && contentType != "audio/mp3" {
		return nil, errors.New("only MP3 files are allowed")
	}

	open, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open audio file")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, open); err != nil {
		return nil, errors.New("failed to copy audio file")
	}

	return &entity.AudioChunk{
		TrackID: trackID,
		Start:   0,
		End:     int64(len(buf.Bytes())),
		Data:    buf.Bytes(),
	}, nil
}
