package repository

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
)

func (r *Repository) GetUser(userID uint) (ds.User, error) {
	var user ds.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) AddNewUser(user *ds.User) (uint, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil

}

func (r *Repository) ChangeUserData(userId uint, userData dto.ChangeUserData) error {
	return r.db.Model(&ds.User{}).
		Where("id = ?", userId).
		Updates(map[string]any{
			"login": userData.Login,
		}).Error

}

func (r *Repository) GetUserByLogin(login string) (ds.User, error) {
	var user ds.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}
