package minio

import (
	"fmt"

	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig = config.MinIOConfig

// type MinIOConfig struct {
// 	Endpoint  string `env:"MINIO_ENDPOINT" env-required:"true"`
// 	AccessKey string `env:"MINIO_ACCESS_KEY" env-required:"true"`
// 	SecretKey string `env:"MINIO_SECRET_KEY" env-required:"true"`
// 	Bucket    string `env:"MINIO_BUCKET" env-required:"true"`
// 	Secure    bool   `env:"MINIO_SECURE" env-default:"false"`
// }

func NewMinioClient(cfg MinIOConfig) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}
	return client, nil
}
