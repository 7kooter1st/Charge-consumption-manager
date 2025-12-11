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
	// 1. Получаем заголовок
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	// 2. Проверяем формат заголовка и извлекаем токен
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header format"})
		return
	}

	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		return
	}
	tokenStr := headerParts[1]

	// 3. Проверяем, не находится ли токен в черном списке
	inBlacklist, err := h.Service.IsInBlacklist(c.Request.Context(), tokenStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check token blacklist"})
		return
	}
	if inBlacklist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been logged out"})
		return
	}

	// 4. Парсим и валидируем подпись токена
	token, err := jwt.ParseWithClaims(tokenStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		// h.Service.GetConfig() должен существовать в вашем сервисе,
		// чтобы получать доступ к секретному ключу.
		return []byte(h.Service.GetConfig().JWT.Secret), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
		return
	}

	// 5. Проверяем claims и сохраняем данные в контекст
	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

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
