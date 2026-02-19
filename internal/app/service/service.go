package service

import (
	"LAB3/internal/app/repository"
	"errors"

	//"mime/multipart"
	"time"
)

var ErrNoRecords = errors.New("записи не найдены")
var ErrForbidden = errors.New("пользователь не имеет доступа к этой заявке")
var ErrBadRequest = errors.New("введены некорректные данные")
var ErrUseCaseDeleted = errors.New("этот сценарий использования удален")

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func formateDate(date time.Time, layout string) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(layout)
}

// ////////////////////
// func CalculateTotalConsumption(useCases []ds.usecase_consumption) uint {
// 	total := uint(0)
// 	for _, uc := range useCases {
// 		// Потребление = потребление за минуту × длительность в минутах
// 		total += uc.UseCase.Consumption * uc.Duration
// 	}
// 	return total
// }

//////////////////////

// Дополнительные методы для бизнес-логики
func (s *Service) ValidateConsumptionData(userId uint, consumptionId uint) error {
	// Проверка, принадлежит ли заявка пользователю
	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		return ErrNoRecords
	}
	if consumption.UserID != userId {
		return ErrForbidden
	}
	return nil
}
