package helpers

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func GetConn() (*sql.DB, error) {
	sqlSrc := "root:@tcp(127.0.0.1:3306)/telem"

	return sql.Open("mysql", sqlSrc)
}
