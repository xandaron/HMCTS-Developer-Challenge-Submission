package main

import (
	"HMCTS-Developer-Challenge/database"
	// "HMCTS-Developer-Challenge/api"
	"html/template"
	"log"
	"net/http"
)

const (
	HomePage = iota
	LoginPage
	SignUpPage
)

var templates = make([]*template.Template, 3)

func main() {
	loadTemplates()

	if err := db.Connect(); err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", servePage(templates[HomePage]))
	http.HandleFunc("/login", servePage(templates[LoginPage]))
	http.HandleFunc("/signup", servePage(templates[SignUpPage]))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Switch to HTTPS in production
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}

	if err := db.Disconnect(); err != nil {
		log.Fatal(err)
	}
}

func loadTemplates() {
	const baseTemplate = "./templates/base.html"
	const navbarTemplate = "./templates/navbar.html"
	templates[HomePage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/home.html"))
	templates[LoginPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/login.html"))
	templates[SignUpPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/signup.html"))
}

func servePage(template *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := template.Execute(w, nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func apiWrapper(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}
