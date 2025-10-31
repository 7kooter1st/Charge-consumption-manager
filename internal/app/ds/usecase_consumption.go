package ds

import (
	"gorm.io/gorm"
)

type Usecase_consumption struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	ConsumptionID uint `gorm:"not null"`
	UseCaseID     uint `gorm:"not null"`
	Duration      uint `gorm:"not null"`

	UseCases    []UseCase     `gorm:"foreignKey:UsecaseID"`     // ОШИБКА: имя поля
	Consumption []Consumption `gorm:"foreginKey:ConsumptionID"` // ОШИБКА: опечатка foreginKey
}
