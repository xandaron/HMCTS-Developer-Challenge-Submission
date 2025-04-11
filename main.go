package main

import (
	"HMCTS-Developer-Challenge/api"
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"bytes"
	"html/template"
	"log"
	"net/http"
)

const (
	HomePage = iota
	LoginSignUpPage
	TasksPage
	TasksAddEditPage

	PageCount
)

var templates = make([]*template.Template, PageCount)

func main() {
	loadTemplates()

	if err := db.Connect(); err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", servePageWithRedirect(templates[HomePage]))

	http.HandleFunc("/api/logout", apiWrapper(api.LogoutHandler))

	http.HandleFunc("/login", servePageSignupLogin(templates[LoginSignUpPage], "login", "Login"))
	http.HandleFunc("/api/login", apiWrapper(api.LoginHandler))

	http.HandleFunc("/signup", servePageSignupLogin(templates[LoginSignUpPage], "signup", "Create Account"))
	http.HandleFunc("/api/signup", apiWrapper(api.SignUpHandler))

	http.HandleFunc("/api/tasks/", apiWrapperWithSessionCheck(api.TasksHandler))
	http.HandleFunc("/tasks", servePageWithRedirect(templates[TasksPage]))
	http.HandleFunc("/tasks/add", servePageTask(templates[TasksAddEditPage], false))
	http.HandleFunc("/tasks/edit/", servePageTask(templates[TasksAddEditPage], true))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	go session.SessionCleanupRoutine()

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
	templates[LoginSignUpPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/login-signup.html"))
	templates[TasksPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/tasks.html"))
	templates[TasksAddEditPage] = template.Must(template.ParseFiles(baseTemplate, navbarTemplate, "./templates/add-edit_task.html"))
}

type pageData struct {
	IsLoggedIn bool
	Edit       bool
	Action     string
	SubmitText string
}

func servePageSignupLogin(template *template.Template, action string, submitText string) http.HandlerFunc {
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
			errors.HandleServerError(w, err, "main.go: servePageWithForm - Execute")
			return
		}
		buf.WriteTo(w)
	}
}

func servePageTask(template *template.Template, edit bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := session.GetUserIDFromSession(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var buf bytes.Buffer
		if err := template.Execute(&buf, &pageData{IsLoggedIn: true, Edit: edit}); err != nil {
			errors.HandleServerError(w, err, "main.go: servePageWithRedirect - Execute")
			return
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

		_, err := session.GetUserIDFromSession(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var buf bytes.Buffer
		if err := template.Execute(&buf, &pageData{IsLoggedIn: true}); err != nil {
			errors.HandleServerError(w, err, "main.go: servePageWithRedirect - Execute")
			return
		}
		buf.WriteTo(w)
	}
}

func apiWrapper(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}

func apiWrapperWithSessionCheck(fn func(http.ResponseWriter, *http.Request, uint)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		userID, err := session.GetUserIDFromSession(w, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fn(w, r, userID)
	}
}
