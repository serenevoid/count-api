package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var data map[string]uint64

func main() {
    data = make(map[string]uint64)
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/hit/", hitkey)
	http.HandleFunc("/get/", getkey)
	http.HandleFunc("/res/", reskey)
	http.HandleFunc("/del/", delkey)
	log.Fatal(http.ListenAndServe(":7456", nil))
}

// "/" to get welcome message
func homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Welcome to the Count API home page")
	fmt.Println("Endpoint hit: homepage")
}

// "/hit/key" or "/hit/namespace/key" to increment and return value of a key
func hitkey(w http.ResponseWriter, r *http.Request) {
	url_params := strings.Split(r.URL.Path, "/")
	if url_params[1] != "hit" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
    ns := ""
	if len(url_params) == 3 {
		ns = "default"
	}
    if len(url_params) == 4 {
        ns = url_params[3]
    }
	key_id := url_params[2]
	value, key_exists := data[ns+"_"+key_id]
	if key_exists {
		data[ns+"_"+key_id] += 1
		fmt.Fprintf(w, "{\"value\":%d}", value+1)
	} else {
		data[ns+"_"+key_id] = uint64(0)
		fmt.Fprintf(w, "{\"value\":%d}", value)
	}
	fmt.Println("Endpoint hit: hit")
    log.Print(data)
}

// "/get/key" or "/get/namespace/key" to get value of a key
func getkey(w http.ResponseWriter, r *http.Request) {
	url_params := strings.Split(r.URL.Path, "/")
	if url_params[1] != "get" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
    ns := ""
	if len(url_params) == 3 {
		ns = "default"
	}
    if len(url_params) == 4 {
        ns = url_params[3]
    }
	key_id := url_params[2]
	value, key_exists := data[ns+"_"+key_id]
	if key_exists {
		fmt.Fprintf(w, "{\"value\":%d}", value)
	} else {
		http.Error(w, "Key not present", http.StatusNotFound)
	}
	fmt.Println("Endpoint hit: get")
    log.Print(data)
}

// "/res`key" or "/res/namespace/key" to reset the value of a key
func reskey(w http.ResponseWriter, r *http.Request) {
	url_params := strings.Split(r.URL.Path, "/")
	if url_params[1] != "res" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
    ns := ""
	if len(url_params) == 3 {
		ns = "default"
	}
    if len(url_params) == 4 {
        ns = url_params[3]
    }
	key_id := url_params[2]
	_, key_exists := data[ns+"_"+key_id]
	if key_exists {
        data[ns+"_"+key_id] = 0
		fmt.Fprintf(w, "{\"value\":%d}", 0)
	} else {
		http.Error(w, "Key not present", http.StatusNotFound)
	}
	fmt.Println("Endpoint hit: res")
    log.Print(data)
}

// "/del/key" or "/del/namespace/key" to delete a key
func delkey(w http.ResponseWriter, r *http.Request) {
	url_params := strings.Split(r.URL.Path, "/")
	if url_params[1] != "del" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
    ns := ""
	if len(url_params) == 3 {
		ns = "default"
	}
    if len(url_params) == 4 {
        ns = url_params[3]
    }
	key_id := url_params[2]
	_, key_exists := data[ns+"_"+key_id]
	if key_exists {
        delete(data, ns+"_"+key_id)
		fmt.Fprint(w, "{\"status\": \"success\"}")
	} else {
		http.Error(w, "Key not present", http.StatusNotFound)
	}
	fmt.Println("Endpoint hit: del")
    log.Print(data)
}
