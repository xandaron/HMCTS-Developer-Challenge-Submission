package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

var errUserExists = fmt.Errorf("user-exists")

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jsonData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := createUser(jsonData.Username, jsonData.Password); err == errUserExists || err == errEmptyUsernameOrPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": err.Error()}); err != nil {
			errors.HandleServerError(w, err, "login.go: HandleLogin - Encode")
			return
		}
		return
	} else if err != nil {
		errors.HandleServerError(w, err, "signup.go: HandleSignUp - createUser")
		return
	}

	userID, err := db.GetUserID(jsonData.Username)
	if err != nil {
		errors.HandleServerError(w, err, "signup.go: HandleSignUp - GetUserID")
		return
	}
	session.CreateUserSessionCookie(w, userID)

	w.WriteHeader(http.StatusOK)
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

	passwordSha256 := sha256.Sum256([]byte(password))
	_, err = dbHandle.Exec("INSERT INTO users (name, password_sha256) VALUES (?, ?)", username, string(passwordSha256[:]))
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
