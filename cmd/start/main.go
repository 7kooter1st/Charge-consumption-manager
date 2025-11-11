package main

import (
	"LAB3/internal/app/dsn"
	"LAB3/internal/app/handler"
	"LAB3/internal/app/repository"
	"LAB3/internal/app/service"
	"log"

	"github.com/joho/godotenv"
)

// Функция для загрузки переменных окружения из файла .env
// Это удобно для локальной разработки
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS environment variables")
	}
}

func main() {
	// Загружаем переменные окружения
	loadEnv()

	// 1. Инициализация репозитория
	// << ИЗМЕНЕНИЕ: Получаем DSN из переменных окружения с помощью вашей функции
	dbDSN := dsn.FromEnv()

	// << ВАЖНО: Добавляем проверку, что DSN был успешно сформирован
	if dbDSN == "" {
		log.Fatal("Database DSN is not configured. Please set DB_HOST, DB_PORT, etc. environment variables or create a .env file.")
	}

	repo, err := repository.New(dbDSN)
	if err != nil {
		log.Fatalf("failed to initialize repository: %s", err.Error())
	}

	log.Println("Successfully connected to the database")

	// 2. Инициализация сервиса
	appService := service.NewService(repo)

	// 3. Инициализация обработчика
	h := handler.NewHandler(appService)

	// 4. Запуск сервера
	router := h.InitRoutes()
	log.Println("Starting server on :8000")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %s", err.Error())
	}
}
