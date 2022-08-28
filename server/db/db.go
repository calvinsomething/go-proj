package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(db:%s)/%s?multiStatements=true",
			USER,
			PASSWORD,
			PORT,
			dbName,
		))
		if err != nil {
			log.Fatalln(err)
		}
		*dbs[i] = db

		start := time.Now()
		for {
			if err := db.Ping(); err == nil {
				break
			} else if time.Since(start) > 30*time.Second {
				log.Println(err)
				CleanUp()
				os.Exit(1)
			}
		}
		log.Printf("Connected to %s...\n", dbName)
	}
}

func CleanUp() {
	log.Println("Closing all db connections...")
	for _, db := range dbs {
		(*db).Close()
	}
}

func Up(db *sql.DB, dbName string) error {
    driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

    m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://db/migrations/%s", dbName),
        "mysql",
        driver,
    )
	if err != nil {
		return err
	}

	return m.Up()
}

func Down(db *sql.DB, dbName string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://db/migrations/%s", dbName),
        "mysql", 
        driver,
    )
	if err != nil {
		return err
	}

	return m.Down()
}

func Test(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO test
		values (1, 2, 3, 4)
	`)

	if err != nil {
		return err
	}

	rows, err := db.Query("select * from test")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var o, t, h, f int

		if err = rows.Scan(&o, &t, &h, &f); err != nil {
			return err
		}

		log.Println(o, t, h, f)
	}

	return nil
}