package main

import (
	_ "database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/message", handler(message))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}

type Message struct {
	Text string `json:"text"`
}

func message(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Message{Text: "Hello from Go!"}
	json.NewEncoder(w).Encode(response)
}
