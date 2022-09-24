package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/calvinsomething/go-proj/db"
)

var (
	hostPort string
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	hostPort = os.Getenv("SERVER_PORT")

	db.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	defer db.Pool.Close()

	execArgs()

	middleware := []http.HandlerFunc{
		logger,
	}
	m := newMux(middleware...)

	m.get("/data", dataHandler)
	m.post("/login", loginHandler)
	m.post("/register", registerHandler)

	log.Printf("Listening on port %s...\n", hostPort)
	log.Fatal(m.ListenAndServe(":" + hostPort))
}

func execArgs() {
	i := 1
	shouldExit := false
	for len(os.Args) > i {
		switch os.Args[i] {
		case "migrate":
			shouldExit = true
			i++
			if len(os.Args) > i {
				if os.Args[i] == "up" {
					log.Println("Migrating database up...")
					if err := db.Pool.Migrate(); err != nil {
						log.Fatal(err)
					}
				} else if os.Args[i] == "down" {
					log.Println("Migrating database down...")
					if err := db.Pool.Migrate(true); err != nil {
						log.Fatal(err)
					}
				} else if steps, err := strconv.ParseInt(os.Args[i], 10, 32); err != nil {
					log.Fatalf("Invalid migrate option: %s", os.Args[i])
				} else {
					log.Printf("Migrating %d steps...", steps)
					if err = db.Pool.MigrateSteps(int(steps)); err != nil {
						log.Fatal(err)
					}
				}
			} else {
				log.Fatal("Missing migrate option")
			}
		default:
			log.Fatalf("Unknown command line argument: '%s'", os.Args[i])
		}
		i++
	}
	if shouldExit {
		os.Exit(0)
	}
}

func httpErr(w http.ResponseWriter, statusCode int, err error, msg ...string) {
	if err != nil {
		log.Output(2, err.Error())
	}

	w.WriteHeader(statusCode)

	if len(msg) != 0 {
		w.Write([]byte(msg[0]))
	}
}

func returnJSON(w http.ResponseWriter, data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		httpErr(w, 500, err)
	}
	w.Write(j)
}
