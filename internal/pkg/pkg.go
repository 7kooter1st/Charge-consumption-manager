package pkg

import (
	"LAB3/internal/app/config"
	"LAB3/internal/app/handler"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// App - это основная структура нашего приложения.
// Она хранит в себе все ключевые зависимости.
type App struct {
	config  *config.Config
	router  *gin.Engine
	handler *handler.Handler
}

// NewApp - это конструктор для создания нового экземпляра App.
func NewApp(cfg *config.Config, router *gin.Engine, h *handler.Handler) *App {
	return &App{
		config:  cfg,
		router:  router,
		handler: h,
	}
}

// RunApp - это метод, который запускает наше приложение.
// Вся логика по определению адреса и запуску Gin-роутера находится здесь.
func (a *App) RunApp() {
	// Получаем хост и порт из конфигурации, которая хранится в структуре App
	logrus.Info("Server start up")
	portStr := strconv.Itoa(a.config.ServicePort)
	address := a.config.ServiceHost + ":" + portStr

	log.Printf("INFO: Starting server on address: %s", address)

	// Запускаем роутер
	if err := a.router.Run(address); err != nil {
		log.Fatalf("FATAL: failed to run server: %s", err.Error())
	}
}
