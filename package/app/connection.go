package app

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func DbConn() (db *sql.DB) {
	var err error

	db, err = sql.Open("postgres", "postgres://postgres:1234@localhost/contact?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	CheckError(err)
	return db
}
