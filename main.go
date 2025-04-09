package main

import (
	"HMCTS-Developer-Challenge/api"
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/session"
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
		log.Println(err)
		return
	}

	http.HandleFunc("/", serveHomePage)

	http.HandleFunc("/login", servePage(templates[LoginPage]))
	http.HandleFunc("/api/login", apiWrapper(api.HandleLogin))

	http.HandleFunc("/signup", servePage(templates[SignUpPage]))
	http.HandleFunc("/api/signup", apiWrapper(api.HandleSignUp))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Switch to HTTPS in production
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Println(err)
	}

	if err := db.Disconnect(); err != nil {
		log.Println(err)
	}
}

func loadTemplates() {
	const baseTemplate = "./templates/base.html"
	const navbarTemplate = "./templates/navbar.html"
	templates[HomePage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/home.html"))
	templates[LoginPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/login.html"))
	templates[SignUpPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/signup.html"))
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type Data struct {
		Name string
	}

	var data *Data
	if userID, err := session.GetUserIDFromSession(w, r); err == nil {
		if name, err := db.GetUserName(userID); err == nil {
			data = &Data{
				Name: name,
			}
		} else {
			data = &Data{
				Name: "Guest",
			}
		}
	} else {
		data = &Data{
			Name: "Guest",
		}
	}

	if err := templates[HomePage].Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func servePage(template *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
