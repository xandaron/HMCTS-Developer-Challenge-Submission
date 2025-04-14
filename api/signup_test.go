package api

import (
	"HMCTS-Developer-Challenge/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUpHandlerSuccess(t *testing.T) {
	// Create valid signup data
	signupData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "newtestuser",
		Password: "password123",
	}

	signupJSON, err := json.Marshal(signupData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a signup request
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(signupJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the signup handler
	handler := http.HandlerFunc(SignUpHandler)
	handler.ServeHTTP(rr, req)

	// Check that the status code is OK
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v, %s", rr.Code, http.StatusOK, rr.Body.String())
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
		t.Error("No session cookie found after successful signup")
	}

	// Try to log in with the new credentials to verify account creation
	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "newtestuser",
		Password: "password123",
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	loginHandler := http.HandlerFunc(LoginHandler)
	loginHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Login after signup failed with status code: %v", status)
	}

	// Log out to end the session
	logoutReq, err := http.NewRequest("POST", "/api/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookies = rr.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		t.Error("No session cookie found after successful signup")
	}
	logoutReq.AddCookie(sessionCookie)

	logoutRR := httptest.NewRecorder()
	logoutHandler := http.HandlerFunc(LogoutHandler)
	logoutHandler.ServeHTTP(logoutRR, logoutReq)

	if logoutRR.Code != http.StatusSeeOther {
		t.Errorf("Logout failed with status code: %v", logoutRR.Code)
	}

	db, err := database.GetDBHandle()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM users WHERE name = ?", signupData.Username); err != nil {
		t.Fatalf("Failed to delete test user from database: %v", err)
	}
}

func TestSignUpHandlerDuplicateUser(t *testing.T) {
	// Try to sign up with an existing username (testuser1)
	signupData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser1", // This user already exists in the test database
		Password: "password123",
	}

	signupJSON, err := json.Marshal(signupData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a signup request
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(signupJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the signup handler
	handler := http.HandlerFunc(SignUpHandler)
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

	if message, exists := responseBody["message"]; !exists || message != "user already exists" {
		t.Errorf("Expected error message 'user already exists', got '%v'", responseBody)
	}
}

func TestSignUpHandlerEmptyCredentials(t *testing.T) {
	// Test cases for empty credentials
	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"Empty Username", "", "password123"},
		{"Empty Password", "emptypassuser", ""},
		{"Empty Both", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signupData := struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				Username: tc.username,
				Password: tc.password,
			}

			signupJSON, err := json.Marshal(signupData)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(signupJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SignUpHandler)
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

func TestSignUpHandlerInvalidJSON(t *testing.T) {
	// Create invalid JSON data
	invalidJSON := []byte(`{"username": "newuser", "password": }`) // Missing value for password

	// Create a signup request
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the signup handler
	handler := http.HandlerFunc(SignUpHandler)
	handler.ServeHTTP(rr, req)

	// Check that the status code is Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSignUpHandlerWrongMethod(t *testing.T) {
	// Test with methods other than POST
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api/signup", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SignUpHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for %s: got %v want %v",
					method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}
