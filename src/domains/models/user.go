package models

import (
	"log"
	"time"

	"github.com/debuqer/telem/src/domains/services"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int
	Name            string
	Email           string
	EmailVerifiedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Password        string
	RememberToken   string
}

func (model *User) Insert() (int64, error) {
	db, _ := services.InitDB()
	defer db.Close()

	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	model.RememberToken = ""

	stmt, err := db.Prepare(`INSERT INTO go.users
	(
		name,
		email,
		password,
		created_at,
		updated_at
	)
	VALUES
	(?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(model.Password), 14)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(
		model.Name,
		model.Email,
		password,
		model.CreatedAt,
		model.UpdatedAt,
	)

	if err != nil {
		return 0, err
	}
	userId, err := result.LastInsertId()

	return userId, err
}
