package main

import (
	"LAB3/internal/app/config"
	"LAB3/internal/app/dsn"
	"LAB3/internal/app/handler"
	"LAB3/internal/app/repository"
	"LAB3/internal/app/service"
	"LAB3/internal/pkg"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes https http
// @BasePath /

// @contact.name API Support
// @contact.url ...
// @contact.email ...

// @license.name AS IS (NO WARRANTY)

// @host localhost:8000
// @BasePath /

// --- ДОБАВЬТЕ ЭТО ОПИСАНИЕ БЕЗОПАСНОСТИ ---
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Введите "Bearer" пробел и затем ваш токен. Пример: "Bearer eyJhbGciOiJI..."

// func main() {
// 	// 1. Инициализация конфигурации
// 	cfg, err := config.NewConfig()
// 	if err != nil {
// 		log.Fatalf("FATAL: error loading config: %v", err)
// 	}

// 	// 2. Инициализация подключения к БД
// 	// godotenv.Load() нужен для DSN, так как он читает .env
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("INFO: No .env file found, using OS environment variables for DB")
// 	}
// 	dbDSN := dsn.FromEnv()
// 	if dbDSN == "" {
// 		log.Fatal("FATAL: Database DSN is not configured.")
// 	}

// 	// 3. Инициализация слоев приложения (Repository -> Service -> Handler)
// 	repo, err := repository.New(dbDSN)
// 	if err != nil {
// 		log.Fatalf("FATAL: error initializing repository: %v", err)
// 	}

// 	appService := service.NewService(repo)

// 	h := handler.NewHandler(appService)

// 	// 4. Инициализация роутера
// 	// Важно: InitRoutes() настраивает все эндпоинты и возвращает готовый роутер
// 	router := h.InitRoutes()

// 	// 5. Создание и запуск приложения
// 	// Мы передаем все созданные зависимости в наше приложение
// 	application := pkg.NewApp(cfg, router, h)
// 	application.RunApp()
// }

func main() {
	// 1. Инициализация конфигурации ИЗ ОДНОГО МЕСТА
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("FATAL: error loading config: %v", err)
	}

	// 2. Инициализация подключения к БД
	// DSN теперь тоже можно собирать из cfg, но ваш способ тоже рабочий
	dbDSN := dsn.FromEnv() // Эта функция читает DB_* из .env
	if dbDSN == "" {
		log.Fatal("FATAL: Database DSN is not configured.")
	}
	repo, err := repository.New(dbDSN)
	if err != nil {
		log.Fatalf("FATAL: error initializing repository: %v", err)
	}

	// 3. Инициализация ВНЕШНИХ КЛИЕНТОВ (MinIO, Redis)
	minioClient, err := minio.New(cfg.Minio.Host+":"+cfg.Minio.Port, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.User, cfg.Minio.Pass, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("FATAL: error initializing minio client: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("FATAL: failed to connect to redis: %v", err)
	}

	// 4. Инициализация СЕРВИСНОГО СЛОЯ с передачей ВСЕХ зависимостей
	appService := service.New(repo, cfg, minioClient, redisClient)

	// 5. Инициализация ОБРАБОТЧИКА
	h := handler.NewHandler(appService) // Используем New вместо NewHandler

	// 6. Инициализация роутера
	router := gin.Default()
	h.InitRoutes() // Предполагается, что InitRoutes настраивает переданный роутер

	// 7. Создание и запуск приложения
	application := pkg.NewApp(cfg, router, h)
	application.RunApp()
}
