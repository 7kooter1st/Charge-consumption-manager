package service

import (
	"encoding/json"
	"time"

	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
	"context"
	"errors"
	"mime/multipart"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const useCasesCacheKey = "usecases:list"
const useCasesCacheTTL = 5 * time.Minute

// GetUseCases возвращает список сценариев использования с фильтрацией по потреблению.
// Список услуг кешируется в Redis (доп. задание).
func (s *Service) GetUseCases(ctx context.Context, startValue uint, endValue uint) ([]ds.UseCase, error) {
	if startValue < 0 || endValue < 0 {
		return nil, ErrBadRequest
	}

	var all []ds.UseCase
	cached, err := s.redisClient.Get(ctx, useCasesCacheKey).Bytes()
	if err == nil {
		if err := json.Unmarshal(cached, &all); err == nil {
			return filterUseCasesByConsumption(all, startValue, endValue)
		}
	}
	if err != nil && err != redis.Nil {
		// при ошибке Redis идём в БД без кеша
	}

	all, err = s.repository.GetAllUseCases()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecords
		}
		return nil, err
	}
	_ = s.redisClient.Set(ctx, useCasesCacheKey, mustMarshal(all), useCasesCacheTTL).Err()

	return filterUseCasesByConsumption(all, startValue, endValue)
}

func filterUseCasesByConsumption(list []ds.UseCase, startValue, endValue uint) ([]ds.UseCase, error) {
	var out []ds.UseCase
	for _, u := range list {
		if startValue > 0 && u.Consumption <= startValue {
			continue
		}
		if endValue > 0 && u.Consumption >= endValue {
			continue
		}
		out = append(out, u)
	}
	if len(out) == 0 {
		return nil, ErrNoRecords
	}
	return out, nil
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (s *Service) invalidateUseCasesCache(ctx context.Context) {
	_ = s.redisClient.Del(ctx, useCasesCacheKey).Err()
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
	s.invalidateUseCasesCache(context.Background())
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
	s.invalidateUseCasesCache(context.Background())
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
	s.invalidateUseCasesCache(context.Background())
	return useCase, nil
}

// func (s *Service) AddImageToUseCase(useCaseId uint, imageUrl string) error {
// 	if imageUrl == "" {
// 		return ErrBadRequest
// 	}
// 	useCase, err := s.repository.GetUseCase(useCaseId)
// 	if err != nil {
// 		return err
// 	}
// 	if useCase.IsDelete {
// 		return ErrUseCaseDeleted
// 	}
// 	return s.repository.AddImageToUseCase(useCaseId, imageUrl)
// }

// UploadUseCaseImage принимает файл изображения (multipart), загружает его в MinIO и сохраняет ссылку в БД.
// Эндпоинт ожидает форму с полем "file" (сама картинка), а не JSON с URL.
func (s *Service) UploadUseCaseImage(ctx context.Context, useCaseId uint, file *multipart.FileHeader) (string, error) {
	_, err := s.repository.GetUseCase(useCaseId)
	if err != nil {
		return "", err
	}
	url, err := s.repository.AddOrReplaceUseCaseImage(ctx, useCaseId, file)
	if err != nil {
		return "", err
	}
	s.invalidateUseCasesCache(ctx)
	return url, nil
}
