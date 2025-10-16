package repository

import (
	"RIP/internal/app/ds"
)

func (r *Repository) GetUseCase(id int) (ds.UseCase, error) {
	var useCase ds.UseCase
	err := r.db.Where("is_active = ? AND id = ?", true, id).First(&useCase).Error
	return useCase, err
}

func (r *Repository) GetUseCases() ([]ds.UseCase, error) {
	var useCases []ds.UseCase
	err := r.db.Where("is_active = ?", true).Find(&useCases).Error
	return useCases, err
}

func (r *Repository) GetUseCasesByTitle(title string) ([]ds.UseCase, error) {
	var useCases []ds.UseCase
	err := r.db.Where("is_active = ? AND title ILIKE ?", true, "%"+title+"%").Find(&useCases).Error
	return useCases, err
}
