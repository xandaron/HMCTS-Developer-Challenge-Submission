package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to perform login
func performLogin(t *testing.T) *http.Cookie {
	// Create login credentials for testuser1
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

	// Check if login was successful
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Login failed with status code: %v", status)
	}

	// Get the session cookie from the response
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Fatal("No session cookie found after login")
	}

	return sessionCookie
}

func TestLogoutHandlerSuccess(t *testing.T) {
	// First perform login to get a valid session cookie
	sessionCookie := performLogin(t)

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the session cookie to the request
	req.AddCookie(sessionCookie)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogoutHandler)

	// Call the handler with our test request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	// Check the redirect location is correct
	location := rr.Header().Get("Location")
	if location != "/login" {
		t.Errorf("handler returned wrong redirect location: got %v want %v", location, "/login")
	}

	// Check that the session cookie has been deleted (MaxAge < 0)
	cookies := rr.Result().Cookies()
	foundSessionCookie := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			foundSessionCookie = true
			if cookie.MaxAge > 0 {
				t.Errorf("session cookie was not properly deleted: MaxAge = %v, expected == 0", cookie.MaxAge)
			}
		}
	}

	if !foundSessionCookie {
		t.Errorf("no session cookie found in response")
	}
}

func TestLogoutHandlerWrongMethod(t *testing.T) {
	// Test with methods other than GET
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		// Create a request with the current method
		req, err := http.NewRequest(method, "/logout", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(LogoutHandler)

		// Call the handler
		handler.ServeHTTP(rr, req)

		// Check the status code is 405 Method Not Allowed
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code for %s: got %v want %v",
				method, status, http.StatusMethodNotAllowed)
		}
	}
}

func TestLogoutHandlerNoSession(t *testing.T) {
	// Test logout when there's no session cookie
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogoutHandler)

	// Call the handler without adding a session cookie
	handler.ServeHTTP(rr, req)

	// Should still redirect to login
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	// Check the redirect location is correct
	location := rr.Header().Get("Location")
	if location != "/login" {
		t.Errorf("handler returned wrong redirect location: got %v want %v", location, "/login")
	}
}
