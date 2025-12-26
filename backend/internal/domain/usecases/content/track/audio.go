package track

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetAudioChunk   = errors.New("failed to get audio chunk")
	ErrUploadAudioFile = errors.New("failed to upload audio file")
	ErrDeleteAudioFile = errors.New("failed to delete audio file")
)

type AudioFileService interface {
	GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error)
	// Admin
	UploadAudioFile(ctx context.Context, claims *entity.Claims, chunk *entity.AudioChunk) error
	DeleteAudioFile(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error
}
