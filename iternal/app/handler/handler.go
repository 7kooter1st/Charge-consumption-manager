package handler

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "RIP/iternal/app/repository"
  "net/http"
  "time"
  "strconv"
)

type Handler struct {
  Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  return &Handler{
    Repository: r,
  }
}

func (h *Handler) GetUsecase(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	usecase, err := h.Repository.GetUsecase(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "usecase.html", gin.H{
		"usecase": usecase,
	})
}

func (h *Handler) GetUsecases(ctx *gin.Context) {
	var usecases []repository.Usecase
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		usecases, err = h.Repository.GetUsecases()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		usecases, err = h.Repository.GetUsecasesByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":   time.Now().Format("15:04:05"),
		"usecases": usecases,
		"query":  searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
		})
	}

func (h *Handler) GetMessage(ctx *gin.Context) {
	// ctx.HTML(http.StatusOK, "message.html", gin.H{})
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	usecase, err := h.Repository.GetUsecase(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "message.html", gin.H{
		"usecase": usecase,
	})
}