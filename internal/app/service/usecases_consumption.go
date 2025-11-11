package service

import (
	"LAB3/internal/app/ds"
	"errors"
	"math"

	"gorm.io/gorm"
)

// AddUseCaseToConsumption добавляет сценарий использования в заявку потребления
func (s *Service) AddUseCaseToConsumption(
	userId uint,
	consumptionId uint,
	useCaseId uint,
) error {
	// Проверяем существование заявки и её принадлежность пользователю
	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoRecords
		}
		return err
	}

	if consumption.UserID != userId {
		return ErrForbidden
	}

	// Проверяем существование и доступность сценария
	useCase, err := s.repository.GetUseCase(useCaseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoRecords
		}
		return err
	}

	if useCase.IsDelete {
		return ErrUseCaseDeleted
	}

	// Добавляем сценарий в заявку
	err = s.repository.AddUseCaseToConsumption(useCaseId, consumptionId)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUseCaseFromConsumption удаляет сценарий из заявки потребления
func (s *Service) DeleteUseCaseFromConsumption(
	userId uint,
	consumptionId uint,
	useCaseId uint,
) error {
	// Проверяем принадлежность заявки пользователю
	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoRecords
		}
		return err
	}

	if consumption.UserID != userId {
		return ErrForbidden
	}

	// Проверяем наличие сценария в заявке
	exists := false
	for _, uc := range consumption.Usecases {
		if uc.UseCaseID == useCaseId {
			exists = true
			break
		}
	}

	if !exists {
		return ErrNoRecords
	}

	// Удаляем сценарий из заявки
	err = s.repository.DeleteUseCaseFromRequest(consumptionId, useCaseId)
	if err != nil {
		return err
	}

	return nil
}

// ChangeUseCaseDurationInConsumption изменяет длительность сценария в заявке
// func (s *Service) ChangeUseCaseDurationInConsumption(
// 	userId uint,
// 	consumptionId uint,
// 	useCaseId uint,
// 	duration uint,
// ) (ds.Usecase_consumption, error) { // <-- Имя типа исправлено на UsecaseConsumption
// 	if duration == 0 || math.IsNaN(float64(duration)) {
// 		return ds.Usecase_consumption{}, ErrBadRequest
// 	}

// 	// Проверяем принадлежность заявки пользователю
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ds.Usecase_consumption{}, ErrNoRecords
// 		}
// 		return ds.Usecase_consumption{}, err
// 	}

// 	if consumption.UserID != userId {
// 		return ds.Usecase_consumption{}, ErrForbidden
// 	}

// 	// Обновляем длительность
// 	// (Предполагается, что ChangeUseCaseConsumption тоже исправлена и работает с UsecaseConsumption)
// 	err = s.repository.ChangeUseCaseConsumption(consumptionId, useCaseId, duration)
// 	if err != nil {
// 		return ds.Usecase_consumption{}, err
// 	}

// 	// Получаем обновленный record. Теперь эта функция вернет правильный тип.
// 	response, err := s.repository.GetUseCaseFromConsumption(consumptionId, useCaseId)
// 	if err != nil {
// 		// Ошибка может быть gorm.ErrRecordNotFound, если запись не найдена
// 		return ds.Usecase_consumption{}, ErrNoRecords
// 	}

// 	// Теперь типы совпадают, и ошибки не будет
// 	return response, nil
// }

// ChangeUseCaseDurationInConsumption изменяет длительность сценария в заявке
func (s *Service) ChangeUseCaseDurationInConsumption(
	userId uint,
	consumptionId uint,
	useCaseId uint,
	duration uint,
) (ds.Usecase_consumption, error) { // <-- Исправлено на UsecaseConsumption
	if duration == 0 || math.IsNaN(float64(duration)) {
		return ds.Usecase_consumption{}, ErrBadRequest // <-- Исправлено
	}

	// ... (остальной код проверки)
	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Usecase_consumption{}, ErrNoRecords // <-- Исправлено
		}
		return ds.Usecase_consumption{}, err // <-- Исправлено
	}
	if consumption.UserID != userId {
		return ds.Usecase_consumption{}, ErrForbidden // <-- Исправлено
	}

	// Обновляем длительность
	err = s.repository.ChangeUseCaseConsumption(consumptionId, useCaseId, duration)
	if err != nil {
		return ds.Usecase_consumption{}, err // <-- Исправлено
	}

	// Получаем обновленный record. Теперь эта функция вернет правильный тип.
	response, err := s.repository.GetUseCaseFromConsumption(consumptionId, useCaseId)
	if err != nil {
		return ds.Usecase_consumption{}, ErrNoRecords // <-- Исправлено
	}

	// Теперь типы совпадают, и ошибки не будет
	return response, nil
}
