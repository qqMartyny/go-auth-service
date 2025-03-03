package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	// Импортируем наш репозиторий и основную "модель" (если нужно)
	mainrepo "github.com/qqMartyny/go-auth-service"
)

// Пример секретного ключа
var secretKey = []byte("your-256-bit-secret")

// ==== РЕГИСТРАЦИЯ (POST /register) ====
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type request struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			BirthDate string `json:"birth_date"` // в формате "YYYY-MM-DD"
			Email     string `json:"email"`
			Password  string `json:"password"`
		}

		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат входных данных"})
			return
		}

		// Проверяем, нет ли пользователя с таким email
		existing, err := mainrepo.FindCustomerByEmail(db, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
			return
		}
		if existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
			return
		}

		// Преобразуем дату (если нужно)
		var birthDate time.Time
		if req.BirthDate != "" {
			parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная дата рождения"})
				return
			}
			birthDate = parsedDate
		}

		// Хешируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
			return
		}

		// Сохраняем нового пользователя
		newCustomer := &mainrepo.Customer{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			BirthDate: birthDate,
			Email:     req.Email,
			Password:  string(hashedPassword),
		}
		if err := mainrepo.InsertCustomer(db, newCustomer); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Пользователь успешно зарегистрирован",
		})
	}
}

// ==== ЛОГИН (POST /login) ====
func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type loginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат входных данных"})
			return
		}

		// Ищем пользователя по email
		user, err := mainrepo.FindCustomerByEmail(db, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
			return
		}
		if user == nil {
			// Нет такого email
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Сравниваем хеши
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Генерируем JWT
		claims := jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),                         // ID в виде строки
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Срок действия 24 часа
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(secretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenStr})
	}
}

// ==== ПОЛУЧЕНИЕ ДАННЫХ О ПОЛЬЗОВАТЕЛЕ (GET /customers/:id) ====
func GetCustomerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ID из URL (строка)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
			return
		}

		// Ищем в базе
		user, err := mainrepo.FindCustomerByID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		// Возвращаем данные (без пароля)
		resp := gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"birth_date": user.BirthDate.Format("2006-01-02"),
			"email":      user.Email,
		}
		c.JSON(http.StatusOK, resp)
	}
}
