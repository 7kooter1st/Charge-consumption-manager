package handler

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body dto.UserRegistration true "User credentials"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
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

	accessToken, refreshToken, err := h.Service.LoginUser(input)
	if err != nil {
		if err == service.ErrNoRecords || err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
			return
		}
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accesstoken":  accessToken,
		"refreshtoken": refreshToken,
	})
}

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
