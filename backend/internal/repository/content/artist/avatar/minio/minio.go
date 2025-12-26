package avatar_minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/minio/minio-go/v7"
)

type ArtistAvatarRepository struct {
	client     *minio.Client
	bucketName string
}

func NewArtistAvatarRepository(ctx context.Context, client *minio.Client, bucketName string) (*ArtistAvatarRepository, error) {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &ArtistAvatarRepository{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (r *ArtistAvatarRepository) SaveCover(ctx context.Context, cover *entity.Cover) error {
	reader := bytes.NewReader(cover.Data)

	_, err := r.client.PutObject(ctx, r.bucketName, cover.ObjectID.String(), reader, int64(len(cover.Data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return fmt.Errorf("put artist avatar to minio: %w", err)
	}

	return nil
}

func (r *ArtistAvatarRepository) GetCover(ctx context.Context, artistID uuid.UUID) (*entity.Cover, error) {
	obj, err := r.client.GetObject(ctx, r.bucketName, artistID.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get artist avatar from minio: %w", err)
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			slog.Error("err", "object close error", err)
		}
	}()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read artist avatar from stream: %w", err)
	}

	return &entity.Cover{
		ObjectID: artistID,
		Data:     data,
	}, nil
}

func (r *ArtistAvatarRepository) DeleteCover(ctx context.Context, artistID uuid.UUID) error {
	err := r.client.RemoveObject(ctx, r.bucketName, artistID.String(), minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete artist avatar from minio: %w", err)
	}

	return nil
}
