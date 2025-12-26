package playlist_cover_minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/minio/minio-go/v7"
)

type PlaylistCoverRepository struct {
	minio      *minio.Client
	bucketName string
}

func NewPlaylistCoverRepository(ctx context.Context, client *minio.Client, bucketName string) (*PlaylistCoverRepository, error) {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &PlaylistCoverRepository{
		minio:      client,
		bucketName: bucketName,
	}, err
}

func (r *PlaylistCoverRepository) SaveCover(ctx context.Context, cover *entity.Cover) error {
	objectName := fmt.Sprintf("playlist_covers/%s", cover.ObjectID)

	reader := bytes.NewReader(cover.Data)

	_, err := r.minio.PutObject(ctx, r.bucketName, objectName, reader, int64(len(cover.Data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return fmt.Errorf("failed to upload cover to MinIO: %w", err)
	}

	return nil
}

func (r *PlaylistCoverRepository) GetCover(ctx context.Context, objectID uuid.UUID) (*entity.Cover, error) {
	objectName := fmt.Sprintf("playlist_covers/%s", objectID)

	// Check if object exists first
	_, err := r.minio.StatObject(ctx, r.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, fmt.Errorf("%w: %v", commonerr.ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to check cover existence in MinIO: %w", err)
	}

	object, err := r.minio.GetObject(ctx, r.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cover from MinIO: %w", err)
	}
	defer func() {
		err := object.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read cover data: %w", err)
	}

	return &entity.Cover{
		ObjectID: objectID,
		Data:     data,
	}, nil
}

func (r *PlaylistCoverRepository) DeleteCover(ctx context.Context, objectID uuid.UUID) error {
	objectName := fmt.Sprintf("playlist_covers/%s", objectID)

	err := r.minio.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete cover from MinIO: %w", err)
	}

	return nil
}
