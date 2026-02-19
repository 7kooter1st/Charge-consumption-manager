package repository

import (
	"LAB3/internal/app/ds"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (r *Repository) GetUseCases(startValue uint, endValue uint) ([]ds.UseCase, error) {
	var UseCases []ds.UseCase
	db := r.db.Model(&ds.UseCase{}).Where("IsDelete = ?", false)
	if startValue > 0 {
		db = db.Where("consumption > ?", startValue)
	}
	if endValue > 0 {
		db = db.Where("consumption < ?", endValue)
	}
	err := r.db.Find(&UseCases).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(UseCases) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return UseCases, nil
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
	err := r.db.Model(&ds.UseCase{}).Where("id = ?", UseCaseID).Update("IsDelete", true).Error
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
// Имя файла генерируется на латинице (uc_<id>_<uuid>.<ext>).
func (r *Repository) AddOrReplaceUseCaseImage(useCaseID uint, header *multipart.FileHeader) error {
	if r.minio == nil || r.minioBucketName == "" {
		return fmt.Errorf("minio не настроен")
	}

	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	ext := getExtensionByContentType(contentType)

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	// Имя файла на латинице: uc_<id>_<uuid>.<ext>
	fileName := fmt.Sprintf("uc_%d_%s%s", useCaseID, strings.ReplaceAll(uuid.New().String(), "-", ""), ext)

	ctx := context.Background()
	_, err = r.minio.PutObject(ctx, r.minioBucketName, fileName, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("не удалось добавить объект в MinIO: %w", err)
	}

	// Ссылка для доступа к файлу в бакете
	objectURL := fmt.Sprintf("http://%s/%s/%s", r.minioEndpoint, r.minioBucketName, fileName)

	err = r.db.Model(&ds.UseCase{}).Where("id = ?", useCaseID).Update("URL", objectURL).Error
	if err != nil {
		_ = r.minio.RemoveObject(ctx, r.minioBucketName, fileName, minio.RemoveObjectOptions{})
		return fmt.Errorf("ошибка сохранения пути к изображению в БД: %w", err)
	}

	return nil
}

func getExtensionByContentType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		return ".jpg"
	case strings.HasPrefix(contentType, "image/png"):
		return ".png"
	case strings.HasPrefix(contentType, "image/gif"):
		return ".gif"
	case strings.HasPrefix(contentType, "image/webp"):
		return ".webp"
	default:
		return ".png"
	}
}

// GetExtensionFromFilename возвращает расширение из имени файла (для совместимости, если понадобится).
func GetExtensionFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return ".png"
	}
	return ext
}
