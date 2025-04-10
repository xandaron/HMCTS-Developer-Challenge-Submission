package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"encoding/json"
	"net/http"
)

var errWrongPassword = fmt.Errorf("incorrect password")
var errUserNotFound = fmt.Errorf("user not found")
var errEmptyUsernameOrPassword = fmt.Errorf("empty username or password")

func HandleLogin(w http.ResponseWriter, r *http.Request) {
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

	userID, err := loginUser(jsonData.Username, jsonData.Password)
	if err == errWrongPassword || err == errUserNotFound || err == errEmptyUsernameOrPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": err.Error()}); err != nil {
			errors.HandleServerError(w, err, "login.go: HandleLogin - Encode")
			return
		}
		return
	} else if err != nil {
		errors.HandleServerError(w, err, "login.go: HandleLogin - loginUser")
		return
	}

	session.CreateUserSessionCookie(w, userID)

	w.WriteHeader(http.StatusOK)
}

func loginUser(username, password string) (uint, error) {
	if username == "" || password == "" {
		return 0, errEmptyUsernameOrPassword
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return 0, err
	}

	var userID uint
	var passwordSha256 string
	if err := dbHandle.QueryRow("SELECT id, password_sha256 FROM users WHERE name = ?", username).Scan(&userID, &passwordSha256); err != nil {
		if err == sql.ErrNoRows {
			return 0, errUserNotFound
		} else {
			return 0, err
		}
	}

	hashedPassword := sha256.Sum256([]byte(password))
	if passwordSha256 != string(hashedPassword[:]) {
		return 0, errWrongPassword
	}

	return userID, nil
}
