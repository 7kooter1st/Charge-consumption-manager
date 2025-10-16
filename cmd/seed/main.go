package main

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Создаем пользователя по умолчанию
	user := ds.Users{
		Username: "default_user",
		Email:    "default@example.com",
		IsActive: true,
	}
	db.FirstOrCreate(&user, ds.Users{Username: "default_user"})

	// Создаем тестовые сценарии использования
	useCases := []ds.UseCase{
		{
			Title:            "Игры",
			Description:      "Использование смартфона для игр",
			ImageURL:         "http://127.0.0.1:9000/lab1/premium_photo-1682125220008-006e43a6896a.jpeg",
			PowerConsumption: 100,
			IsActive:         true,
		},
		{
			Title:            "Просмотр видео",
			Description:      "Просмотр видео на YouTube и других платформах",
			ImageURL:         "http://127.0.0.1:9000/lab1/video.jpeg",
			PowerConsumption: 20,
			IsActive:         true,
		},
		{
			Title:            "Просмотр reels",
			Description:      "Просмотр коротких видео в Instagram",
			ImageURL:         "http://127.0.0.1:9000/lab1/instagram.jpeg",
			PowerConsumption: 22,
			IsActive:         true,
		},
		{
			Title:            "Мессенджеры",
			Description:      "Использование WhatsApp, Telegram и других мессенджеров",
			ImageURL:         "http://127.0.0.1:9000/lab1/messanger.jpeg",
			PowerConsumption: 20,
			IsActive:         true,
		},
		{
			Title:            "Навигатор",
			Description:      "Использование GPS навигации",
			ImageURL:         "http://127.0.0.1:9000/lab1/maps.jpeg",
			PowerConsumption: 1448,
			IsActive:         true,
		},
		{
			Title:            "Звонки",
			Description:      "Голосовые звонки",
			ImageURL:         "http://127.19.0.1:9000/lab1/phone.jpeg",
			PowerConsumption: 427,
			IsActive:         true,
		},
		{
			Title:            "Видеозвонки",
			Description:      "Видеозвонки через Zoom, Skype и другие приложения",
			ImageURL:         "http://127.19.0.1:9000/lab1/zoom.jpeg",
			PowerConsumption: 100,
			IsActive:         true,
		},
		{
			Title:            "Музыка",
			Description:      "Прослушивание музыки",
			ImageURL:         "http://127.19.0.1:9000/lab1/music.jpeg",
			PowerConsumption: 15,
			IsActive:         true,
		},
	}

	for _, uc := range useCases {
		db.FirstOrCreate(&uc, ds.UseCase{Title: uc.Title})
	}

	println("Тестовые данные добавлены успешно!")
}
