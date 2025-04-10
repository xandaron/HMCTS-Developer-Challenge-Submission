package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"crypto/sha256"
	"fmt"
	"net/http"
)

var errUserExists = fmt.Errorf("user-exists")

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if err := createUser(username, password); err == errUserExists || err == errEmptyUsernameOrPassword {
		http.Redirect(w, r, fmt.Sprintf("/signup?error=%s", err.Error()), http.StatusSeeOther)
		return
	} else if err != nil {
		errors.HandleServerError(w, err, "signup.go: HandleSignUp - createUser")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/api/login?username=%s&password=%s", username, password), http.StatusSeeOther)
}

func createUser(username, password string) error {
	if username == "" || password == "" {
		return errEmptyUsernameOrPassword
	}

	exists, err := checkUserExists(username)
	if err != nil {
		return err
	} else if exists {
		return errUserExists
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return err
	}
	_, err = dbHandle.Exec("INSERT INTO users (name, password_sha256) VALUES (?, ?)", username, sha256.Sum256([]byte(password)))
	return err
}

func checkUserExists(username string) (bool, error) {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return false, err
	}

	var count int
	if err := dbHandle.QueryRow("SELECT COUNT(1) FROM users WHERE name = ?", username).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
