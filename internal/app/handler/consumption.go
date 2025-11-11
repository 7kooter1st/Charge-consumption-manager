package handler

import (
	dto "LAB3/internal/app/DTO"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getUserIdFromQuery(c *gin.Context) (uint, error) {
	// 1. Получаем значение параметра 'user_id' из URL (например, /consumptions?user_id=123)
	userIdStr := c.Query("user_id")

	// 2. Проверяем, был ли параметр вообще передан. Если нет, возвращаем ошибку.
	if userIdStr == "" {
		return 0, errors.New("query parameter 'user_id' is required")
	}

	// 3. Конвертируем строковое значение в число (uint64).
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		// Если конвертация не удалась (например, передали не число), возвращаем ошибку.
		return 0, errors.New("invalid 'user_id' format, must be a number")
	}

	// 4. Возвращаем результат в формате uint и nil в качестве ошибки.
	return uint(userId), nil
}

func (h *Handler) getUseCasesInConsumption(c *gin.Context) {
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumptionId, count, err := h.Service.GetUseCasesInConsumption(userId)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NumberOfUseCasesResponse{
		ConsumptionId:    consumptionId,
		NumberOfUseCases: count,
	})
}

func (h *Handler) getFilteredConsumptions(c *gin.Context) {
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filter dto.ConsumptionFilter
	filter.Status = c.Query("status")
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		// Пример парсинга даты, формат можно изменить
		filter.Start_date, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		filter.End_date, _ = time.Parse("2006-01-02", endDateStr)
	}

	consumptions, err := h.Service.GetFilteredConsumptions(userId, filter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, consumptions)
}

func (h *Handler) createNewConsumption(c *gin.Context) {
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.CreateNewConsumption(userId)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) getOneConsumption(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumption, err := h.Service.GetOneConsumption(uint(consumptionId), userId)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, consumption)
}

func (h *Handler) deleteConsumption(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.DeleteConsumption(uint(consumptionId), userId); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consumption draft deleted"})
}

func (h *Handler) formateConsumption(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.FormateConsumption(uint(consumptionId), userId)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) addUseCaseToConsumption(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payload struct {
		UseCaseID uint `json:"use_case_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.AddUseCaseToConsumption(userId, uint(consumptionId), payload.UseCaseID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Use case added to consumption"})
}

func (h *Handler) deleteUseCaseFromConsumption(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}
	useCaseId, err := strconv.ParseUint(c.Param("usecase_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID format"})
		return
	}
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.DeleteUseCaseFromConsumption(userId, uint(consumptionId), uint(useCaseId))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Use case removed from consumption"})
}

func (h *Handler) changeUseCaseDurationInConsumption(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}
	useCaseId, err := strconv.ParseUint(c.Param("usecase_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID format"})
		return
	}
	userId, err := getUserIdFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input dto.ChangeUseCaseDurationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.ChangeUseCaseDurationInConsumption(userId, uint(consumptionId), uint(useCaseId), input.Duration)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) moderatorAction(c *gin.Context) {
	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
		return
	}

	var input dto.ModeratorAction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.ModeratorAction(uint(consumptionId), input.Action, input.ModeratorID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
