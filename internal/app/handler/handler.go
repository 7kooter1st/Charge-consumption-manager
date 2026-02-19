package handler

import (
	"LAB3/internal/app/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
}

var STATUS_CODES = map[int]string{
	403: "Forbidden",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Группа /api для всех методов согласно требованиям задания
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/:id", h.getUserData)
			users.PUT("/:id", h.changeUserData)
			users.POST("/register", h.registerUser)
		}

		usecases := api.Group("/usecases")
		{
			usecases.GET("/", h.getUseCases)
			usecases.GET("/:id", h.getUseCaseByID)
			usecases.POST("/", h.createUseCase)
			usecases.PUT("/:id", h.updateUseCase)
			usecases.DELETE("/:id", h.deleteUseCase)
			// Маршрут для обновления картинки сценария
			usecases.PUT("/:id/image", h.addImageToUseCase)
		}

		consumptions := api.Group("/consumptions")
		{
			// Основные маршруты для заявок
			consumptions.GET("/", h.getFilteredConsumptions)
			consumptions.POST("/", h.createNewConsumption)
			consumptions.GET("/draft", h.getUseCasesInConsumption)
			consumptions.GET("/:id", h.getOneConsumption)
			consumptions.DELETE("/:id", h.deleteConsumption)

			// Маршруты для изменения состояния заявки
			consumptions.PUT("/:id/formate", h.formateConsumption)
			consumptions.PUT("/:id/moderate", h.moderatorAction)

			// Маршруты для управления сценариями ВНУТРИ заявки
			consumptions.POST("/:id/usecases", h.addUseCaseToConsumption)
			consumptions.PUT("/:id/usecases/:usecase_id", h.changeUseCaseDurationInConsumption)
			consumptions.DELETE("/:id/usecases/:usecase_id", h.deleteUseCaseFromConsumption)
		}
	}

	return router
}

func (h *Handler) handleError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrNoRecords) || errors.Is(err, service.ErrUseCaseDeleted) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, service.ErrBadRequest) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Для всех остальных ошибок возвращаем текст ошибки (для отладки MinIO и т.д.)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, message string) {
	ctx.JSON(errorStatusCode, gin.H{
		"status":  strconv.Itoa(errorStatusCode) + " " + STATUS_CODES[errorStatusCode],
		"message": message,
	})
}
