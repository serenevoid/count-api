package main

import (
	"count-api/api"
	"count-api/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const version = "0.1.0"

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/get/{namespace}/{key}", api.Get).Methods("GET")
    r.HandleFunc("/set/{namespace}/{key}", api.Set).Methods("GET")
    r.HandleFunc("/hit/{namespace}/{key}", api.Hit).Methods("GET")
    r.HandleFunc("/stats", api.Stats).Methods("GET")
    go database.Monitor()
    log.Fatal(http.ListenAndServe(":8080", r))
    defer database.CloseDB()
}
