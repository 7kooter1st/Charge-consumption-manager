package handler

import (
	"LAB3/internal/app/ds"
	"LAB3/internal/app/role" // Используем сервис для доступа к конфигу
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "userRole"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}

	token, err := jwt.ParseWithClaims(headerParts[1], &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(h.Service.GetConfig().JWT.Secret), nil // <-- Получаем секрет через сервис
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token claims are not of type *JWTClaims"})
		return
	}

	// Сохраняем ID и роль пользователя в контекст для дальнейшего использования
	c.Set(userCtx, claims.UserID)
	c.Set(roleCtx, claims.Role)
}

func (h *Handler) requireRole(requiredRole role.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.Get(roleCtx)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user role not found in context"})
			return
		}

		if userRole.(role.Role) < requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you don't have enough permissions"})
			return
		}
	}
}
