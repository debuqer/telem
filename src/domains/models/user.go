package models

import (
	"log"
	"time"

	"github.com/debuqer/telem/src/domains/services"
)

type User struct {
	Id        int
	Username  string
	CreatedAt time.Time
	Password  string
}

func (model *User) Insert() error {
	db, _ := services.InitDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users( username, password, created_at) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt.Exec(model.Username, model.Password, model.CreatedAt)

	return err
}
