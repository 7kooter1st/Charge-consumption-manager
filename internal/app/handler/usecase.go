package handler

import (
	"net/http"
	"strconv"
	"time"

	"RIP/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetUseCases отображает список сценариев использования
func (h *Handler) GetUseCases(ctx *gin.Context) {
	var useCases []ds.UseCase
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		useCases, err = h.Repository.GetUseCases()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		useCases, err = h.Repository.GetUseCasesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// Получаем текущую заявку пользователя для проверки активности карточек
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
	}

	currentConsumption, err := h.Repository.GetCurrentConsumption(user.ID)
	hasActiveConsumption := err == nil

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":                 time.Now().Format("15:04:05"),
		"UseCases":             useCases,
		"query":                searchQuery,
		"hasActiveConsumption": hasActiveConsumption,
		"currentConsumption":   currentConsumption,
	})
}

// GetUseCaseById отображает детали сценария использования
func (h *Handler) GetUseCaseById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID",
		})
		return
	}

	useCase, err := h.Repository.GetUseCase(id)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Сценарий использования не найден",
		})
		return
	}

	ctx.HTML(http.StatusOK, "usecase.html", gin.H{
		"useCase": useCase,
	})
}
