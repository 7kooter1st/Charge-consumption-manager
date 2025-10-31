package handler

import (
	"LAB3/internal/app/Service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service.Service
}

var STATUS_CODES = map[int]string{
	403: "Forbidden",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

func NewHandler(s *Service.Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Группа маршрутов для аутентификации и регистрации пользователей
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.registerUser)
		// Здесь обычно добавляют эндпоинт для логина, который возвращает JWT
	}

	// Группа маршрутов для пользователей (защищается middleware в реальном приложении)
	users := router.Group("/users")
	{
		users.GET("/:id", h.getUserData)
		users.PUT("/:id", h.changeUserData)
	}

	// Группа маршрутов для сценариев использования (Use Cases)
	usecases := router.Group("/usecases")
	{
		usecases.GET("/", h.getUseCases)
		usecases.POST("/", h.createUseCase) // Может быть доступно только модераторам
		usecases.GET("/:id", h.getUseCaseByID)
		usecases.PUT("/:id", h.updateUseCase)    // Может быть доступно только модераторам
		usecases.DELETE("/:id", h.deleteUseCase) // Может быть доступно только модераторам
	}

	// Группа маршрутов для заявок на потребление (Consumptions)
	consumptions := router.Group("/consumptions")
	{
		// Маршруты для обычных пользователей
		consumptions.GET("/", h.getFilteredConsumptions)
		consumptions.POST("/", h.createNewConsumption)
		consumptions.GET("/draft", h.getUseCasesInConsumption)
		consumptions.GET("/:id", h.getOneConsumption)
		consumptions.DELETE("/:id", h.deleteConsumption)
		consumptions.PUT("/:id/formate", h.formateConsumption)

		// Маршруты для управления сценариями внутри заявки
		consumptions.POST("/:id/usecases", h.addUseCaseToConsumption)
		consumptions.DELETE("/:id/usecases/:usecase_id", h.deleteUseCaseFromConsumption)
		consumptions.PUT("/:id/usecases/:usecase_id", h.changeUseCaseDurationInConsumption)

		// Маршруты только для модераторов
		consumptions.PUT("/:id/moderate", h.moderatorAction)
	}

	return router
}

func (h *Handler) handleError(c *gin.Context, err error) {
	if errors.Is(err, Service.ErrNoRecords) || errors.Is(err, Service.ErrUseCaseDeleted) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, Service.ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, Service.ErrBadRequest) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Для всех остальных непредвиденных ошибок
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, message string) {
	ctx.JSON(errorStatusCode, gin.H{
		"status":  strconv.Itoa(errorStatusCode) + " " + STATUS_CODES[errorStatusCode],
		"message": message,
	})
}
