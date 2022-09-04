package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/calvinsomething/go-proj/db"
)

var (
	host_port string
)

type (
	handler_t func(http.ResponseWriter, *http.Request)

	mux_t struct {
		mux        *http.ServeMux
		middleware []handler_t
	}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Server
	host_port = os.Getenv("SERVER_PORT")

	// DB
	db.PORT = os.Getenv("DB_PORT")
	db.USER = os.Getenv("DB_USER")
	db.PASSWORD = os.Getenv("DB_PASSWORD")
	db.NAME = os.Getenv("DB_NAME")
}

func main() {
	db.Initialize()
	defer db.Pool.Close()

	execArgs()

	log.Printf("Listening on port %s...\n", host_port)
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}

func execArgs() {
	i := 1
	shouldExit := false
	for len(os.Args) > i {
		switch os.Args[i] {
		case "migrate":
			i++
			if len(os.Args) > i {
				if os.Args[i] == "up" {
					log.Println("Migrating database up...")
					if err := db.Migrate(db.Pool); err != nil {
						log.Fatal(err)
					}
				} else if os.Args[i] == "down" {
					log.Println("Migrating database down...")
					if err := db.Migrate(db.Pool, true); err != nil {
						log.Fatal(err)
					}
					shouldExit = true
				}
			}
		default:
			log.Fatalf("unknown command line argument: '%s'", os.Args[i])
		}
		i++
	}
	if shouldExit {
		os.Exit(0)
	}
}

// Custom Mux
func newMux() *mux_t {
	middleware := []handler_t{
		logger,
	}
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/err", errHandler)
	mux.HandleFunc("/data", dataHandler)

	return &mux_t{mux, middleware}
}

func (m *mux_t) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mw := range m.middleware {
		mw(w, r)
	}
	m.mux.ServeHTTP(w, r)
}

// Middleware
func logger(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Header)
}

// temp...
func errHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(struct {
		Err string
	}{"test err"})
	if err != nil {
		log.Println(err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Test(db.Pool); err != nil {
		log.Fatal(err)
	}
	log.Println("REMOTE=", r.RemoteAddr)
	w.Write([]byte("hiya"))
}
