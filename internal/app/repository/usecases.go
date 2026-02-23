package repository

import (
	"LAB3/internal/app/ds"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/minio/minio-go/v7"
)

func (r *Repository) GetUseCases(startValue uint, endValue uint) ([]ds.UseCase, error) {
	var UseCases []ds.UseCase
	db := r.db.Model(&ds.UseCase{}).Where("is_delete = ?", false)
	if startValue > 0 {
		db = db.Where("consumption > ?", startValue)
	}
	if endValue > 0 {
		db = db.Where("consumption < ?", endValue)
	}
	err := db.Find(&UseCases).Error
	if err != nil {
		return nil, err
	}
	if len(UseCases) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return UseCases, nil
}

// GetAllUseCases возвращает все неудалённые сценарии (для кеша Redis).
func (r *Repository) GetAllUseCases() ([]ds.UseCase, error) {
	var list []ds.UseCase
	err := r.db.Model(&ds.UseCase{}).Where("is_delete = ?", false).Find(&list).Error
	return list, err
}

func (r *Repository) GetUseCase(id uint) (ds.UseCase, error) {
	UseCase := ds.UseCase{}
	err := r.db.Where("id = ?", id).First(&UseCase).Error
	if err != nil {
		return ds.UseCase{}, err
	}
	return UseCase, nil
}

func (r *Repository) AddUseCase(UseCase ds.UseCase) (ds.UseCase, error) {
	err := r.db.Create(&UseCase).Error
	if err != nil {
		return ds.UseCase{}, err
	}
	return UseCase, nil
}

func (r *Repository) DeleteUseCase(UseCaseID uint) error {
	err := r.db.Model(&ds.UseCase{}).Where("id = ?", UseCaseID).Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateUseCase(UseCaseID uint) error {
	err := r.db.Model(&ds.UseCase{}).Where("id = ?", UseCaseID).Updates(UseCaseID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateConsumption(Consumption *ds.Consumption) error {
	err := r.db.Create(Consumption).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddUseCaseToConsumption(UseCase_ID uint, Consumption_ID uint) error {
	err := r.db.Create(&ds.Usecase_consumption{
		UseCaseID:     UseCase_ID,
		ConsumptionID: Consumption_ID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddImageToUseCase(UseCaseID uint, imageURL string) error {
	return r.db.Model(&ds.UseCase{}).
		Where("id = ?", UseCaseID).
		Update("URL", imageURL).Error
}

// AddOrReplaceUseCaseImage загружает файл изображения в MinIO и сохраняет ссылку в БД.
// Принимает сам файл (multipart), а не URL. Имя объекта в MinIO — ID сценария.
func (r *Repository) AddOrReplaceUseCaseImage(ctx context.Context, useCaseID uint, header *multipart.FileHeader) (string, error) {
	if r.minio == nil {
		return "", fmt.Errorf("minio client not configured")
	}

	filename := strconv.FormatUint(uint64(useCaseID), 10)

	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	// Определяем Content-Type по содержимому (первые 512 байт)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}
	contentType := http.DetectContentType(buffer[:n])

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	_, err = r.minio.PutObject(ctx, r.minioBucketName, filename, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	imageURL := r.minioPublicUrl + "/" + r.minioBucketName + "/" + filename
	err = r.db.Model(&ds.UseCase{}).Where("id = ?", useCaseID).Update("URL", imageURL).Error
	if err != nil {
		_ = r.minio.RemoveObject(ctx, r.minioBucketName, filename, minio.RemoveObjectOptions{})
		return "", fmt.Errorf("ошибка сохранения пути к изображению: %w", err)
	}

	return imageURL, nil
}
