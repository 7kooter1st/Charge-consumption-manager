package service

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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

func (s *Service) UploadUseCaseImage(ctx context.Context, useCaseId uint, file *multipart.FileHeader) (string, error) {
	// 1. Проверяем существование UseCase (можно пропустить, если Repository.AddImageToUseCase это проверяет, но лучше проверить)
	_, err := s.repository.GetUseCase(useCaseId)
	if err != nil {
		return "", err
	}

	// 2. Подготавливаем файл к загрузке
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Генерируем уникальное имя файла, чтобы не перезатереть существующие
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	bucketName := s.config.Minio.Bucket

	// 3. Проверяем, существует ли бакет, если нет — создаем
	exists, err := s.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("minio bucket check failed: %w", err)
	}
	if !exists {
		err = s.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
		// Устанавливаем политику доступа (чтобы картинки были публичными для чтения)
		policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`, bucketName)
		_ = s.minioClient.SetBucketPolicy(ctx, bucketName, policy)
	}

	// 4. Загружаем файл
	// Content-Type берем из заголовка файла или ставим "application/octet-stream"
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.minioClient.PutObject(ctx, bucketName, newFilename, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to minio: %w", err)
	}

	// 5. Формируем публичную ссылку
	// http://localhost:9000/bucket-name/filename.ext
	imageURL := fmt.Sprintf("%s/%s/%s", s.config.Minio.PublicUrl, bucketName, newFilename)

	// 6. Сохраняем ссылку в БД
	err = s.repository.AddImageToUseCase(useCaseId, imageURL)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}
