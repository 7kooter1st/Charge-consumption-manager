package repository

import (
	"RIP/internal/app/ds"

	"gorm.io/gorm"
)

// GetUserByID получает пользователя по ID
func (r *Repository) GetUserByID(id uint) (ds.Users, error) {
	var user ds.Users
	err := r.db.First(&user, id).Error
	return user, err
}

// CreateUser создает нового пользователя
func (r *Repository) CreateUser(username, email string) (ds.Users, error) {
	user := ds.Users{
		Username: username,
		Email:    email,
		IsActive: true,
	}
	err := r.db.Create(&user).Error
	return user, err
}

// GetOrCreateDefaultUser получает или создает пользователя по умолчанию
func (r *Repository) GetOrCreateDefaultUser() (ds.Users, error) {
	var user ds.Users
	err := r.db.Where("username = ?", "default_user").First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// Создаем пользователя по умолчанию
		user, err = r.CreateUser("default_user", "default@example.com")
	}

	return user, err
}
