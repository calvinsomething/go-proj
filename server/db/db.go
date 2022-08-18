package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	USER string
	PASSWORD string
	PORT string

	Players *sql.DB
	dbs = []**sql.DB{&Players}
	dbNames = []string{"players"}
)

func Initialize() {
	for i, dbName := range dbNames {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(db:%s)/%s",
			USER,
			PASSWORD,
			PORT,
			dbName,
		))
		if err != nil {
			log.Fatal(err)
		}
		*dbs[i] = db

		start := time.Now()
		for {
			if err := db.Ping(); err == nil {
				break
			} else if time.Since(start) > 30*time.Second {
				log.Fatal(err)
			}
		}
		log.Printf("Connected to %s...\n", dbName)
	}
}
