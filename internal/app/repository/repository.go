package repository

import (
	"context"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// RepositorySettings — параметры подключения к БД и MinIO
type RepositorySettings struct {
	PostgresDSN      string
	MinioEndpoint    string
	MinioAccessKey   string
	MinioSecretKey   string
	MinioBucketName  string
}

type Repository struct {
	db                *gorm.DB
	minio             *minio.Client
	minioBucketName   string
	minioEndpoint     string
}

func New(settings *RepositorySettings) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(settings.PostgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// MinIO опционален: если эндпоинт не задан, работаем без MinIO
	var minioClient *minio.Client
	if settings.MinioEndpoint != "" && settings.MinioBucketName != "" {
		minioClient, err = minio.New(settings.MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
			Secure: false,
		})
		if err != nil {
			return nil, err
		}
		// Создаём бакет при старте, если его ещё нет
		ctx := context.Background()
		exists, err := minioClient.BucketExists(ctx, settings.MinioBucketName)
		if err != nil {
			return nil, err
		}
		if !exists {
			err = minioClient.MakeBucket(ctx, settings.MinioBucketName, minio.MakeBucketOptions{})
			if err != nil {
				return nil, err
			}
			log.Printf("MinIO: бакет %q создан", settings.MinioBucketName)
		}
	}

	return &Repository{
		db:              db,
		minio:           minioClient,
		minioBucketName: settings.MinioBucketName,
		minioEndpoint:   settings.MinioEndpoint,
	}, nil
}
