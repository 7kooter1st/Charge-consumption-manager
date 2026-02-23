package repository

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// RepositorySettings параметры подключения к БД и MinIO.
type RepositorySettings struct {
	PostgresDSN      string
	MinioClient      *minio.Client
	MinioBucketName  string
	MinioPublicUrl   string
}

type Repository struct {
	db               *gorm.DB
	minio            *minio.Client
	minioBucketName  string
	minioPublicUrl   string
}

func New(settings *RepositorySettings) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(settings.PostgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:              db,
		minio:           settings.MinioClient,
		minioBucketName:  settings.MinioBucketName,
		minioPublicUrl:   settings.MinioPublicUrl,
	}, nil
}