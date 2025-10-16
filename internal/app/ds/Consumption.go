package ds

import "gorm.io/gorm"

type Consumption struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null"`                            // ID пользователя
	Status     string `gorm:"size:50;not null;default:'черновик'"` // Статус заявки
	TotalPower uint   `gorm:"default:0"`                           // Общее потребление энергии (рассчитывается при завершении)
	CreatedAt  int64  `gorm:"autoCreateTime"`                      // Время создания
	UpdatedAt  int64  `gorm:"autoUpdateTime"`                      // Время обновления

	User     Users  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UseCases []UCCP `gorm:"foreignKey:ConsumptionID" json:"use_cases,omitempty"`
}
