package repository

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
	"time"
)

func (r *Repository) GetUseCasesInConsumption(userId uint) (uint, int64, error) {
	// Вернуть id заявки черновика и количество сценариев в этой заявке
	var (
		consumptionId    uint
		numberOfUseCases int64
	)

	// Ищем черновик заявки на потребление
	err := r.db.Model(&ds.Consumption{}).
		Where("user_id = ? AND status = ?", userId, "черновик").
		Select("id").
		First(&consumptionId).Error
	if err != nil {
		return 0, 0, err
	}

	// Считаем количество сценариев в заявке
	err = r.db.Model(&ds.Usecase_consumption{}).
		Where("ConsumptionID = ?", consumptionId).
		Count(&numberOfUseCases).Error
	if err != nil {
		return consumptionId, 0, err
	}

	return consumptionId, numberOfUseCases, nil
}

func (r *Repository) GetFilteredConsumptions(userId uint, filter dto.ConsumptionFilter) ([]ds.Consumption, error) {
	// Вернуть заявки для пользователя, кроме заявок со статусом удален и черновик, отфильтрованные по дате и статусу
	var consumptions []ds.Consumption
	db := r.db.Model(&ds.Consumption{}).
		Where("user_id = ? AND status NOT IN ('черновик','удален')", userId)

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if !filter.Start_date.IsZero() {
		db = db.Where("created_at >= ?", filter.Start_date)
	}
	if !filter.End_date.IsZero() {
		db = db.Where("created_at <= ?", filter.End_date)
	}

	err := db.Preload("User").Preload("Moderator").Find(&consumptions).Error
	if err != nil {
		return []ds.Consumption{}, err
	}
	return consumptions, nil
}

func (r *Repository) GetOneConsumption(consumptionId uint, status string) (ds.Consumption, error) {
	// Вернуть черновик заявки и её сценарии использования
	var consumption ds.Consumption
	err := r.db.Where("id = ? AND status = ?", consumptionId, status).
		Preload("UseCases.UseCase").
		First(&consumption).Error
	if err != nil {
		return ds.Consumption{}, err
	}
	return consumption, nil
}

func (r *Repository) ChangeConsumptionDuration(consumptionId uint, useCaseId uint, duration uint) error {
	// Изменить длительность использования сценария в заявке
	err := r.db.Model(&ds.Usecase_consumption{}).
		Where("consumption_id = ? AND use_case_id = ?", consumptionId, useCaseId).
		Update("duration", duration).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FormateConsumption(consumptionId uint) error {
	// Изменить статус черновика пользователя и проставить дату формирования
	err := r.db.Model(&ds.Consumption{}).
		Where("id = ? AND status = 'черновик'", consumptionId).
		Updates(map[string]any{
			"updated_at": time.Now(),
			"status":     "сформирован",
		}).Error
	if err != nil {
		return err
	}
	return nil
}

// /////

func (r *Repository) ModeratorAction(ConsumptionId uint, action string, totalConsumption float64, moderatorId uint) error {
	//TODO Отклонение/Завершение заявки модератором, проставить модератора, дату действия, рассчитать поле итоговой мощности
	err := r.db.Model(&ds.Consumption{}).
		Where("id = ? AND status = 'сформирован'", ConsumptionId).
		Updates(map[string]any{
			"UpdatedAt":   time.Now(),
			"Status":      action,
			"Consumption": totalConsumption,
			"Moderator":   moderatorId,
		}).Error

	if err != nil {
		return err
	}
	return nil
}

//////

func (r *Repository) DeleteConsumption(consumptionId uint) error {
	// Проставить статус заявки удален и прописать дату удаления
	err := r.db.Model(&ds.Consumption{}).
		Where("id = ? AND status = 'черновик'", consumptionId).
		Updates(map[string]any{
			"DeletedAt": time.Now(),
			"Status":    "удален",
		}).Error
	if err != nil {
		return err
	}
	return nil
}
