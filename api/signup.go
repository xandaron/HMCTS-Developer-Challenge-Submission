package api

import (
	db "HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/argon2"
)

var errUserExists = errors.Error("user already exists")

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

var config = PasswordConfig{
	time:    3,
	memory:  64 * 1024,
	threads: 4,
	keyLen:  32,
}

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
		return errors.AddContext(err, "signup.go: createUser - checkUserExists")
	} else if exists {
		return errUserExists
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.AddContext(err, "signup.go: createUser - GetDBHandle")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return errors.AddContext(err, "signup.go: createUser - rand.Read")
	}

	// Hash the password
	hash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)

	// Base64 encode for storage
	b64Hash := base64.StdEncoding.EncodeToString(hash)
	b64Salt := base64.StdEncoding.EncodeToString(salt)

	_, err = dbHandle.Exec(
		"INSERT INTO users (name, password_hash) VALUES (?, ?)",
		username,
		fmt.Sprintf(
			"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
			config.memory,
			config.time,
			config.threads,
			b64Salt,
			b64Hash,
		),
	)
	return err
}

func checkUserExists(username string) (bool, error) {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return false, errors.AddContext(err, "signup.go: checkUserExists - GetDBHandle")
	}

	var count int
	if err := dbHandle.QueryRow("SELECT COUNT(1) FROM users WHERE name = ?", username).Scan(&count); err != nil {
		return false, errors.AddContext(err, "signup.go: checkUserExists - QueryRow")
	}
	return count > 0, nil
}
