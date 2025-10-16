package repository

import (
	"RIP/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

func (r *Repository) CheckIfUseCaseInConsumption(consumptionID, useCaseID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ds.UCCP{}).
		Where("consumption_id = ? AND use_case_id = ?", consumptionID, useCaseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DecreaseUseCaseDuration уменьшает время существующего сценария в заявке.
// Если время становится 0 или меньше, запись удаляется.
func (r *Repository) DecreaseUseCaseDuration(consumptionID, useCaseID, decrement uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var uccp ds.UCCP
		// Находим запись
		err := tx.Where("consumption_id = ? AND use_case_id = ?", consumptionID, useCaseID).First(&uccp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Если записи нет, ничего не делаем, это не ошибка
				return nil
			}
			return err // Другая ошибка
		}

		// Если уменьшение приведет к нулю или отрицательному значению
		if uccp.Duration <= decrement {
			// Удаляем запись
			return tx.Delete(&uccp).Error
		}

		// Иначе, просто уменьшаем время
		return tx.Model(&uccp).Update("duration", gorm.Expr("duration - ?", decrement)).Error
	})
}

// GetCurrentConsumption получает текущую заявку пользователя в статусе "черновик"
func (r *Repository) GetCurrentConsumption(userID uint) (ds.Consumption, error) {
	var consumption ds.Consumption
	err := r.db.Where("user_id = ? AND status = ?", userID, "черновик").First(&consumption).Error
	return consumption, err
}

// CreateConsumption создает новую заявку
func (r *Repository) CreateConsumption(userID uint) (ds.Consumption, error) {
	consumption := ds.Consumption{
		UserID: userID,
		Status: "черновик",
	}
	err := r.db.Create(&consumption).Error
	return consumption, err
}

// GetConsumptionWithUseCases получает заявку с включенными сценариями использования
func (r *Repository) GetConsumptionWithUseCases(consumptionID uint) (ds.Consumption, error) {
	var consumption ds.Consumption
	err := r.db.Preload("UseCases.UseCase").First(&consumption, consumptionID).Error
	return consumption, err
}

func (r *Repository) AddUseCaseToConsumption(consumptionID, useCaseID, duration uint) error {
	uccp := ds.UCCP{
		ConsumptionID: consumptionID,
		UseCaseID:     useCaseID,
		Duration:      duration,
	}
	return r.db.Create(&uccp).Error
}

// IncreaseUseCaseDuration увеличивает время существующего сценария в заявке
func (r *Repository) IncreaseUseCaseDuration(consumptionID, useCaseID, increment uint) error {
	return r.db.Model(&ds.UCCP{}).
		Where("consumption_id = ? AND use_case_id = ?", consumptionID, useCaseID).
		Update("duration", gorm.Expr("duration + ?", increment)).Error
}

// DeleteConsumption логически удаляет заявку (меняет статус на "удален")
func (r *Repository) DeleteConsumption(consumptionID uint) error {
	// Используем SQL UPDATE напрямую, как требуется
	result := r.db.Exec("UPDATE consumptions SET status = 'удален' WHERE id = ?", consumptionID)
	return result.Error
}

// CalculateTotalPower рассчитывает общее потребление энергии для заявки
func (r *Repository) CalculateTotalPower(consumptionID uint) (uint, error) {
	var totalPower uint
	err := r.db.Raw(`
		SELECT COALESCE(SUM(uc.power_consumption * uccp.duration), 0) as total_power
		FROM uccps uccp
		JOIN use_cases uc ON uccp.use_case_id = uc.id
		WHERE uccp.consumption_id = ?
	`, consumptionID).Scan(&totalPower).Error
	return totalPower, err
}

// UpdateConsumptionTotalPower обновляет общее потребление энергии в заявке
func (r *Repository) UpdateConsumptionTotalPower(consumptionID uint, totalPower uint) error {
	return r.db.Model(&ds.Consumption{}).Where("id = ?", consumptionID).Update("total_power", totalPower).Error
}

// RemoveUseCaseFromConsumption удаляет сценарий из заявки
func (r *Repository) RemoveUseCaseFromConsumption(consumptionID, useCaseID uint) error {
	return r.db.Where("consumption_id = ? AND use_case_id = ?", consumptionID, useCaseID).
		Delete(&ds.UCCP{}).Error
}
