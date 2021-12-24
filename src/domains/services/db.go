package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	connection := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
		return db, nil
	}

	return db, err
}
