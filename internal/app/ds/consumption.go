package ds

import (
	"gorm.io/gorm"
)

type Consumption struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`                            // ID пользователя
	Status      string `gorm:"size:50;not null;default:'черновик'"` // Статус заявки
	TotalPower  uint   `gorm:"default:0"`                           // Общее потребление энергии (рассчитывается при завершении)
	CreatedAt   int64  `gorm:"autoCreateTime"`
	ModeratedAt int64  `gorm:"autoModerateTime"` // Время создания
	UpdatedAt   int64  `gorm:"autoUpdateTime"`   // Время обновления
	ModeratorID *uint

	Usecases  []Usecase_consumption `gorm:"foreignKey:UseCaseID"`
	User      User
	Moderator User
}
