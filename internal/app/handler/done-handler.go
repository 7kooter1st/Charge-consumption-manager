package handler

// import (
// 	"LAB3/internal/app/dto"
// 	"LAB3/internal/app/Service"
// 	"errors"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// Handler является оберткой для сервисного слоя

// New создает новый экземпляр Handler
// func New(s *Service.Service) *Handler {
// 	return &Handler{
// 		Service: s,
// 	}
// }

// handleError централизованно обрабатывает ошибки от сервисного слоя
// func (h *Handler) handleError(c *gin.Context, err error) {
// 	if errors.Is(err, Service.ErrNoRecords) || errors.Is(err, Service.ErrUseCaseDeleted) {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if errors.Is(err, Service.ErrForbidden) {
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if errors.Is(err, Service.ErrBadRequest) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// Для всех остальных непредвиденных ошибок
// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// }

// // getUserIdFromQuery - Вспомогательная функция.
// func getUserIdFromQuery(c *gin.Context) (uint, error) {
// 	userIdStr := c.Query("user_id")
// 	if userIdStr == "" {
// 		return 0, errors.New("user_id query parameter is required")
// 	}
// 	userId, err := strconv.ParseUint(userIdStr, 10, 32)
// 	if err != nil {
// 		return 0, errors.New("invalid user_id format")
// 	}
// 	return uint(userId), nil
// }

// --- Обработчики для Пользователей (Users) ---

// func (h *Handler) registerUser(c *gin.Context) {
// 	var input dto.UserRegistration
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user, err := h.Service.AddNewUser(input)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, user)
// }

// func (h *Handler) getUserData(c *gin.Context) {
// 	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
// 		return
// 	}

// 	user, err := h.Service.GetUserData(uint(userId))
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// func (h *Handler) changeUserData(c *gin.Context) {
// 	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
// 		return
// 	}

// 	var input dto.ChangeUserData
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user, err := h.Service.ChangeUserData(uint(userId), input)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// --- Обработчики для Сценариев (UseCases) ---

// func (h *Handler) getUseCases(c *gin.Context) {
// 	startValue, _ := strconv.ParseUint(c.Query("start_value"), 10, 32)
// 	endValue, _ := strconv.ParseUint(c.Query("end_value"), 10, 32)

// 	useCases, err := h.Service.GetUseCases(uint(startValue), uint(endValue))
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, useCases)
// }

// func (h *Handler) createUseCase(c *gin.Context) {
// 	var input dto.AddUseCase
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	useCase, err := h.Service.AddUseCase(input)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, useCase)
// }

// func (h *Handler) getUseCaseByID(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	useCase, err := h.Service.GetUseCase(uint(id))
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, useCase)
// }

// func (h *Handler) updateUseCase(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	var input dto.ChangeUseCase
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedUseCase, err := h.Service.UpdateUseCase(uint(id), input)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, updatedUseCase)
// }

// func (h *Handler) deleteUseCase(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	if err := h.Service.DeleteUseCase(uint(id)); err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Use case marked as deleted"})
// }

// --- Обработчики для Заявок (Consumptions) ---

// func (h *Handler) getUseCasesInConsumption(c *gin.Context) {
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	consumptionId, count, err := h.Service.GetUseCasesInConsumption(userId)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, dto.NumberOfUseCasesResponse{
// 		ConsumptionId:    consumptionId,
// 		NumberOfUseCases: count,
// 	})
// }

// func (h *Handler) getFilteredConsumptions(c *gin.Context) {
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var filter dto.ConsumptionFilter
// 	filter.Status = c.Query("status")
// 	if startDateStr := c.Query("start_date"); startDateStr != "" {
// 		// Пример парсинга даты, формат можно изменить
// 		filter.Start_date, _ = time.Parse("2006-01-02", startDateStr)
// 	}
// 	if endDateStr := c.Query("end_date"); endDateStr != "" {
// 		filter.End_date, _ = time.Parse("2006-01-02", endDateStr)
// 	}

// 	consumptions, err := h.Service.GetFilteredConsumptions(userId, filter)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, consumptions)
// }

// func (h *Handler) createNewConsumption(c *gin.Context) {
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response, err := h.Service.CreateNewConsumption(userId)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, response)
// }

// func (h *Handler) getOneConsumption(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	consumption, err := h.Service.GetOneConsumption(uint(consumptionId), userId)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, consumption)
// }

// func (h *Handler) deleteConsumption(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.Service.DeleteConsumption(uint(consumptionId), userId); err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Consumption draft deleted"})
// }

// func (h *Handler) formateConsumption(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response, err := h.Service.FormateConsumption(uint(consumptionId), userId)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func (h *Handler) addUseCaseToConsumption(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var payload struct {
// 		UseCaseID uint `json:"use_case_id" binding:"required"`
// 	}
// 	if err := c.ShouldBindJSON(&payload); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = h.Service.AddUseCaseToConsumption(userId, uint(consumptionId), payload.UseCaseID)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Use case added to consumption"})
// }

// func (h *Handler) deleteUseCaseFromConsumption(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}
// 	useCaseId, err := strconv.ParseUint(c.Param("usecase_id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID format"})
// 		return
// 	}
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = h.Service.DeleteUseCaseFromConsumption(userId, uint(consumptionId), uint(useCaseId))
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Use case removed from consumption"})
// }

// func (h *Handler) changeUseCaseDurationInConsumption(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}
// 	useCaseId, err := strconv.ParseUint(c.Param("usecase_id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID format"})
// 		return
// 	}
// 	userId, err := getUserIdFromQuery(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var input dto.ChangeUseCaseDurationRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response, err := h.Service.ChangeUseCaseDurationInConsumption(userId, uint(consumptionId), uint(useCaseId), input.Duration)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func (h *Handler) moderatorAction(c *gin.Context) {
// 	consumptionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID format"})
// 		return
// 	}

// 	var input dto.ModeratorAction
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response, err := h.Service.ModeratorAction(uint(consumptionId), input.Action, input.ModeratorID)
// 	if err != nil {
// 		h.handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
