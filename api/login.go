package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/session"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, err := loginUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sessionID, sessionTimout := session.CreateUserSession(userID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  sessionTimout,
		HttpOnly: true,
		Secure:   false, // Set to true in production
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func loginUser(username, password string) (uint, error) {
	if username == "" || password == "" {
		return 0, fmt.Errorf("username and password cannot be empty")
	}

	query := fmt.Sprintf("SELECT id, password_sha256 FROM users WHERE name = '%s'", username)
	var queryResponse struct {
		ID             uint
		PasswordSha256 string
	}

	if err := db.GetDBHandle().QueryRow(query).Scan(&queryResponse); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user does not exist")
		} else {
			return 0, err
		}
	}

	hashedPassword := sha256.Sum256([]byte(password))
	if queryResponse.PasswordSha256 != string(hashedPassword[:]) {
		return 0, fmt.Errorf("incorrect password")
	}

	return queryResponse.ID, nil
}
