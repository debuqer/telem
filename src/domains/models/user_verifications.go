package models

import (
	"log"
	"time"

	"github.com/debuqer/telem/src/domains/services"
)

type UserVerification struct {
	UserId    int
	Type      string
	Code      string
	CreatedAt time.Time
}

func (userVerficiation *UserVerification) GenerateVerificationCode() error {
	db, _ := services.InitDB()

	userVerficiation.CreatedAt = time.Now()

	stmt, err := db.Prepare(`
	 INSERT INTO go.user_verifications
	 (user_id,
	 type,
	 code,
	 created_at)
	 VALUES
	 (?,?,?,?);
	`)

	if err != nil {
		log.Fatal(err)
		return err
	}

	stmt.Exec(
		userVerficiation.UserId,
		userVerficiation.Type,
		userVerficiation.Code,
		userVerficiation.CreatedAt,
	)

	return nil
}
