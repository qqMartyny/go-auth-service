package repository

import (
	"database/sql"
	"errors"
	"time"
)

// Customer — модель для таблицы customers
type Customer struct {
	ID        int
	FirstName string
	LastName  string
	BirthDate time.Time
	Email     string
	Password  string // хранит bcrypt-хеш
}

// FindCustomerByEmail возвращает запись о покупателе, если email существует
func FindCustomerByEmail(db *sql.DB, email string) (*Customer, error) {
	query := `SELECT id, first_name, last_name, birth_date, email, password
              FROM customers
              WHERE email=$1`
	var c Customer
	err := db.QueryRow(query, email).Scan(
		&c.ID, &c.FirstName, &c.LastName, &c.BirthDate, &c.Email, &c.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // не нашли пользователя
		}
		return nil, err
	}
	return &c, nil
}

// InsertCustomer вставляет новую запись
func InsertCustomer(db *sql.DB, c *Customer) error {
	query := `INSERT INTO customers (first_name, last_name, birth_date, email, password)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, c.FirstName, c.LastName, c.BirthDate, c.Email, c.Password)
	return err
}

// FindCustomerByID ищет пользователя по ID
func FindCustomerByID(db *sql.DB, id int) (*Customer, error) {
	query := `SELECT id, first_name, last_name, birth_date, email, password
              FROM customers
              WHERE id=$1`
	var c Customer
	err := db.QueryRow(query, id).Scan(
		&c.ID, &c.FirstName, &c.LastName, &c.BirthDate, &c.Email, &c.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
