package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Тот же самый ключ, что и в handlers
var secretKey = []byte("your-256-bit-secret")

// AuthRequired проверяет валидность JWT-токена,
// если токен невалиден — возвращаем 403 (Forbidden).
func AuthRequired(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Токен не предоставлен"})
		return
	}

	// Проверяем, что начинается с "Bearer "
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Некорректный заголовок Authorization"})
		return
	}

	// Парсим токен
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Невалидный токен: %v", err)})
		return
	}

	// Проверяем, действительно ли токен валиден
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Невалидный токен"})
		return
	}

	// Дополнительно можно проверить, что срок действия не истёк, но
	// jwt.ParseWithClaims уже смотрит на ExpiresAt, если оно в будущем.

	// Можно сохранить в Context ID пользователя (claims.Subject) для дальнейшего использования
	c.Set("user_id", claims.Subject)

	// Продолжаем выполнение цепочки
	c.Next()
}
