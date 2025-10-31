
// // consumption.go
// package service

// import (
// 	dto "lab/internal/app/DTO"
// 	"lab/internal/app/ds"
// 	"errors"
// 	"gorm.io/gorm"
// 	"math"
// 	"time"
// )

// // GetUseCasesInConsumption возвращает ID черновика потребления и количество сценариев в нём
// func (s *Service) GetUseCasesInConsumption(userId uint) (uint, int64, error) {
// 	consumptionId, numberOfUseCases, err := s.repository.GetUseCasesInConsumption(userId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return 0, 0, ErrNoRecords
// 		}
// 		return 0, 0, err
// 	}
// 	return consumptionId, numberOfUseCases, nil
// }

// // GetFilteredConsumptions возвращает отфильтрованный список заявок потребления для пользователя
// func (s *Service) GetFilteredConsumptions(
// 	userId uint,
// 	filter dto.ConsumptionFilter,
// ) ([]dto.ConsumptionsResponse, error) {
// 	consumptions, err := s.repository.GetFilteredConsumptions(userId, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(consumptions) == 0 {
// 		return nil, ErrNoRecords
// 	}

// 	var consumptionsResponse []dto.ConsumptionsResponse
// 	for _, consumption := range consumptions {
// 		moderatorLogin := ""
// 		if consumption.Moderator != nil {
// 			moderatorLogin = consumption.Moderator.Login
// 		}

// 		consumptionsResponse = append(consumptionsResponse, dto.ConsumptionsResponse{
// 			ID:          consumption.ID,
// 			Status:      consumption.Status,
// 			Creator:     consumption.User.Login,
// 			CreatedAt:   consumption.CreatedAt,
// 			UpdatedAt:   consumption.UpdatedAt,
// 			ModeratedAt: consumption.ModeratedAt,
// 			Moderator:   moderatorLogin,
// 			TotalPower:  consumption.TotalPower,
// 			UserID:      consumption.UserID,
// 		})
// 	}
// 	return consumptionsResponse, nil
// }

// // GetOneConsumption возвращает одну заявку потребления со всеми сценариями использования
// func (s *Service) GetOneConsumption(consumptionId uint, userId uint) (dto.OneConsumptionResponse, error) {
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return dto.OneConsumptionResponse{}, ErrNoRecords
// 		}
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	if consumption.UserID != userId {
// 		return dto.OneConsumptionResponse{}, ErrForbidden
// 	}

// 	return s.ValidateConsumptionResponse(&consumption), nil
// }

// // ValidateConsumptionResponse преобразует Consumption в DTO, исключая удаленные сценарии
// func (s *Service) ValidateConsumptionResponse(consumption *ds.Consumption) dto.OneConsumptionResponse {
// 	var useCasesDTO []dto.UseCaseFromConsumptionResponse

// 	if consumption.Usecases != nil {
// 		for _, uc := range consumption.Usecases {
// 			if uc.UseCase != nil && !uc.UseCase.IsDelete {
// 				useCasesDTO = append(useCasesDTO, dto.UseCaseFromConsumptionResponse{
// 					ID:          uc.UseCase.ID,
// 					Name:        uc.UseCase.Name,
// 					Description: uc.UseCase.Description,
// 					Consumption: uc.UseCase.Consumption,
// 					Image:       uc.UseCase.URL,
// 					IsDelete:    uc.UseCase.IsDelete,
// 					Duration:    uc.Duration,
// 				})
// 			}
// 		}
// 	}

// 	response := dto.OneConsumptionResponse{
// 		ID:         consumption.ID,
// 		TotalPower: consumption.TotalPower,
// 		Status:     consumption.Status,
// 		UserID:     consumption.UserID,
// 		UseCases:   useCasesDTO,
// 	}
// 	return response
// }

// // ChangeConsumptionDuration изменяет длительность сценария в заявке потребления
// func (s *Service) ChangeConsumptionDuration(
// 	userId uint,
// 	consumptionId uint,
// 	useCaseId uint,
// 	duration uint,
// ) (dto.OneConsumptionResponse, error) {
// 	if duration == 0 {
// 		return dto.OneConsumptionResponse{}, ErrBadRequest
// 	}

// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return dto.OneConsumptionResponse{}, ErrNoRecords
// 		}
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	if consumption.UserID != userId {
// 		return dto.OneConsumptionResponse{}, ErrForbidden
// 	}

// 	err = s.repository.ChangeConsumptionDuration(consumptionId, useCaseId, duration)
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	consumption, err = s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	return s.ValidateConsumptionResponse(&consumption), nil
// }

// // FormateConsumption формирует (завершает черновик) заявку потребления
// func (s *Service) FormateConsumption(
// 	consumptionId uint,
// 	userId uint,
// ) (dto.OneConsumptionResponse, error) {
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return dto.OneConsumptionResponse{}, ErrNoRecords
// 		}
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	if consumption.UserID != userId {
// 		return dto.OneConsumptionResponse{}, ErrForbidden
// 	}

// 	// Проверяем, что есть хотя бы один сценарий
// 	if len(consumption.Usecases) == 0 {
// 		return dto.OneConsumptionResponse{}, ErrBadRequest
// 	}

// 	// Проверяем корректность данных
// 	for _, uc := range consumption.Usecases {
// 		if uc.Duration == 0 || math.IsNaN(float64(uc.Duration)) {
// 			return dto.OneConsumptionResponse{}, ErrBadRequest
// 		}
// 	}

// 	err = s.repository.FormateConsumption(consumptionId)
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	consumption, err = s.repository.GetOneConsumption(consumptionId, "сформирован")
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	return s.ValidateConsumptionResponse(&consumption), nil
// }

// // ModeratorAction выполняет действие модератора (завершение/отклонение заявки)
// func (s *Service) ModeratorAction(
// 	consumptionId uint,
// 	action string,
// 	moderatorId uint,
// ) (dto.OneConsumptionResponse, error) {
// 	// Проверяем, что пользователь модератор
// 	moderator, err := s.repository.GetUser(moderatorId)
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, ErrNoRecords
// 	}
// 	if !moderator.IsModerator {
// 		return dto.OneConsumptionResponse{}, ErrForbidden
// 	}

// 	// Проверяем корректность действия
// 	if action != "завершен" && action != "отклонен" {
// 		return dto.OneConsumptionResponse{}, ErrBadRequest
// 	}

// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "сформирован")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return dto.OneConsumptionResponse{}, ErrNoRecords
// 		}
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	totalConsumption := 0.0
// 	if action == "завершен" {
// 		totalConsumption = float64(CalculateTotalConsumption(consumption.Usecases))
// 	}

// 	err = s.repository.ModeratorAction(consumptionId, action, totalConsumption, moderatorId)
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	consumption, err = s.repository.GetOneConsumption(consumptionId, action)
// 	if err != nil {
// 		return dto.OneConsumptionResponse{}, err
// 	}

// 	return s.ValidateConsumptionResponse(&consumption), nil
// }

// // DeleteConsumption удаляет черновик заявки потребления
// func (s *Service) DeleteConsumption(consumptionId uint, userId uint) error {
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ErrNoRecords
// 		}
// 		return err
// 	}

// 	if consumption.UserID != userId {
// 		return ErrForbidden
// 	}

// 	err = s.repository.DeleteConsumption(consumptionId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // CreateNewConsumption создает новую заявку потребления для пользователя
// func (s *Service) CreateNewConsumption(userId uint) (dto.NumberOfUseCasesResponse, error) {
// 	consumption := ds.Consumption{
// 		UserID:    userId,
// 		Status:    "черновик",
// 		CreatedAt: time.Now().Unix(),
// 		UpdatedAt: time.Now().Unix(),
// 	}

// 	err := s.repository.CreateConsumption(&consumption)
// 	if err != nil {
// 		return dto.NumberOfUseCasesResponse{}, err
// 	}

// 	return dto.NumberOfUseCasesResponse{
// 		ConsumptionId:    consumption.ID,
// 		NumberOfUseCases: 0,
// 	}, nil
// }

// // CalculateTotalConsumption вычисляет общее потребление энергии на основе сценариев
// func CalculateTotalConsumption(useCases []ds.usecase_consumption) uint {
// 	total := uint(0)
// 	for _, uc := range useCases {
// 		if uc.UseCase != nil {
// 			total += uc.UseCase.Consumption * uc.Duration
// 		}
// 	}
// 	return total
// }

// // usecase_consumption.go
// package service

// import (
// 	"lab/internal/app/ds"
// 	"errors"
// 	"gorm.io/gorm"
// 	"math"
// )

// // AddUseCaseToConsumption добавляет сценарий использования в заявку потребления
// func (s *Service) AddUseCaseToConsumption(
// 	userId uint,
// 	consumptionId uint,
// 	useCaseId uint,
// ) error {
// 	// Проверяем существование заявки и её принадлежность пользователю
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ErrNoRecords
// 		}
// 		return err
// 	}

// 	if consumption.UserID != userId {
// 		return ErrForbidden
// 	}

// 	// Проверяем существование и доступность сценария
// 	useCase, err := s.repository.GetUseCase(useCaseId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ErrNoRecords
// 		}
// 		return err
// 	}

// 	if useCase.IsDelete {
// 		return ErrUseCaseDeleted
// 	}

// 	// Добавляем сценарий в заявку
// 	err = s.repository.AddUseCaseToConsumption(useCaseId, consumptionId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // DeleteUseCaseFromConsumption удаляет сценарий из заявки потребления
// func (s *Service) DeleteUseCaseFromConsumption(
// 	userId uint,
// 	consumptionId uint,
// 	useCaseId uint,
// ) error {
// 	// Проверяем принадлежность заявки пользователю
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ErrNoRecords
// 		}
// 		return err
// 	}

// 	if consumption.UserID != userId {
// 		return ErrForbidden
// 	}

// 	// Проверяем наличие сценария в заявке
// 	exists := false
// 	for _, uc := range consumption.Usecases {
// 		if uc.UseCaseID == useCaseId {
// 			exists = true
// 			break
// 		}
// 	}

// 	if !exists {
// 		return ErrNoRecords
// 	}

// 	// Удаляем сценарий из заявки
// 	err = s.repository.DeleteUseCaseFromRequest(consumptionId, useCaseId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // ChangeUseCaseDurationInConsumption изменяет длительность сценария в заявке
// func (s *Service) ChangeUseCaseDurationInConsumption(
// 	userId uint,
// 	consumptionId uint,
// 	useCaseId uint,
// 	duration uint,
// ) (ds.usecase_consumption, error) {
// 	if duration == 0 || math.IsNaN(float64(duration)) {
// 		return ds.usecase_consumption{}, ErrBadRequest
// 	}

// 	// Проверяем принадлежность заявки пользователю
// 	consumption, err := s.repository.GetOneConsumption(consumptionId, "черновик")
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ds.usecase_consumption{}, ErrNoRecords
// 		}
// 		return ds.usecase_consumption{}, err
// 	}

// 	if consumption.UserID != userId {
// 		return ds.usecase_consumption{}, ErrForbidden
// 	}

// 	// Обновляем длительность
// 	err = s.repository.ChangeUseCaseConsumption(consumptionId, useCaseId, duration)
// 	if err != nil {
// 		return ds.usecase_consumption{}, err
// 	}

// 	// Получаем обновленный record
// 	response, err := s.repository.GetUseCaseFromConsumption(consumptionId, useCaseId)
// 	if err != nil {
// 		return ds.usecase_consumption{}, ErrNoRecords
// 	}

// 	return response, nil
// }

// // usecases.go
// package service

// import (
// 	dto "lab/internal/app/DTO"
// 	"lab/internal/app/ds"
// 	"errors"
// 	"gorm.io/gorm"
// )

// // GetUseCases возвращает список сценариев использования с фильтрацией по потреблению
// func (s *Service) GetUseCases(startValue uint, endValue uint) ([]ds.UseCase, error) {
// 	if startValue < 0 || endValue < 0 {
// 		return nil, ErrBadRequest
// 	}

// 	useCases, err := s.repository.GetUseCases(startValue, endValue)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrNoRecords
// 		}
// 		return nil, err
// 	}

// 	if len(useCases) == 0 {
// 		return nil, ErrNoRecords
// 	}

// 	return useCases, nil
// }

// // GetUseCase возвращает один сценарий использования по ID
// func (s *Service) GetUseCase(useCaseId uint) (ds.UseCase, error) {
// 	useCase, err := s.repository.GetUseCase(useCaseId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ds.UseCase{}, ErrNoRecords
// 		}
// 		return ds.UseCase{}, err
// 	}

// 	if useCase.IsDelete {
// 		return ds.UseCase{}, ErrUseCaseDeleted
// 	}

// 	return useCase, nil
// }

// // AddUseCase добавляет новый сценарий использования
// func (s *Service) AddUseCase(newUseCase dto.AddUseCase) (ds.UseCase, error) {
// 	if newUseCase.Name == "" || newUseCase.URL == "" || newUseCase.Consumption == 0 {
// 		return ds.UseCase{}, ErrBadRequest
// 	}

// 	useCaseToDB := ds.UseCase{
// 		Name:        newUseCase.Name,
// 		URL:         newUseCase.URL,
// 		Description: newUseCase.Description,
// 		Consumption: newUseCase.Consumption,
// 		IsDelete:    false,
// 	}

// 	response, err := s.repository.AddUseCase(useCaseToDB)
// 	if err != nil {
// 		return ds.UseCase{}, err
// 	}

// 	return response, nil
// }

// // DeleteUseCase удаляет (мягкое удаление) сценарий использования
// func (s *Service) DeleteUseCase(useCaseId uint) error {
// 	useCase, err := s.repository.GetUseCase(useCaseId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ErrNoRecords
// 		}
// 		return err
// 	}

// 	if useCase.IsDelete {
// 		return ErrUseCaseDeleted
// 	}

// 	err = s.repository.DeleteUseCase(useCaseId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // UpdateUseCase обновляет данные сценария использования
// func (s *Service) UpdateUseCase(useCaseId uint, updateData dto.ChangeUseCase) (ds.UseCase, error) {
// 	useCase, err := s.repository.GetUseCase(useCaseId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return ds.UseCase{}, ErrNoRecords
// 		}
// 		return ds.UseCase{}, err
// 	}

// 	if useCase.IsDelete {
// 		return ds.UseCase{}, ErrUseCaseDeleted
// 	}

// 	// Обновляем только переданные поля
// 	if updateData.Name != "" {
// 		useCase.Name = updateData.Name
// 	}
// 	if updateData.URL != "" {
// 		useCase.URL = updateData.URL
// 	}
// 	if updateData.Description != "" {
// 		useCase.Description = updateData.Description
// 	}
// 	if updateData.Consumption != 0 {
// 		useCase.Consumption = updateData.Consumption
// 	}

// 	err = s.repository.UpdateUseCase(useCaseId)
// 	if err != nil {
// 		return ds.UseCase{}, err
// 	}

// 	return useCase, nil
// }

// //users.go

// package service

// import (
// 	"errors"
// 	dto "lab/internal/app/DTO"
// 	"lab/internal/app/ds"

// 	"gorm.io/gorm"
// )

// // GetUserData возвращает информацию о пользователе
// func (s *Service) GetUserData(userId uint) (dto.UserDataResposne, error) {
// 	user, err := s.repository.GetUser(userId)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return dto.UserDataResposne{}, ErrNoRecords
// 		}
// 		return dto.UserDataResposne{}, err
// 	}

// 	return dto.UserDataResposne{
// 		ID:    user.ID,
// 		Login: user.Login,
// 	}, nil
// }

// // AddNewUser регистрирует нового пользователя
// func (s *Service) AddNewUser(user dto.UserRegistration) (dto.UserDataResposne, error) {
// 	if user.Login == "" || user.Password == "" {
// 		return dto.UserDataResposne{}, ErrBadRequest
// 	}

// 	userId, err := s.repository.AddNewUser(&ds.User{
// 		Login:       user.Login,
// 		Password:    user.Password,
// 		IsModerator: false,
// 	})
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}

// 	return s.GetUserData(userId)
// }

// // ChangeUserData изменяет данные пользователя
// func (s *Service) ChangeUserData(userId uint, userData dto.ChangeUserData) (dto.UserDataResposne, error) {
// 	if userData.Login == "" {
// 		return dto.UserDataResposne{}, ErrBadRequest
// 	}

// 	err := s.repository.ChangeUserData(userId, userData)
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}

// 	return s.GetUserData(userId)
// }

// // service.go
// package service

// import (
// 	"lab/internal/app/repository"
// 	"errors"
// 	"log"
// 	"os"

// 	"github.com/minio/minio-go/v7"
// 	"github.com/minio/minio-go/v7/pkg/credentials"
// )

// var (
// 	ErrNoRecords       = errors.New("записи не найдены")
// 	ErrForbidden       = errors.New("пользователь не имеет доступа к этой заявке")
// 	ErrBadRequest      = errors.New("введены некорректные данные")
// 	ErrUseCaseDeleted  = errors.New("этот сценарий использования удален")
// )

// type Service struct {
// 	repository  *repository.Repository
// 	minioClient *minio.Client
// }

// func NewService(repository *repository.Repository) *Service {
// 	minioClient, err := minio.New(os.Getenv("MINIO_HOST")+":"+os.Getenv("MINIO_PORT"), &minio.Options{
// 		Creds:  credentials.NewStaticV4("minio", "minio124", ""),
// 		Secure: false,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return &Service{
// 		repository:  repository,
// 		minioClient: minioClient,
// 	}
// }