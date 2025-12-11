package handler

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/role"
	"LAB3/internal/app/service"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) refreshToken(c *gin.Context) {
	var input dto.RefreshRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.Service.RefreshTokens(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) logout(c *gin.Context) {
	// 1. Получаем заголовок Authorization
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	// 2. Извлекаем токен (убираем "Bearer ")
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}
	tokenStr := headerParts[1]

	// 3. Добавляем токен в черный список через сервис
	// (Метод AddToBlacklist мы добавили в service.go в предыдущих шагах)
	err := h.Service.AddToBlacklist(c.Request.Context(), tokenStr)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

// func (h *Handler) InitRoutes() *gin.Engine {
// 	router := gin.Default()

// 	router.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // В проде лучше указать конкретный домен фронта
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})

// 	router.POST("/users/register", h.registerUser)
// 	router.POST("/users/login", h.login)

// 	api := router.Group("/api", h.userIdentity)
// 	{
// 		api.GET("/users/me", h.getMe) // Нужно будет создать обработчик getMe
// 		consumptions := api.Group("/consumptions")
// 		{
// 			consumptions.GET("/", h.getFilteredConsumptions)
// 			consumptions.POST("/", h.createNewConsumption)
// 			consumptions.GET("/draft", h.getUseCasesInConsumption)
// 			consumptions.GET("/:id", h.getOneConsumption)
// 			consumptions.DELETE("/:id", h.deleteConsumption)
// 			consumptions.PUT("/:id/formate", h.formateConsumption)

// 			// Управление сценариями внутри заявки
// 			consumptions.POST("/:id/usecases", h.addUseCaseToConsumption)
// 			consumptions.PUT("/:id/usecases/:usecase_id", h.changeUseCaseDurationInConsumption)
// 			consumptions.DELETE("/:id/usecases/:usecase_id", h.deleteUseCaseFromConsumption)
// 		}

// 		// Чтение сценариев доступно всем залогиненным пользователям
// 		api.GET("/usecases", h.getUseCases)
// 		api.GET("/usecases/:id", h.getUseCaseByID)
// 		// api.GET("/me", h.getMe)
// 		moderator := api.Group("/", h.requireRole(role.Moderator))
// 		{
// 			// Модератор может модерировать заявки
// 			moderator.PUT("/consumptions/:id/moderate", h.moderatorAction)

// 			// Модератор может управлять сценариями (создавать, изменять, удалять)
// 			usecases := moderator.Group("/usecases")
// 			{
// 				usecases.POST("/", h.createUseCase)
// 				usecases.PUT("/:id", h.updateUseCase)
// 				usecases.DELETE("/:id", h.deleteUseCase)
// 				usecases.PUT("/:id/image", h.addImageToUseCase)
// 			}
// 		}
// 	}
// 	return router
// }

func (h *Handler) InitRoutes(router *gin.Engine) {

	// Настройка CORS (ОЧЕНЬ ВАЖНО ДЛЯ ФРОНТЕНДА)
	// Добавь middleware для CORS, если его еще нет, иначе фронт не сможет слать запросы
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // В проде лучше указать конкретный домен фронта
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- ПУБЛИЧНАЯ ЗОНА (Доступна без токена) ---

	// Auth
	router.POST("/users/register", h.registerUser)
	router.POST("/users/login", h.login)
	router.POST("/users/refresh", h.refreshToken) // Если реализовано

	// Usecases (Чтение доступно всем!)
	publicUseCases := router.Group("/usecases")
	{
		publicUseCases.GET("/", h.getUseCases)       // Получить список
		publicUseCases.GET("/:id", h.getUseCaseByID) // Получить подробности
	}

	// Swagger (если подключен)
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- ЗАЩИЩЕННАЯ ЗОНА (Требуется логин) ---
	api := router.Group("/api", h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/me", h.getMe)
			users.POST("/logout", h.logout)
		}

		consumptions := api.Group("/consumptions")
		{
			// ... твои маршруты заявок ...
			consumptions.GET("/", h.getFilteredConsumptions)
			consumptions.POST("/", h.createNewConsumption)
			// и так далее
		}

		// --- ЗОНА МОДЕРАТОРА ---
		moderator := api.Group("/", h.requireRole(role.Moderator))
		{
			// Управление сценариями (только модератор может менять их)
			modUseCases := moderator.Group("/usecases")
			{
				modUseCases.POST("/", h.createUseCase)
				modUseCases.PUT("/:id", h.updateUseCase)
				modUseCases.DELETE("/:id", h.deleteUseCase)

				// Загрузка картинки
				modUseCases.PUT("/:id/image", h.addImageToUseCase)
			}

			moderator.PUT("/consumptions/:id/moderate", h.moderatorAction)
		}
	}
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
	// Для всех остальных непредвиденных ошибок
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, message string) {
	ctx.JSON(errorStatusCode, gin.H{
		"status":  strconv.Itoa(errorStatusCode) + " " + STATUS_CODES[errorStatusCode],
		"message": message,
	})
}
