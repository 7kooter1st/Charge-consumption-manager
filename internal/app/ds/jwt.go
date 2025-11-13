package ds

import (
	"LAB3/internal/app/role"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint      `json:"user_id"`
	Role   role.Role `json:"role"`
}
