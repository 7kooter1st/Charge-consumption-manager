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

// GetUserData возвращает информацию о пользователе (включая роль для фронта)
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
		Role:  uint8(user.Role),
	}, nil
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
	access, refresh, _, _, err := s.LoginUserWithMeta(creds)
	return access, refresh, err
}

// LoginUserWithMeta то же, что LoginUser, плюс userID и role для фронта (интерфейс модератора/клиента, Redux).
func (s *Service) LoginUserWithMeta(creds dto.UserRegistration) (accessToken, refreshToken string, userID uint, userRole role.Role, err error) {
	user, err := s.repository.GetUserByLogin(creds.Login)
	if err != nil {
		return "", "", 0, 0, ErrNoRecords
	}

	if !checkPasswordHash(creds.Password, user.Password) {
		return "", "", 0, 0, errors.New("invalid password")
	}

	accessToken, refreshToken, err = s.generateTokenPair(user.ID, user.Role)
	if err != nil {
		return "", "", 0, 0, err
	}
	return accessToken, refreshToken, user.ID, user.Role, nil
}

// RefreshTokens - валидирует refresh токен и выдает новую пару
func (s *Service) RefreshTokens(refreshTokenStr string) (string, string, error) {
	// 1. Парсим токен
	token, err := jwt.ParseWithClaims(refreshTokenStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return "", "", err
	}

	// 2. Валидируем claims
	if claims, ok := token.Claims.(*ds.JWTClaims); ok && token.Valid {
		// ВАЖНО: Проверяем, что нам подсунули именно Refresh токен, а не старый Access
		if !claims.IsRefresh {
			return "", "", fmt.Errorf("token is not a refresh token")
		}

		// 3. (Опционально) Проверяем, существует ли пользователь в БД до сих пор
		// user, err := s.repository.GetUser(claims.UserID)
		// if err != nil { return "", "", err }

		// 4. Генерируем новую пару токенов (ротация refresh токена)
		return s.generateTokenPair(claims.UserID, claims.Role)
	}

	return "", "", fmt.Errorf("invalid token")
}

// Вспомогательная функция для генерации пары токенов (чтобы не дублировать код)
func (s *Service) generateTokenPair(userID uint, userRole role.Role) (string, string, error) {
	// --- Access Token ---
	accessClaims := ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "charge-consumption-manager",
		},
		UserID:    userID,
		Role:      userRole,
		IsRefresh: false, // Это Access токен
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// --- Refresh Token ---
	refreshClaims := ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "charge-consumption-manager",
		},
		UserID:    userID,
		Role:      userRole,
		IsRefresh: true, // Это Refresh токен
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
