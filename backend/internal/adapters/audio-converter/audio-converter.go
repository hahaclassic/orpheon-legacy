package audioconverter

import (
	"context"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AudioConverter struct{}

func New() *AudioConverter {
	return &AudioConverter{}
}

func (a AudioConverter) ChangeBitrate(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error) {
	return chunk, nil
}
