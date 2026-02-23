package handler

import (
	_ "embed"
	"errors"
	"net/http"
	"strconv"
	"strings"

	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/role"
	"LAB3/internal/app/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed openapi.json
var openAPISpec []byte

//go:embed swagger.html
var swaggerHTML []byte

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

// refreshToken godoc
// @Summary      Обновление токенов
// @Description  По refresh_token выдаёт новую пару access_token и refresh_token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body dto.RefreshRequest true "refresh_token"
// @Success      200  {object} dto.TokenResponse
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Router       /users/refresh [post]
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

func (h *Handler) InitRoutes(router *gin.Engine) {

	config := cors.DefaultConfig()
	// Разрешаем запросы с любых доменов (для разработки удобно, на проде лучше указать конкретный)
	config.AllowAllOrigins = true
	// Разрешаем методы
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"}
	// Разрешаем заголовки (важно добавить Authorization для токенов)
	config.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	// Разрешаем фронтенду видеть определенные заголовки ответа
	config.ExposeHeaders = []string{"Content-Length"}
	// Разрешаем передачу куки/креденшелов (если нужно)
	config.AllowCredentials = true

	// Применяем middleware
	router.Use(cors.New(config))

	// --- ПУБЛИЧНАЯ ЗОНА (Доступна без токена) ---

	// Auth
	router.POST("/users/register", h.registerUser)
	router.POST("/users/login", h.login)
	router.POST("/users/refresh", h.refreshToken) // Если реализовано

	// Usecases (чтение доступно всем)
	publicUseCases := router.Group("/useCases")
	{
		publicUseCases.GET("/", h.getUseCases)
		publicUseCases.GET("/:id", h.getUseCaseByID)
	}
	// Те же эндпоинты под /api/usecases/ — без авторизации
	router.GET("/api/usecases/", h.getUseCases)
	router.GET("/api/usecases/:id", h.getUseCaseByID)

	// Спецификация OpenAPI и Swagger UI
	router.GET("/swagger/doc.json", h.serveOpenAPISpec)
	router.GET("/swagger/index.html", h.serveSwaggerUI)
	router.GET("/swagger/", h.serveSwaggerUI)

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
			consumptions.GET("/draft", h.getUseCasesInConsumption)
			consumptions.GET("/:id", h.getOneConsumption)
			consumptions.DELETE("/:id", h.deleteConsumption)
			consumptions.PUT("/:id/formate", h.formateConsumption)

			// Управление сценариями внутри заявки
			consumptions.POST("/:id/usecases", h.addUseCaseToConsumption)
			consumptions.PUT("/:id/usecases/:usecase_id", h.changeUseCaseDurationInConsumption)
			consumptions.DELETE("/:id/usecases/:usecase_id", h.deleteUseCaseFromConsumption)
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

func (h *Handler) serveOpenAPISpec(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", openAPISpec)
}

func (h *Handler) serveSwaggerUI(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Data(http.StatusOK, "text/html; charset=utf-8", swaggerHTML)
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
