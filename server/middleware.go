package main

import (
	"log"
	"net/http"
)

func logger(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Header)
}
