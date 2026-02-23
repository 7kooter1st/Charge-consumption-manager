package handler

// login godoc
// @Summary      Вход (логин)
// @Description  Аутентификация по логину и паролю. Возвращает access_token, refresh_token, user_id и role для фронта.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body dto.UserRegistration true "Логин и пароль"
// @Success      200  {object} dto.LoginResponse
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string "invalid login or password"
// @Router       /users/login [post]

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) login(c *gin.Context) {
	var input dto.UserRegistration
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, userID, role, err := h.Service.LoginUserWithMeta(input)
	if err != nil {
		if err == service.ErrNoRecords || err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
			return
		}
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:      userID,
		Role:        uint8(role),
	})
}

// registerUser godoc
// @Summary      Регистрация
// @Description  Создание нового пользователя (роль User по умолчанию).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body dto.UserRegistration true "Логин и пароль"
// @Success      201  {object} dto.UserDataResposne
// @Failure      400  {object} map[string]string
// @Router       /users/register [post]
func (h *Handler) registerUser(c *gin.Context) {
	var input dto.UserRegistration
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.AddNewUser(input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) getUserData(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.Service.GetUserData(uint(userId))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) changeUserData(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var input dto.ChangeUserData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.ChangeUserData(uint(userId), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// getMe godoc
// @Summary      Текущий пользователь
// @Description  Возвращает данные авторизованного пользователя (id, login, role).
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object} dto.UserDataResposne
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/users/me [get]
func (h *Handler) getMe(c *gin.Context) {
	userIdCtx, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}

	userId, ok := userIdCtx.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user id is of an invalid type"})
		return
	}

	userData, err := h.Service.GetUserData(userId)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, userData)
}
