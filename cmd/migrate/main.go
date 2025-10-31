package main

import (
	"LAB3/internal/app/ds"
	"LAB3/internal/app/dsn"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// init() вызывается перед main()
func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS environment variables")
	}
}

func main() {
	log.Println("Starting database migration...")

	// Получаем DSN для подключения к БД из переменных окружения
	dbDSN := dsn.FromEnv()
	if dbDSN == "" {
		log.Fatal("Database DSN is not configured. Please set DB_HOST, etc. environment variables.")
	}

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful.")

	// Запускаем AutoMigrate
	// GORM автоматически создаст таблицы, отсутствующие столбцы,
	// индексы и внешние ключи.
	// Порядок важен для GORM, чтобы правильно определить внешние ключи.
	err = db.AutoMigrate(
		&ds.User{},
		&ds.UseCase{},
		&ds.Consumption{},
		&ds.Usecase_consumption{}, // Наша исправленная промежуточная таблица
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("✅ Database migration completed successfully!")
}
