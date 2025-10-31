package dto

type UserRegistration struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDataResposne struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
}

type ChangeUserData struct {
	Login string `json:"login,omitempty"`
}
