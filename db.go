package main

import (
	"database/sql"
	"fmt"
)

// initDB инициализирует подключение к PostgreSQL.
func initDB() (*sql.DB, error) {
	connStr := "user=postgres password=1209348756 dbname=auth_service_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение с БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось установить соединение с БД: %w", err)
	}

	fmt.Println("Успешное подключение к БД!")
	return db, nil
}
