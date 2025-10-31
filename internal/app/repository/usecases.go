package repository

import (
	"LAB3/internal/app/ds"
	"fmt"
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

func (r *Repository) AddImageToSolarPanel(UseCaseID uint, imageURL string) error {
	return r.db.Model(&ds.UseCase{}).
		Where("id = ?", UseCaseID).
		Update("URL", imageURL).Error

}
