package handler

import (
	dto "LAB3/internal/app/DTO"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUseCases(c *gin.Context) {
	startValue, _ := strconv.ParseUint(c.Query("start_value"), 10, 32)
	endValue, _ := strconv.ParseUint(c.Query("end_value"), 10, 32)

	useCases, err := h.Service.GetUseCases(uint(startValue), uint(endValue))
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID format"})
		return
	}

	// Получаем файл из multipart/form-data (поле "image" или "file")
	fileHeader, err := c.FormFile("image")
	if err != nil {
		fileHeader, err = c.FormFile("file")
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется файл изображения (form-data, поле 'image' или 'file')"})
		return
	}

	err = h.Service.AddImageToUseCase(uint(id), fileHeader)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image added/updated successfully"})
}
