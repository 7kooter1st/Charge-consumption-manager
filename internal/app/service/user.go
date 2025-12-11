package service

import (
	dto "LAB3/internal/app/DTO"
	"LAB3/internal/app/ds"
	"LAB3/internal/app/role"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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
// func (s *Service) AddNewUser(user dto.UserRegistration) (dto.UserDataResposne, error) {
// 	if user.Login == "" || user.Password == "" {
// 		return dto.UserDataResposne{}, ErrBadRequest
// 	}

// 	userId, err := s.repository.AddNewUser(&ds.User{
// 		Login:       user.Login,
// 		Password:    user.Password,
// 		IsModerator: false,
// 	})
// 	if err != nil {
// 		return dto.UserDataResposne{}, err
// 	}

// 	return s.GetUserData(userId)
// }

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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Проверка хеша пароля
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AddNewUser регистрирует нового пользователя
func (s *Service) AddNewUser(user dto.UserRegistration) (dto.UserDataResposne, error) {
	if user.Login == "" || user.Password == "" {
		return dto.UserDataResposne{}, ErrBadRequest
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return dto.UserDataResposne{}, err
	}

	// По умолчанию роль - User
	newUser := ds.User{
		Login:    user.Login,
		Password: hashedPassword,
		Role:     role.User,
	}

	userId, err := s.repository.AddNewUser(&newUser)
	if err != nil {
		return dto.UserDataResposne{}, err
	}

	return s.GetUserData(userId)
}

// LoginUser аутентифицирует пользователя и возвращает JWT токен
func (s *Service) LoginUser(creds dto.UserRegistration) (string, string, error) {
	user, err := s.repository.GetUserByLogin(creds.Login)
	if err != nil {
		return "", "", ErrNoRecords // Пользователь не найден
	}

	if !checkPasswordHash(creds.Password, user.Password) {
		return "", "", errors.New("invalid password")
	}

	// Создаем JWT Claims
	accessClaims := ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.ExpiresIn)), // Используем конфиг
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "charge-consumption-manager",
		},
		UserID:    user.ID,
		Role:      user.Role,
		IsRefresh: false,
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signedToken, err := token.SignedString([]byte(s.config.JWT.Secret)) // Используем секрет из конфига
	// if err != nil {
	// 	return "", fmt.Errorf("failed to sign token: %w", err)
	// }

	// return signedToken, nil
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// --- 4. Создаем Refresh Token (долгоживущий) ---
	refreshClaims := ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "charge-consumption-manager",
		},
		UserID:    user.ID,
		Role:      user.Role,
		IsRefresh: true, // <-- Указываем, что это refresh токен
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// 5. Возвращаем оба токена
	return accessToken, refreshToken, nil
}
