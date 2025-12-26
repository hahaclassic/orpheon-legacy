package audio_minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/minio/minio-go/v7"
)

type AudioFileRepository struct {
	minioClient *minio.Client
	bucketName  string
}

func NewAudioFileRepository(ctx context.Context, client *minio.Client, bucketName string) (*AudioFileRepository, error) {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &AudioFileRepository{
		minioClient: client,
		bucketName:  bucketName,
	}, nil
}

func (r *AudioFileRepository) UploadAudioFile(ctx context.Context, chunk *entity.AudioChunk) error {
	objectName := chunk.TrackID.String()

	_, err := r.minioClient.PutObject(ctx, r.bucketName, objectName,
		bytes.NewReader(chunk.Data), int64(len(chunk.Data)),
		minio.PutObjectOptions{ContentType: "audio/mpeg"})
	if err != nil {
		return fmt.Errorf("failed to upload audio file: %w", err)
	}

	return nil
}

func (r *AudioFileRepository) GetAudioChunk(ctx context.Context, chunk *entity.AudioChunk) (*entity.AudioChunk, error) {
	objectName := chunk.TrackID.String()

	opts := minio.GetObjectOptions{}
	if err := opts.SetRange(chunk.Start, chunk.End-1); err != nil {
		return nil, fmt.Errorf("invalid range: %w", err)
	}

	obj, err := r.minioClient.GetObject(ctx, r.bucketName, objectName, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio chunk: %w", err)
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	if chunk.End == math.MaxInt64 {
		chunk.End = int64(len(data))
	}

	return &entity.AudioChunk{
		Data:    data,
		TrackID: chunk.TrackID,
		Start:   chunk.Start,
		End:     chunk.End,
	}, nil
}

func (r *AudioFileRepository) DeleteFile(ctx context.Context, trackID uuid.UUID) error {
	objectName := trackID.String()
	err := r.minioClient.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
