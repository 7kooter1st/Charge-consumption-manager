package ds

import "gorm.io/gorm"

type UseCase struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey"`
	Title            string `gorm:"size:255;not null"` // Название сценария использования
	Description      string `gorm:"size:500"`          // Описание сценария
	ImageURL         string `gorm:"size:500"`          // URL изображения
	PowerConsumption uint   `gorm:"not null"`          // Потребление энергии в мА
	IsActive         bool   `gorm:"default:true"`      // Активен ли сценарий
}
