package main

import (
	"HMCTS-Developer-Challenge/api"
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/session"
	"bytes"
	"html/template"
	"log"
	"net/http"
)

const (
	HomePage = iota
	LoginSignUpPage
)

var templates = make([]*template.Template, 3)

func main() {
	loadTemplates()

	if err := db.Connect(); err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", servePageWithRedirect(templates[HomePage]))

	http.HandleFunc("/login", servePageWithForm(templates[LoginSignUpPage], "/api/login", "Login"))
	http.HandleFunc("/api/login", apiWrapper(api.HandleLogin))
	http.HandleFunc("/api/logout", apiWrapper(api.HandleLogout))

	// You probably don't want to allow users to sign up. This is just for testing purposes.
	http.HandleFunc("/signup", servePageWithForm(templates[LoginSignUpPage], "/api/signup", "Create Account"))
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
	templates[LoginSignUpPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/login-signup_form.html"))
}

type pageData struct {
	IsLoggedIn   bool
	Action     string
	SubmitText string
}

func servePageWithForm(template *template.Template, action string, submitText string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		loggedIn := false
		if _, err := session.GetUserIDFromSession(w, r); err == nil {
			loggedIn = true
		}

		var buf bytes.Buffer
		if err := template.Execute(&buf, &pageData{IsLoggedIn: loggedIn, Action: action, SubmitText: submitText}); err != nil {
			log.Printf("serverPageWithForm: %s\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		buf.WriteTo(w)
	}
}

func servePageWithRedirect(template *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if _, err := session.GetUserIDFromSession(w, r); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var buf bytes.Buffer
		if err := template.Execute(&buf, &pageData{IsLoggedIn: true}); err != nil {
			log.Printf("serverPageWithRedirect: %s\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		buf.WriteTo(w)
	}
}

func apiWrapper(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}
