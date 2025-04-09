package api

import (
	"HMCTS-Developer-Challenge/database"
	"crypto/sha256"
	"fmt"
	"log"
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
		log.Printf("HandleSignUp: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	query := fmt.Sprintf("INSERT INTO users (name, password_sha256) VALUES ('%s', '%s')", username, sha256.Sum256([]byte(password)))

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return err
	}
	_, err = dbHandle.Exec(query)
	return err
}

func checkUserExists(username string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM users WHERE name='%s'", username)

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return false, err
	}

	var count int
	if err := dbHandle.QueryRow(query).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
