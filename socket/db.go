package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	user     = "postgres"
	dbname   = ""
	password = ""
	port     = 5432
	sslmode  = "disable"
)

func NewDatabase() (*sql.DB, error) {
	formattedSource := fmt.Sprintf(
		"host=%s port=%d user=%s password=%d dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)
	db, err := sql.Open("postgres", formattedSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}

func CreateQR(db *sql.DB, token string) (string, error) {
	query := fmt.Sprintf("INSERT INTO tokenqr (qr) VALUES ($1)")
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(query, token)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return token, nil
}

func ValidateQR(db *sql.DB, token string) (bool, error) {
	getQuery := fmt.Sprintf("SELECT qr FROM tokenqr WHERE qr = $1")
	rmQuery := fmt.Sprintf("DELETE FROM tokenqr WHERE qr = $1")

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	rows, err := tx.Query(getQuery, token)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	isQRExist := false
	if rows.Next() {
		isQRExist = true
	}
	rows.Close()

	if isQRExist {
		_, err := tx.Exec(rmQuery, token)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		tx.Commit()
		return true, nil
	} else {
		tx.Commit()
		return false, nil
	}
}
