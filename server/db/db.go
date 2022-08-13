package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Initialize() {
	var err error
	db, err = sql.Open("mysql", "root:my-secret-pw@tcp(db:3306)/setup")
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	for {
		if err := db.Ping(); err == nil {
			break
		} else if time.Since(start) > 30*time.Second {
			log.Fatal(err)
		}
	}
	log.Println("Connected to database...")
}
