package repository

import "LAB3/internal/app/ds"

func (r *Repository) DeleteUseCaseFromRequest(Consumption_ID uint, UseCase_ID uint) error {
	//TODO Выполнить DELETE из таблицы request_panels
	err := r.db.
		Where("ConsumptionID = ? AND UseCaseID = ?", Consumption_ID, UseCase_ID).
		Delete(&ds.Consumption{}).Error
	if err != nil {
		return err
	}
	return nil
}

// func (r *Repository) ChangeUseCaseConsumption(Consumption_ID uint, UseCase_ID uint, duration_ uint) error {
// 	// Выполнить изменение площади у услуги в заявке
// 	err := r.db.Model(&ds.Consumption{}).
// 		Where("ConsumptionID = ? AND UseCaseID= ?", Consumption_ID, UseCase_ID).
// 		Update("duration", duration_).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *Repository) ChangeUseCaseConsumption(consumptionId uint, useCaseId uint, duration uint) error {
	err := r.db.Model(&ds.Usecase_consumption{}).
		Where("consumption_id = ? AND use_case_id = ?", consumptionId, useCaseId).
		Update("duration", duration).Error

	// Если в процессе выполнения запроса произошла ошибка, возвращаем ее.
	if err != nil {
		return err
	}

	// Если ошибок не было, возвращаем nil, сигнализируя об успешном выполнении.
	return nil
}

// func (r *Repository) GetUseCaseFromConsumption(Consumption_ID uint, UseCase_ID uint) (ds.Consumption, error) {
// 	var Consumption ds.Consumption
// 	err := r.db.Where("ConsumptionID = ? AND UseCAseID = ?", Consumption_ID, UseCase_ID).Preload("UseCase").
// 		First(&Consumption).Error
// 	if err != nil {
// 		return ds.Consumption{}, err
// 	}
// 	return Consumption, nil
// }

// GetUseCaseFromConsumption находит одну конкретную связь "сценарий-в-заявке"
// и возвращает ее.
func (r *Repository) GetUseCaseFromConsumption(consumptionId uint, useCaseId uint) (ds.Usecase_consumption, error) {
	var usecaseConsumption ds.Usecase_consumption
	// Ищем в таблице usecase_consumptions
	err := r.db.Model(&ds.Usecase_consumption{}).
		Where("consumption_id = ? AND use_case_id = ?", consumptionId, useCaseId).
		Preload("UseCase"). // Сразу подгружаем данные самого сценария, это удобно
		First(&usecaseConsumption).Error
	if err != nil {
		return ds.Usecase_consumption{}, err
	}
	return usecaseConsumption, nil
}
