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
	USER     string
	PASSWORD string
	PORT     string
	NAME     string

	Pool *sql.DB
)

func Initialize() {
	var err error
	Pool, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(db:%s)/%s?multiStatements=true",
		USER,
		PASSWORD,
		PORT,
		NAME,
	))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("SQL OPENED@!!", fmt.Sprintf("%s:%s@tcp(db:%s)/%s?multiStatements=true",
	USER,
	PASSWORD,
	PORT,
	NAME,
))

	start := time.Now()
	for {
		if err = Pool.Ping(); err == nil {
			break
		} else if time.Since(start) > 30*time.Second {
			Pool.Close()
			log.Println(err)
			os.Exit(1)
		}
		time.Sleep(500 * time.Millisecond)
	}
	log.Printf("Connected to database %s...\n", NAME)
}

func Migrate(db *sql.DB, down ...bool) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://db/migrations"),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if len(down) != 0 {
		if down[0] {
			return m.Down()
		}
	}

	return m.Up()
}

func Test(db *sql.DB) error {
	rows, err := db.Query("select * from players")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var i, h int
		var f, r, c, p1, p2 string

		if err = rows.Scan(&i, &f, &r, &c, &p1, &p2, &h); err != nil {
			return err
		}

		log.Println(i, f, r, c, p1, p2, h)
	}

	return nil
}
