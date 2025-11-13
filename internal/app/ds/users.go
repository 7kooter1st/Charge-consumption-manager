package ds

import (
	"LAB3/internal/app/role"
	"sync"
)

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Login       string    `gorm:"type:varchar(255);unique"`
	Password    string    `gorm:"type:varchar(255)" json:"-"` // <-- Добавьте json:"-"
	IsModerator bool      // Мы заменим это поле на Role
	Role        role.Role `gorm:"default:0"` // <-- ДОБАВЛЕНО
}

type singletonUser struct {
	id uint
}

var (
	instance *singletonUser
	once     sync.Once
)

func GetUser() *singletonUser {
	once.Do(func() {
		instance = &singletonUser{1}
	})
	return instance
}

func (s *singletonUser) GetId() uint {
	return s.id
}
