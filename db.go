package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func getConn() (*sql.DB, error) {
	return sql.Open("mysql", sqlSrc)
}
