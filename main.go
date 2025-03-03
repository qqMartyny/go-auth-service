package main

import (
	"log"

	"github.com/gin-gonic/gin"

	// Наши внутренние пакеты
	"github.com/qqMartyny/go-auth-service/handlers"
	"github.com/qqMartyny/go-auth-service/middleware"
)

func main() {
	// Инициализируем базу данных
	db, err := initDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v\n", err)
	}
	defer db.Close()

	// Инициализируем роутер Gin
	r := gin.Default()

	// Регистрация и логин — без авторизации
	r.POST("/register", handlers.RegisterHandler(db))
	r.POST("/login", handlers.LoginHandler(db))

	// Группа маршрутов с авторизацией
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthRequired)
	{
		// Пример защищённого маршрута: GET /customers/:id
		authGroup.GET("/customers/:id", handlers.GetCustomerHandler(db))
	}

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v\n", err)
	}
}
