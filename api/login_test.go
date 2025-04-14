package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandlerSuccess(t *testing.T) {
	// Create valid login data
	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser1",
		Password: "demo123",
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a login request
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the login handler
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	// Check that the status code is OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify a session cookie was created
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Error("No session cookie found after successful login")
	}
}

func TestLoginHandlerWrongPassword(t *testing.T) {
	// Create login data with incorrect password
	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser1",
		Password: "wrongpassword",
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a login request
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the login handler
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	// Check that the status code is Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body for error message
	var responseBody map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	if message, exists := responseBody["message"]; !exists || message != "incorrect password" {
		t.Errorf("Expected error message 'incorrect password', got '%v'", responseBody)
	}
}

func TestLoginHandlerUserNotFound(t *testing.T) {
	// Create login data with non-existent user
	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "nonexistentuser",
		Password: "anypassword",
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a login request
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the login handler
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	// Check that the status code is Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body for error message
	var responseBody map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	if message, exists := responseBody["message"]; !exists || message != "user not found" {
		t.Errorf("Expected error message 'user not found', got '%v'", responseBody)
	}
}

func TestLoginHandlerEmptyCredentials(t *testing.T) {
	// Test cases for empty credentials
	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"Empty Username", "", "password"},
		{"Empty Password", "username", ""},
		{"Empty Both", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginData := struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				Username: tc.username,
				Password: tc.password,
			}

			loginJSON, err := json.Marshal(loginData)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(LoginHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}

			var responseBody map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Failed to parse response as JSON: %v", err)
			}

			if message, exists := responseBody["message"]; !exists || message != "empty username or password" {
				t.Errorf("Expected error message 'empty username or password', got '%v'", responseBody)
			}
		})
	}
}

func TestLoginHandlerInvalidJSON(t *testing.T) {
	// Create invalid JSON data
	invalidJSON := []byte(`{"username": "testuser1", "password": }`) // Missing value for password

	// Create a login request
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the login handler
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	// Check that the status code is Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLoginHandlerWrongMethod(t *testing.T) {
	// Test with methods other than POST
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api/login", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(LoginHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for %s: got %v want %v",
					method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestParseHash(t *testing.T) {
	// Test valid password hash parsing
	validHash := "$argon2id$v=19$m=65536,t=3,p=4$pqJ2kWwyHs6Uszb0saO8wQ==$i93hewu5pDcYtrEjUSaKnd6yB00FLwIjWzpuOK5o9/Q="

	info, err := parseHash(validHash)
	if err != nil {
		t.Errorf("Failed to parse valid hash: %v", err)
	}

	if info.Algorithm != "argon2id" {
		t.Errorf("Expected algorithm 'argon2id', got '%s'", info.Algorithm)
	}
	if info.Version != 19 {
		t.Errorf("Expected version 19, got %d", info.Version)
	}
	if info.Memory != 65536 {
		t.Errorf("Expected memory 65536, got %d", info.Memory)
	}
	if info.Time != 3 {
		t.Errorf("Expected time 3, got %d", info.Time)
	}
	if info.Threads != 4 {
		t.Errorf("Expected threads 4, got %d", info.Threads)
	}
	if len(info.Salt) == 0 {
		t.Error("Salt should not be empty")
	}
	if len(info.Hash) == 0 {
		t.Error("Hash should not be empty")
	}
}

func TestParseHashInvalid(t *testing.T) {
	// Test cases with invalid hash formats
	invalidHashes := []struct {
		name string
		hash string
	}{
		{"Missing Parts", "$argon2id$v=19$m=65536,t=3,p=4"},
		{"Invalid Version", "$argon2id$vXYZ$m=65536,t=3,p=4$salt$hash"},
		{"Invalid Memory", "$argon2id$v=19$m=XYZ,t=3,p=4$salt$hash"},
		{"Invalid Time", "$argon2id$v=19$m=65536,t=XYZ,p=4$salt$hash"},
		{"Invalid Threads", "$argon2id$v=19$m=65536,t=3,p=XYZ$salt$hash"},
		{"Invalid Params Format", "$argon2id$v=19$m=65536;t=3;p=4$salt$hash"},
		{"Invalid Salt Encoding", "$argon2id$v=19$m=65536,t=3,p=4$!!!$hash"},
		{"Invalid Hash Encoding", "$argon2id$v=19$m=65536,t=3,p=4$c2FsdA==$!!!"},
	}

	for _, tc := range invalidHashes {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseHash(tc.hash)
			if err == nil {
				t.Errorf("Expected error for invalid hash '%s', got nil", tc.hash)
			}
		})
	}
}
