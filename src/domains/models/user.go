package models

import (
	"log"
	"time"

	"github.com/debuqer/telem/src/domains/services"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int
	Username  string
	CreatedAt time.Time
	Password  string
}

func (model *User) Insert() (int64, error) {
	db, _ := services.InitDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users( username, password, created_at) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(model.Password), 14)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(model.Username, password, model.CreatedAt)
	if err != nil {
		return 0, err
	}
	userId, err := result.LastInsertId()

	return userId, err
}
