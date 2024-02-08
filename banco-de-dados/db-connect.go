package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	// connect to database
	db, err := sql.Open("mysql", "buddhilw:pass@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	return db
}
