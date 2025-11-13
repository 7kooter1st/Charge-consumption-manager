package service

import (
	"LAB3/internal/app/config"
	"LAB3/internal/app/repository"
	"errors"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrNoRecords = errors.New("записи не найдены")
var ErrForbidden = errors.New("пользователь не имеет доступа к этой заявке")
var ErrBadRequest = errors.New("введены некорректные данные")
var ErrUseCaseDeleted = errors.New("этот сценарий использования удален")

type Service struct {
	repository  *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

func NewService(repository *repository.Repository) *Service {
	minioClient, err := minio.New(os.Getenv("MINIO_HOST")+":"+os.Getenv("MINIO_PORT"), &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig()
	return &Service{
		repository:  repository,
		minioClient: minioClient,
		config:      cfg,
	}
}

func (s *Service) GetConfig() *config.Config {
	return s.config
}

// Вспомогательные функции остаются без изменений
func formateDate(date time.Time, layout string) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(layout)
}

func extractFilenameFromURL(imageURL string) string {
	parsedUrl, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}
	return path.Base(parsedUrl.Path)
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
