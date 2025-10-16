package ds

import "gorm.io/gorm"

type UCCP struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	ConsumptionID uint   `gorm:"not null"`      // ID заявки
	UseCaseID     uint   `gorm:"not null"`      // ID сценария использования
	Duration      uint   `gorm:"not null"`      // Время использования в минутах
	Order         uint   `gorm:"default:0"`     // Порядок в заявке
	IsMain        bool   `gorm:"default:false"` // Главный сценарий
	Notes         string `gorm:"size:500"`      // Дополнительные заметки

	// Связи
	Consumption Consumption `gorm:"foreignKey:ConsumptionID" json:"-"`
	UseCase     UseCase     `gorm:"foreignKey:UseCaseID" json:"use_case,omitempty"`

	// Составной уникальный ключ
	// gorm:"uniqueIndex:idx_consumption_usecase,unique" - будет добавлено в миграции
}
