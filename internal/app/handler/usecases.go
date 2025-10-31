package handler

import (
	"LAB3/internal/app/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUseCases(c *gin.Context) {
	startValue, _ := strconv.ParseUint(c.Query("start_value"), 10, 32)
	endValue, _ := strconv.ParseUint(c.Query("end_value"), 10, 32)

	useCases, err := h.service.GetUseCases(uint(startValue), uint(endValue))
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

	useCase, err := h.service.AddUseCase(input)
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

	useCase, err := h.service.GetUseCase(uint(id))
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

	updatedUseCase, err := h.service.UpdateUseCase(uint(id), input)
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

	if err := h.service.DeleteUseCase(uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Use case marked as deleted"})
}
