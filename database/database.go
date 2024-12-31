package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver mysql
)

func ConnectDatabase() (*sql.DB, error) {
	conn := "golang:golang@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, error := sql.Open("mysql", conn)

	if error != nil {
		return nil, error
	}

	if error = db.Ping(); error != nil {
		return nil, error
	}

	return db, nil
}
