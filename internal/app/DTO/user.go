package dto

type UserRegistration struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDataResposne struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  uint8  `json:"role"` // 0=User, 1=Moderator — для фронта (интерфейс модератора/клиента)
}

// LoginResponse — ответ логина для фронта (роль и id для роутинга и Redux).
type LoginResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
	UserID      uint   `json:"user_id"`
	Role        uint8  `json:"role"` // 0=User, 1=Moderator
}

type ChangeUserData struct {
	Login string `json:"login,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
