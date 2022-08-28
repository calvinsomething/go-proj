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
}

func main() {
	db.Initialize()
	defer db.CleanUp()

	if err := db.Up(db.Players, "players"); err != nil {
		log.Fatal(err)
	}

	if err := db.Test(db.Players); err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Listening on port %s...\n", host_port)
	log.Fatal(http.ListenAndServe(":8080", newMux()))
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
	w.Write([]byte("hi"))
}
