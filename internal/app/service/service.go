package service

import (
	"LAB3/internal/app/config"
	"LAB3/internal/app/ds"
	"LAB3/internal/app/repository"
	"context"
	"errors"
	"net/url"
	"path"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/minio/minio-go/v7"
)

var ErrNoRecords = errors.New("записи не найдены")
var ErrForbidden = errors.New("пользователь не имеет доступа к этой заявке")
var ErrBadRequest = errors.New("введены некорректные данные")
var ErrUseCaseDeleted = errors.New("этот сценарий использования удален")

const jwtBlacklistPrefix = "jwt_blacklist:"

type Service struct {
	repository  *repository.Repository
	config      *config.Config
	minioClient *minio.Client
	redisClient *redis.Client
}

func (s *Service) AddToBlacklist(ctx context.Context, tokenStr string) error {
	// Парсим токен без проверки подписи, чтобы просто получить из него время истечения (ExpiresAt)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &ds.JWTClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Ключ в Redis будет, например, "jwt_blacklist:eyJhbGciOi..."
	key := jwtBlacklistPrefix + tokenStr

	// Время жизни ключа в Redis (TTL) равно времени, оставшемуся до истечения токена.
	// Это нужно, чтобы не засорять Redis просроченными токенами.
	ttl := time.Until(claims.ExpiresAt.Time)

	// Если токен уже истек, нет смысла добавлять его в блэклист.
	if ttl <= 0 {
		return nil // Не является ошибкой
	}

	// Добавляем ключ в Redis с указанным временем жизни.
	return s.redisClient.Set(ctx, key, "revoked", ttl).Err()
}

func (s *Service) IsInBlacklist(ctx context.Context, tokenStr string) (bool, error) {
	key := jwtBlacklistPrefix + tokenStr

	// Пытаемся получить ключ из Redis
	err := s.redisClient.Get(ctx, key).Err()

	if err == redis.Nil {
		// redis.Nil - это специальная ошибка, означающая "ключ не найден".
		// Для нас это не ошибка, а нормальная ситуация: токен НЕ в блэклисте.
		return false, nil
	}
	if err != nil {
		// Любая другая ошибка (например, Redis недоступен) является проблемой.
		return false, err
	}

	// Если мы дошли сюда, значит err == nil, ключ был найден. Токен в блэклисте.
	return true, nil
}

// func NewService(repository *repository.Repository) *Service {
// 	minioClient, err := minio.New(os.Getenv("MINIO_HOST")+":"+os.Getenv("MINIO_PORT"), &minio.Options{
// 		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
// 		Secure: false,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	cfg, err := config.NewConfig()
// 	return &Service{
// 		repository:  repository,
// 		minioClient: minioClient,
// 		config:      cfg,
// 	}
// }

func New(
	repo *repository.Repository,
	cfg *config.Config,
	minioClient *minio.Client,
	redisClient *redis.Client,
) *Service {
	return &Service{
		repository:  repo,
		config:      cfg,
		minioClient: minioClient,
		redisClient: redisClient,
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
