package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"HMCTS-Developer-Challenge/session"
	"bytes"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashInfo struct {
	Algorithm string
	Version   int
	Memory    uint32
	Time      uint32
	Threads   uint8
	Salt      []byte
	Hash      []byte
}

var errUserNotFound = errors.Error("user not found")
var errWrongPassword = errors.Error("incorrect password")
var errEmptyUsernameOrPassword = errors.Error("empty username or password")

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	dbHandle, err := database.GetDBHandle()
	if err != nil {
		return 0, errors.AddContext(err, "login.go: loginUser - GetDBHandle")
	}

	var userID uint
	var passwordHash string
	if err := dbHandle.QueryRow("SELECT id, password_hash FROM users WHERE name = ?", username).Scan(&userID, &passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return 0, errUserNotFound
		} else {
			return 0, errors.AddContext(err, "login.go: loginUser - QueryRow")
		}
	}

	info, err := parseHash(passwordHash)
	if err != nil {
		return 0, errors.AddContext(err, "login.go: loginUser - parseHash")
	}

	newHash := argon2.IDKey(
		[]byte(password),
		info.Salt,
		info.Time,
		info.Memory,
		info.Threads,
		uint32(len(info.Hash)),
	)

	if subtle.ConstantTimeCompare(info.Hash, newHash) == 0 {
		return 0, errWrongPassword
	}

	return userID, nil
}

func parseHash(encodedHash string) (*HashInfo, error) {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return nil, errors.Errorf("invalid hash format: expected 6 parts, got %d", len(parts))
	}

	algorithm := parts[1]

	versionStr := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return nil, errors.Errorf("invalid version format: %v", err)
	}

	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return nil, errors.Error("invalid parameters format")
	}

	memoryStr := strings.TrimPrefix(params[0], "m=")
	memory, err := strconv.ParseUint(memoryStr, 10, 32)
	if err != nil {
		return nil, errors.Errorf("invalid memory parameter: %v", err)
	}

	timeStr := strings.TrimPrefix(params[1], "t=")
	time, err := strconv.ParseUint(timeStr, 10, 32)
	if err != nil {
		return nil, errors.Errorf("invalid time parameter: %v", err)
	}

	threadsStr := strings.TrimPrefix(params[2], "p=")
	threads, err := strconv.ParseUint(threadsStr, 10, 8)
	if err != nil {
		return nil, errors.Errorf("invalid threads parameter: %v", err)
	}

	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, errors.Errorf("invalid salt encoding: %v", err)
	}

	hash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, errors.Errorf("invalid hash encoding: %v", err)
	}

	return &HashInfo{
		Algorithm: algorithm,
		Version:   version,
		Memory:    uint32(memory),
		Time:      uint32(time),
		Threads:   uint8(threads),
		Salt:      salt,
		Hash:      hash,
	}, nil
}
