package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "github.com/joho/godotenv"
)

var (
	_port string
)

type (
	handler func(http.ResponseWriter, *http.Request)

	multiplexer struct {
		mux *http.ServeMux
		middleware []handler
	}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// if err := godotenv.Load("../.env"); err != nil {
	// 	log.Fatal(err)
	// }

	_port = os.Getenv("GO_PORT")
}

func main() {
	middleware := []handler{
		logger,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	log.Printf("Listening on port %s...\n", _port)
	log.Fatal(http.ListenAndServe(":" + _port, &multiplexer{mux, middleware}))
}

func (m *multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mw := range m.middleware {
		mw(w, r)
	}
	m.mux.ServeHTTP(w, r)
}

func logger(w http.ResponseWriter, r *http.Request) {
	log.Println("LOGGING:", r.URL.Path)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(struct{
		Err string
	}{"test err"})
	if err != nil {
		log.Println(err)
	}
}