package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
)

var errWrongPassword = fmt.Errorf("incorrect-password")
var errUserNotFound = fmt.Errorf("user-not-found")
var errEmptyUsernameOrPassword = fmt.Errorf("username-password-cannot-be-empty")

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, err := loginUser(username, password)
	if err == errWrongPassword || err == errUserNotFound || err == errEmptyUsernameOrPassword {
		http.Redirect(w, r, fmt.Sprintf("/login?error=%s", err.Error()), http.StatusSeeOther)
		return
	} else if err != nil {
		errors.HandleServerError(w, err, "login.go: HandleLogin - loginUser")
		return
	}

	session.CreateUserSessionCookie(w, userID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
