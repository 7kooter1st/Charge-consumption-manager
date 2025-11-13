package service

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

// GetUseCases возвращает список сценариев использования с фильтрацией по потреблению
func (s *Service) GetUseCases(startValue uint, endValue uint) ([]ds.UseCase, error) {
	if startValue < 0 || endValue < 0 {
		return nil, ErrBadRequest
	}

	useCases, err := s.repository.GetUseCases(startValue, endValue)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecords
		}
		return nil, err
	}

	if len(useCases) == 0 {
		return nil, ErrNoRecords
	}

	return useCases, nil
}

// GetUseCase возвращает один сценарий использования по ID
func (s *Service) GetUseCase(useCaseId uint) (ds.UseCase, error) {
	useCase, err := s.repository.GetUseCase(useCaseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.UseCase{}, ErrNoRecords
		}
		return ds.UseCase{}, err
	}

	if useCase.IsDelete {
		return ds.UseCase{}, ErrUseCaseDeleted
	}

	return useCase, nil
}

// AddUseCase добавляет новый сценарий использования
func (s *Service) AddUseCase(newUseCase dto.AddUseCase) (ds.UseCase, error) {
	if newUseCase.Name == "" || newUseCase.Consumption == 0 {
		return ds.UseCase{}, ErrBadRequest
	}

	useCaseToDB := ds.UseCase{
		Name:        newUseCase.Name,
		URL:         newUseCase.URL,
		Description: newUseCase.Description,
		Consumption: newUseCase.Consumption,
		IsDelete:    false,
	}

	response, err := s.repository.AddUseCase(useCaseToDB)
	if err != nil {
		return ds.UseCase{}, err
	}

	return response, nil
}

// DeleteUseCase удаляет (мягкое удаление) сценарий использования
func (s *Service) DeleteUseCase(useCaseId uint) error {
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

	err = s.repository.DeleteUseCase(useCaseId)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUseCase обновляет данные сценария использования
func (s *Service) UpdateUseCase(useCaseId uint, updateData dto.ChangeUseCase) (ds.UseCase, error) {
	useCase, err := s.repository.GetUseCase(useCaseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.UseCase{}, ErrNoRecords
		}
		return ds.UseCase{}, err
	}

	if useCase.IsDelete {
		return ds.UseCase{}, ErrUseCaseDeleted
	}

	// Обновляем только переданные поля
	if updateData.Name != "" {
		useCase.Name = updateData.Name
	}
	if updateData.URL != "" {
		useCase.URL = updateData.URL
	}
	if updateData.Description != "" {
		useCase.Description = updateData.Description
	}
	if updateData.Consumption != 0 {
		useCase.Consumption = updateData.Consumption
	}

	err = s.repository.UpdateUseCase(useCaseId)
	if err != nil {
		return ds.UseCase{}, err
	}

	return useCase, nil
}

func (s *Service) AddImageToUseCase(useCaseId uint, imageUrl string) error {
	if imageUrl == "" {
		return ErrBadRequest
	}
	useCase, err := s.repository.GetUseCase(useCaseId)
	if err != nil {
		return err
	}
	if useCase.IsDelete {
		return ErrUseCaseDeleted
	}
	return s.repository.AddImageToUseCase(useCaseId, imageUrl)
}
