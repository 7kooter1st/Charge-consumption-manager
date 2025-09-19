package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"RIP/iternal/app/handler"    // исправлено
	"RIP/iternal/app/repository" // исправлено
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
    r.LoadHTMLGlob("templates/*")
    // Serve static assets (CSS, images) under /static
    r.Static("/static", "resources")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/hello", handler.GetUsecases)
	r.GET("/usecase/:id", handler.GetUsecase)
    r.GET("/message", handler.GetMessage)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}