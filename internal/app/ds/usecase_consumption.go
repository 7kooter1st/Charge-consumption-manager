package ds

import (
	"gorm.io/gorm"
)

type Usecase_consumption struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	ConsumptionID uint `gorm:"not null"`
	UseCaseID     uint `gorm:"not null"`
	Duration      uint `gorm:"not null;default:1"`

	// Добавляем поле для хранения данных самого сценария при Preload
	UseCase UseCase `gorm:"foreignKey:UseCaseID"`
}
