package audio_fs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AudioFileRepository struct {
	baseDir string
}

func NewAudioFileRepository(baseDir string) (*AudioFileRepository, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	return &AudioFileRepository{
		baseDir: baseDir,
	}, nil
}

func (r *AudioFileRepository) filePath(trackID uuid.UUID) string {
	return filepath.Join(r.baseDir, fmt.Sprintf("%s.mp3", trackID.String()))
}

func (r *AudioFileRepository) UploadAudioFile(ctx context.Context, chunk *entity.AudioChunk) error {
	path := r.filePath(chunk.TrackID)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer f.Close()

	if _, err := f.Seek(chunk.Start, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}

	if _, err := f.Write(chunk.Data); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}

func (r *AudioFileRepository) GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error) {
	path := r.filePath(chunk.TrackID)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close file", "error", err)
		}
	}()

	if _, err := f.Seek(chunk.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	size := chunk.End - chunk.Start
	if size <= 0 {
		return nil, fmt.Errorf("invalid chunk range: start=%d, end=%d", chunk.Start, chunk.End)
	}

	var data []byte
	if chunk.End == math.MaxInt64 {
		data, err = io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk: %w", err)
		}
	} else {
		data = make([]byte, size)
		limitReader := io.LimitReader(f, size)
		_, err = limitReader.Read(data)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("failed to read chunk: %w", err)
		}
	}

	return &entity.AudioChunk{
		Data:    data,
		TrackID: chunk.TrackID,
		Start:   chunk.Start,
		End:     chunk.Start + int64(len(data)),
	}, nil
}

func (r *AudioFileRepository) DeleteFile(ctx context.Context, trackID uuid.UUID) error {
	path := r.filePath(trackID)

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
