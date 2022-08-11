package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Initialize() {
	var err error
	db, err = sql.Open("mysql", "root:my-secret-pw@tcp(db:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing database...")
}

func Ping() {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database online...")
}