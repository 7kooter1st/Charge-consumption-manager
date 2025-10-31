package service

import (
	"errors"
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
	"math"
	"time"

	"gorm.io/gorm"
)

// GetUseCasesInConsumption возвращает ID черновика потребления и количество сценариев в нём
func (s *Service) GetUseCasesInConsumption(userId uint) (uint, int64, error) {
	consumptionId, numberOfUseCases, err := s.repository.GetUseCasesInConsumption(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, ErrNoRecords
		}
		return 0, 0, err
	}
	return consumptionId, numberOfUseCases, nil
}

// GetFilteredConsumptions возвращает отфильтрованный список заявок потребления для пользователя
func (s *Service) GetFilteredConsumptions(
	userId uint,
	filter dto.ConsumptionFilter,
) ([]dto.ConsumptionsResponse, error) {
	consumptions, err := s.repository.GetFilteredConsumptions(userId, filter)
	if err != nil {
		return nil, err
	}
	if len(consumptions) == 0 {
		return nil, ErrNoRecords
	}

	var consumptionsResponse []dto.ConsumptionsResponse
	for _, consumption := range consumptions {
		moderatorLogin := ""
		if consumption.Moderator != 0 {
			moderatorLogin = consumption.Moderator.Login
		}

		consumptionsResponse = append(consumptionsResponse, dto.ConsumptionsResponse{
			ID:          consumption.ID,
			Status:      consumption.Status,
			Creator:     consumption.User.Login,
			CreatedAt:   consumption.CreatedAt,
			UpdatedAt:   consumption.UpdatedAt,
			ModeratedAt: consumption.ModeratedAt,
			Moderator:   moderatorLogin,
			TotalPower:  consumption.TotalPower,
			UserID:      consumption.UserID,
		})
	}
	return consumptionsResponse, nil
}

// GetOneConsumption возвращает одну заявку потребления со всеми сценариями использования
func (s *Service) GetOneConsumption(consumptionId uint, userId uint) (dto.OneConsumptionResponse, error) {
	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.OneConsumptionResponse{}, ErrNoRecords
		}
		return dto.OneConsumptionResponse{}, err
	}

	if consumption.UserID != userId {
		return dto.OneConsumptionResponse{}, ErrForbidden
	}

	return s.ValidateConsumptionResponse(&consumption), nil
}

// ValidateConsumptionResponse преобразует Consumption в DTO, исключая удаленные сценарии
func (s *Service) ValidateConsumptionResponse(consumption *ds.Consumption) dto.OneConsumptionResponse {
	var useCasesDTO []dto.UseCaseFromConsumptionResponse

	if consumption.Usecases != nil {
		for _, uc := range consumption.Usecases {
			if uc.UseCase != nil && !uc.UseCase.IsDelete {
				useCasesDTO = append(useCasesDTO, dto.UseCaseFromConsumptionResponse{
					ID:          uc.UseCase.ID,
					Name:        uc.UseCase.Name,
					Description: uc.UseCase.Description,
					Consumption: uc.UseCase.Consumption,
					Image:       uc.UseCase.URL,
					IsDelete:    uc.UseCase.IsDelete,
					Duration:    uc.Duration,
				})
			}
		}
	}

	response := dto.OneConsumptionResponse{
		ID:         consumption.ID,
		TotalPower: consumption.TotalPower,
		Status:     consumption.Status,
		UserID:     consumption.UserID,
		UseCases:   useCasesDTO,
	}
	return response
}

// ChangeConsumptionDuration изменяет длительность сценария в заявке потребления
func (s *Service) ChangeConsumptionDuration(
	userId uint,
	consumptionId uint,
	useCaseId uint,
	duration uint,
) (dto.OneConsumptionResponse, error) {
	if duration == 0 {
		return dto.OneConsumptionResponse{}, ErrBadRequest
	}

	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.OneConsumptionResponse{}, ErrNoRecords
		}
		return dto.OneConsumptionResponse{}, err
	}

	if consumption.UserID != userId {
		return dto.OneConsumptionResponse{}, ErrForbidden
	}

	err = s.repository.ChangeConsumptionDuration(consumptionId, useCaseId, duration)
	if err != nil {
		return dto.OneConsumptionResponse{}, err
	}

	consumption, err = s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		return dto.OneConsumptionResponse{}, err
	}

	return s.ValidateConsumptionResponse(&consumption), nil
}

// FormateConsumption формирует (завершает черновик) заявку потребления
func (s *Service) FormateConsumption(
	consumptionId uint,
	userId uint,
) (dto.OneConsumptionResponse, error) {
	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.OneConsumptionResponse{}, ErrNoRecords
		}
		return dto.OneConsumptionResponse{}, err
	}

	if consumption.UserID != userId {
		return dto.OneConsumptionResponse{}, ErrForbidden
	}

	// Проверяем, что есть хотя бы один сценарий
	if len(consumption.Usecases) == 0 {
		return dto.OneConsumptionResponse{}, ErrBadRequest
	}

	// Проверяем корректность данных
	for _, uc := range consumption.Usecases {
		if uc.Duration == 0 || math.IsNaN(float64(uc.Duration)) {
			return dto.OneConsumptionResponse{}, ErrBadRequest
		}
	}

	err = s.repository.FormateConsumption(consumptionId)
	if err != nil {
		return dto.OneConsumptionResponse{}, err
	}

	consumption, err = s.repository.GetOneConsumption(consumptionId, "сформирован")
	if err != nil {
		return dto.OneConsumptionResponse{}, err
	}

	return s.ValidateConsumptionResponse(&consumption), nil
}

// ModeratorAction выполняет действие модератора (завершение/отклонение заявки)
func (s *Service) ModeratorAction(
	consumptionId uint,
	action string,
	moderatorId uint,
) (dto.OneConsumptionResponse, error) {
	// Проверяем, что пользователь модератор
	moderator, err := s.repository.GetUser(moderatorId)
	if err != nil {
		return dto.OneConsumptionResponse{}, ErrNoRecords
	}
	if !moderator.IsModerator {
		return dto.OneConsumptionResponse{}, ErrForbidden
	}

	// Проверяем корректность действия
	if action != "завершен" && action != "отклонен" {
		return dto.OneConsumptionResponse{}, ErrBadRequest
	}

	consumption, err := s.repository.GetOneConsumption(consumptionId, "сформирован")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.OneConsumptionResponse{}, ErrNoRecords
		}
		return dto.OneConsumptionResponse{}, err
	}

	totalConsumption := 0.0
	if action == "завершен" {
		totalConsumption = float64(CalculateTotalConsumption(consumption.Usecases))
	}

	err = s.repository.ModeratorAction(consumptionId, action, totalConsumption, moderatorId)
	if err != nil {
		return dto.OneConsumptionResponse{}, err
	}

	consumption, err = s.repository.GetOneConsumption(consumptionId, action)
	if err != nil {
		return dto.OneConsumptionResponse{}, err
	}

	return s.ValidateConsumptionResponse(&consumption), nil
}

// DeleteConsumption удаляет черновик заявки потребления
func (s *Service) DeleteConsumption(consumptionId uint, userId uint) error {
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

	err = s.repository.DeleteConsumption(consumptionId)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewConsumption создает новую заявку потребления для пользователя
func (s *Service) CreateNewConsumption(userId uint) (dto.NumberOfUseCasesResponse, error) {
	consumption := ds.Consumption{
		UserID:    userId,
		Status:    "черновик",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err := s.repository.CreateConsumption(&consumption)
	if err != nil {
		return dto.NumberOfUseCasesResponse{}, err
	}

	return dto.NumberOfUseCasesResponse{
		ConsumptionId:    consumption.ID,
		NumberOfUseCases: 0,
	}, nil
}

// CalculateTotalConsumption вычисляет общее потребление энергии на основе сценариев
func CalculateTotalConsumption(useCases []ds.usecase_consumption) uint {
	total := uint(0)
	for _, uc := range useCases {
		if uc.UseCase != nil {
			total += uc.UseCase.Consumption * uc.Duration
		}
	}
	return total
}
