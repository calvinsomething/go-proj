package db

import (
	"context"
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
	Pool ConnectionPool
)

type (
	ConnectionPool struct {
		*sql.DB
	}

	Logger struct {
		*log.Logger
		verbose bool
	}
)

// Verbose is a method for implementing migrate.Logger.
func (l Logger) Verbose() bool {
	return l.verbose
}

// Initialize connects the database and pings to confirm.
func Initialize(USER, PASSWORD, PORT, NAME string) {
	var err error
	pool, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(db:%s)/%s?multiStatements=true",
		USER,
		PASSWORD,
		PORT,
		NAME,
	))
	if err != nil {
		log.Fatalln(err)
	}
	Pool = ConnectionPool{pool}

	fmt.Printf("%s:%s@tcp(db:%s)/%s?multiStatements=true\n",
		USER,
		PASSWORD,
		PORT,
		NAME,
	)

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

func (p *ConnectionPool) getMigrateInstance() (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(p.DB, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return nil, err
	}

	m.Log = Logger{log.Default(), true}

	return m, nil
}

// Migrate runs all migrations up, or down if the down param is true.
func (p *ConnectionPool) Migrate(down ...bool) error {
	m, err := p.getMigrateInstance()
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

// MigrateSteps takes a signed int and migrates up or down the number of steps passed in.
func (p *ConnectionPool) MigrateSteps(steps int) error {
	m, err := p.getMigrateInstance()
	if err != nil {
		return err
	}

	return m.Steps(steps)
}

func Test() {
	_, err := Pool.Exec(`
		INSERT INTO players
		values ('123451234', 'A', 'gnome', 'warlock', 'tailoring', null, 5)
	`)

	if err != nil {
		log.Println("ERRR", err)
	}
}

type ErrNoAffect struct {
	s string
}

func (e *ErrNoAffect) Error() string {
	return e.s
}

func (p *ConnectionPool) MustAffect(ctx context.Context, stmt string, args ...interface{}) error {
	res, err := p.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	} else if ra, err := res.RowsAffected(); err != nil {
		return err
	} else if ra == 0 {
		return &ErrNoAffect{"no rows affected"}
	}
	return nil
}
