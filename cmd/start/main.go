package main

import (
	"LAB3/internal/app/config"
	"LAB3/internal/app/dsn"
	"LAB3/internal/app/handler"
	"LAB3/internal/app/repository"
	"LAB3/internal/app/service"
	"LAB3/internal/pkg"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("FATAL: error loading config: %v", err)
	}

	// 2. Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: No .env file found, using OS environment variables")
	}
	dbDSN := dsn.FromEnv()
	if dbDSN == "" {
		log.Fatal("FATAL: Database DSN is not configured.")
	}

	// MinIO: endpoint = MINIO_HOST:MINIO_PORT, ключи и бакет из .env
	minioHost := os.Getenv("MINIO_HOST")
	minioPort := os.Getenv("MINIO_PORT")
	minioEndpoint := ""
	if minioHost != "" && minioPort != "" {
		minioEndpoint = fmt.Sprintf("%s:%s", minioHost, minioPort)
	}
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "minio"
	}
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "minio124"
	}
	minioBucketName := os.Getenv("MINIO_BUCKET_NAME")
	if minioBucketName == "" {
		minioBucketName = "usecase-images"
	}

	// 3. Инициализация слоев приложения (Repository -> Service -> Handler)
	repo, err := repository.New(&repository.RepositorySettings{
		PostgresDSN:     dbDSN,
		MinioEndpoint:   minioEndpoint,
		MinioAccessKey:  minioAccessKey,
		MinioSecretKey:  minioSecretKey,
		MinioBucketName: minioBucketName,
	})
	if err != nil {
		log.Fatalf("FATAL: error initializing repository: %v", err)
	}

	appService := service.NewService(repo)

	h := handler.NewHandler(appService)

	// 4. Инициализация роутера
	// Важно: InitRoutes() настраивает все эндпоинты и возвращает готовый роутер
	router := h.InitRoutes()

	// 5. Создание и запуск приложения
	// Мы передаем все созданные зависимости в наше приложение
	application := pkg.NewApp(cfg, router, h)
	application.RunApp()
}
