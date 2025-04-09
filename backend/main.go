package main

import (
	"HMCTS-Developer-Challenge/database"
	"log"
	"net/http"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
		return
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	if err := db.Disconnect(); err != nil {
		log.Fatal(err)
	}
}

func handler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}
