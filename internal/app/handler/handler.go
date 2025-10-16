package handler

import (
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// GET методы
	router.GET("/usecases", h.GetUseCases)              // Список сценариев использования
	router.GET("/usecase/:id", h.GetUseCaseById)        // Детали сценария использования
	router.GET("/consumption", h.GetCurrentConsumption) // Текущая заявка пользователя

	// POST методы

	router.POST("/consumption/add-duration", h.AddUseCaseToConsumption)
	router.POST("/consumption/remove-duration", h.DeleteUseCaseFromConsumption)
	// router.POST("/consumption/add", h.AddUseCaseToConsumption) // Добавить сценарий в заявку
	router.POST("/consumption/delete", h.DeleteConsumption) // Удалить заявку
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
