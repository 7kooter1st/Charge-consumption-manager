package handler

import (
	dto "LAB3/internal/app/DTO"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUseCase godoc
// @Summary      Create a new use case
// @Description  Creates a new use case. Requires moderator permissions.
// @Tags         Use Cases (Moderator)
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth // <--- Указываем, что этот эндпоинт защищен
// @Param        useCase body dto.AddUseCase true "Use Case data"
// @Success      201  {object} ds.UseCase
// @Failure      400  {object} map[string]string "Bad Request"
// @Failure      401  {object} map[string]string "Unauthorized"
// @Failure      403  {object} map[string]string "Forbidden"
// @Failure      500  {object} map[string]string "Internal Server Error"
// @Router       /api/usecases [post] // <--- Обновляем путь

func (h *Handler) getUseCases(c *gin.Context) {
	startValue, _ := strconv.ParseUint(c.Query("start_value"), 10, 32)
	endValue, _ := strconv.ParseUint(c.Query("end_value"), 10, 32)

	useCases, err := h.Service.GetUseCases(c.Request.Context(), uint(startValue), uint(endValue))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, useCases)
}

func (h *Handler) createUseCase(c *gin.Context) {
	var input dto.AddUseCase
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	useCase, err := h.Service.AddUseCase(input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, useCase)
}

func (h *Handler) getUseCaseByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	useCase, err := h.Service.GetUseCase(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, useCase)
}

func (h *Handler) updateUseCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input dto.ChangeUseCase
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUseCase, err := h.Service.UpdateUseCase(uint(id), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedUseCase)
}

func (h *Handler) deleteUseCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.Service.DeleteUseCase(uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Use case marked as deleted"})
}

func (h *Handler) addImageToUseCase(c *gin.Context) {
	// 1. Парсим ID из URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID format"})
		return
	}

	// 2. Получаем файл из формы (ключ "file" или "image")
	// На фронтенде FormData должна иметь поле с этим именем!
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required (key 'file')"})
		return
	}

	// 3. Вызываем сервис
	url, err := h.Service.UploadUseCaseImage(c.Request.Context(), uint(id), file)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"url":     url,
	})
}
