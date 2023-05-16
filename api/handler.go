package api

import (
	"count-api/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const version = "0.1.0"

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["namespace"] + "_" + vars["key"]
	value := database.Get(key)
	if value != 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"value\":%v}", value)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{\"value\":null}")
	}
}

func Set(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["namespace"] + "_" + vars["key"]
	value := r.URL.Query().Get("value")
	byte_value, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		fmt.Println("Error: Cannot parse value to Uint64")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error: Please enter an positive integer value")
	} else {
		old_value, err := database.Set(key, byte_value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"old_value\":null,\"value\":null}")
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "{\"old_value\":%v,\"value\":%v}", old_value, value)
		}
	}
}

func Hit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["namespace"] + "_" + vars["key"]
	value := database.Get(key)
	_, err := database.Set(key, value+1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"value\":null}")
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"value\":%v}", value+1)
	}
}

func Stats(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
    key_count := database.CountKeys()
    fmt.Fprintf(w, "{\"total_keys\":%v,\"version\":%v}", key_count, version)
}
