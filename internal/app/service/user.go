package service

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

// GetUserData возвращает информацию о пользователе
func (s *Service) GetUserData(userId uint) (dto.UserDataResposne, error) {
	user, err := s.repository.GetUser(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserDataResposne{}, ErrNoRecords
		}
		return dto.UserDataResposne{}, err
	}

	return dto.UserDataResposne{
		ID:    user.ID,
		Login: user.Login,
	}, nil
}

// AddNewUser регистрирует нового пользователя
func (s *Service) AddNewUser(user dto.UserRegistration) (dto.UserDataResposne, error) {
	if user.Login == "" || user.Password == "" {
		return dto.UserDataResposne{}, ErrBadRequest
	}

	userId, err := s.repository.AddNewUser(&ds.User{
		Login:       user.Login,
		Password:    user.Password,
		IsModerator: false,
	})
	if err != nil {
		return dto.UserDataResposne{}, err
	}

	return s.GetUserData(userId)
}

// ChangeUserData изменяет данные пользователя
func (s *Service) ChangeUserData(userId uint, userData dto.ChangeUserData) (dto.UserDataResposne, error) {
	if userData.Login == "" {
		return dto.UserDataResposne{}, ErrBadRequest
	}

	err := s.repository.ChangeUserData(userId, userData)
	if err != nil {
		return dto.UserDataResposne{}, err
	}

	return s.GetUserData(userId)
}

// package service

// import (
// 	dto "lab/internal/app/DTO"
// 	"lab/internal/app/ds"
// )

// func (s *Service) GetUserData(userId uint) (dto.UserDataResposne, error) {

// 	user, err := s.repository.GetUser(userId)
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}

// 	return dto.UserDataResposne{Login: user.Login,
// 		ID: user.ID}, nil
// }

// func (s *Service) AddNewUser(user dto.UserRegistration) (dto.UserDataResposne, error) {
// 	userId, err := s.repository.AddNewUser(&ds.User{Login: user.Login,
// 		Password: user.Password})
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}
// 	return s.GetUserData(userId)
// }

// func (s *Service) ChangeUserData(user dto.ChangeUserData) (dto.UserDataResposne, error) {
// 	err := s.repository.ChangeUserData(ds.GetUser().GetId(), user)
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}
// 	response, err := s.repository.GetUser(ds.GetUser().GetId())
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}
// 	return dto.UserDataResposne{ID: response.ID, Login: response.Login}, nil
// }
