package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetCurrentConsumption отображает текущую заявку пользователя
func (h *Handler) GetCurrentConsumption(ctx *gin.Context) {
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка получения пользователя",
		})
		return
	}

	// Получаем все активные сценарии для отображения
	useCases, err := h.Repository.GetUseCases()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки сценариев",
		})
		return
	}

	// Получаем текущую заявку
	consumption, err := h.Repository.GetCurrentConsumption(user.ID)
	if err != nil {
		// Если заявки нет, показываем страницу с пустой заявкой
		ctx.HTML(http.StatusOK, "consumption.html", gin.H{
			"consumption": nil,
			"user":        user,
			"UseCases":    useCases,
		})
		return
	}

	// Получаем заявку с включенными сценариями
	consumptionWithUseCases, err := h.Repository.GetConsumptionWithUseCases(consumption.ID)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки заявки",
		})
		return
	}

	ctx.HTML(http.StatusOK, "consumption.html", gin.H{
		"consumption": consumptionWithUseCases,
		"user":        user,
		"UseCases":    useCases,
	})
}

// AddUseCaseToConsumption добавляет сценарий использования в заявку (обработка формы)
func (h *Handler) AddUseCaseToConsumption(ctx *gin.Context) {
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/usecases") // Редирект в случае ошибки пользователя
		return
	}

	// Получаем параметры из формы
	useCaseIDStr := ctx.PostForm("useCaseID")
	durationStr := ctx.PostForm("duration") // должно быть 15 из формы

	useCaseID, err := strconv.ParseUint(useCaseIDStr, 10, 32)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/usecases")
		return
	}

	duration, err := strconv.ParseUint(durationStr, 10, 32)
	if err != nil {
		duration = 15 // значение по умолчанию
	}

	///////////////////////////////////////////////////////
	// Получаем или создаем текущую заявку
	consumption, err := h.Repository.GetCurrentConsumption(user.ID)
	if err != nil {
		// Создаем новую заявку, если не найдена
		consumption, err = h.Repository.CreateConsumption(user.ID)
		if err != nil {
			logrus.Error(err)
			ctx.Redirect(http.StatusSeeOther, "/usecases")
			return
		}
	}
	/////////////////////////////////////////////////////////

	// Проверяем, существует ли сценарий, и обновляем или добавляем его
	exists, err := h.Repository.CheckIfUseCaseInConsumption(consumption.ID, uint(useCaseID))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/usecases")
		return
	}

	if exists {
		err = h.Repository.IncreaseUseCaseDuration(consumption.ID, uint(useCaseID), uint(duration))
	} else {
		err = h.Repository.AddUseCaseToConsumption(consumption.ID, uint(useCaseID), uint(duration))
	}

	if err != nil {
		logrus.Error(err)
		// Не паникуем, просто редирект, чтобы не показывать ошибку
	}

	// Пересчет и обновление общей мощности
	totalPower, _ := h.Repository.CalculateTotalPower(consumption.ID)
	h.Repository.UpdateConsumptionTotalPower(consumption.ID, totalPower)

	// *** Ключевое изменение: Редирект обратно на главную страницу ***
	ctx.Redirect(http.StatusSeeOther, "/usecases")
}

// DeleteUseCaseFromConsumption уменьшает время сценария или удаляет его (обработка формы)
func (h *Handler) DeleteUseCaseFromConsumption(ctx *gin.Context) {
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/usecases") // Редирект в случае ошибки пользователя
		return
	}

	useCaseIDStr := ctx.PostForm("useCaseID")
	useCaseID, err := strconv.ParseUint(useCaseIDStr, 10, 32)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/usecases")
		return
	}

	consumption, err := h.Repository.GetCurrentConsumption(user.ID)
	if err != nil {
		// Если заявки нет, ничего не делаем
		ctx.Redirect(http.StatusSeeOther, "/usecases")
		return
	}

	// Уменьшаем время на 15 минут
	err = h.Repository.DecreaseUseCaseDuration(consumption.ID, uint(useCaseID), 15)
	if err != nil {
		logrus.Error(err)
		// Не паникуем, просто редирект, чтобы не показывать ошибку
	}

	// Пересчет и обновление общей мощности
	totalPower, _ := h.Repository.CalculateTotalPower(consumption.ID)
	h.Repository.UpdateConsumptionTotalPower(consumption.ID, totalPower)

	// *** Ключевое изменение: Редирект обратно на главную страницу ***
	ctx.Redirect(http.StatusSeeOther, "/usecases")
}

// DeleteConsumption логически удаляет заявку
func (h *Handler) DeleteConsumption(ctx *gin.Context) {
	// Получаем пользователя по умолчанию
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	// Получаем текущую заявку
	consumption, err := h.Repository.GetCurrentConsumption(user.ID)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	// Удаляем заявку (меняем статус на "удален")
	err = h.Repository.DeleteConsumption(consumption.ID)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/consumption")
}

// IncreaseUseCaseDuration увеличивает время использования сценария
func (h *Handler) IncreaseUseCaseDuration(ctx *gin.Context) {
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	useCaseIDStr := ctx.PostForm("useCaseID")
	incrementStr := ctx.PostForm("increment")

	useCaseID, err := strconv.ParseUint(useCaseIDStr, 10, 32)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	increment, err := strconv.ParseUint(incrementStr, 10, 32)
	if err != nil {
		increment = 15
	}

	consumption, err := h.Repository.GetCurrentConsumption(user.ID)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	err = h.Repository.IncreaseUseCaseDuration(consumption.ID, uint(useCaseID), uint(increment))
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/consumption")
}

// RemoveUseCaseFromConsumption удаляет сценарий из заявки
func (h *Handler) RemoveUseCaseFromConsumption(ctx *gin.Context) {
	user, err := h.Repository.GetOrCreateDefaultUser()
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	useCaseIDStr := ctx.PostForm("useCaseID")
	useCaseID, err := strconv.ParseUint(useCaseIDStr, 10, 32)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	consumption, err := h.Repository.GetCurrentConsumption(user.ID)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/consumption")
		return
	}

	err = h.Repository.RemoveUseCaseFromConsumption(consumption.ID, uint(useCaseID))
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/consumption")
}
