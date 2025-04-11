package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

var errUserExists = fmt.Errorf("user already exists")

// This should probably require system admin credentials in production
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(map[string]string{"message": err.Error()}); err != nil {
			errors.HandleServerError(w, err, "login.go: HandleLogin - Encode")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		buf.WriteTo(w)
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
		return errors.New(err, "signup.go: createUser - checkUserExists")
	} else if exists {
		return errUserExists
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.New(err, "signup.go: createUser - GetDBHandle")
	}

	passwordSha256 := sha256.Sum256([]byte(password))
	_, err = dbHandle.Exec("INSERT INTO users (name, password_sha256) VALUES (?, ?)", username, string(passwordSha256[:]))
	return err
}

func checkUserExists(username string) (bool, error) {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return false, errors.New(err, "signup.go: checkUserExists - GetDBHandle")
	}

	var count int
	if err := dbHandle.QueryRow("SELECT COUNT(1) FROM users WHERE name = ?", username).Scan(&count); err != nil {
		return false, errors.New(err, "signup.go: checkUserExists - QueryRow")
	}
	return count > 0, nil
}
