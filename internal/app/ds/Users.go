package ds

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"size:100;not null;unique"` // Имя пользователя
	Email    string `gorm:"size:255;unique"`          // Email пользователя
	IsActive bool   `gorm:"default:true"`             // Активен ли пользователь

	// Связи
	Consumptions []Consumption `gorm:"foreignKey:UserID" json:"-"`
}
