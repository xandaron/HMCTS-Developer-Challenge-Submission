package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateUserSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()
	CreateUserSessionCookie(w, 1)

	// Check if the cookie is set
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("Expected a cookie to be set, but none were found")
	}

	if cookies[0].Name != "session_id" {
		t.Errorf("Expected cookie name 'session_id', got %v", cookies[0].Name)
	}

	if cookies[0].Value == "" {
		t.Errorf("Expected a non-empty cookie value, got an empty string")
	}
}

func TestGetUserIDFromSession(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	sessions["test-session-id"] = Session{
		UserID:    1,
		Timestamp: time.Now(),
	}

	userID, err := GetUserIDFromSession(w, r)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if userID != 1 {
		t.Errorf("Expected user ID 1, got %v", userID)
	}

	delete(sessions, "test-session-id")
}

func TestDeleteUserSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	sessions["test-session-id"] = Session{
		UserID:    1,
		Timestamp: time.Now(),
	}

	err := DeleteUserSessionCookie(w, r)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := sessions["test-session-id"]; exists {
		t.Errorf("Expected session to be deleted, but it still exists")
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 0 && cookies[0].Name == "session_id" && cookies[0].Value != "" {
		t.Errorf("Expected cookie to be cleared, but it was not")
	}
}

func TestCreateUserSession(t *testing.T) {
	sessionID, timeout := createUserSession(1)

	if sessionID == "" {
		t.Fatalf("Expected a valid session ID, got an empty string")
	}

	session, exists := sessions[sessionID]
	if !exists {
		t.Fatalf("Expected session to exist, but it does not")
	}

	if session.UserID != 1 {
		t.Errorf("Expected UserID 1, got %v", session.UserID)
	}

	if temp := timeout.Sub(session.Timestamp); temp != sessionTimeout {
		t.Errorf("Expected timeout to be %v, got %v", sessionTimeout, temp)
	}

	delete(sessions, sessionID)
}

func TestGetSessionID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	sessions["test-session-id"] = Session{
		UserID:    1,
		Timestamp: time.Now(),
	}

	sessionID, err := getSessionID(w, r)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if sessionID != "test-session-id" {
		t.Errorf("Expected session ID 'test-session-id', got %v", sessionID)
	}

	delete(sessions, "test-session-id")
}

func TestGetUserID(t *testing.T) {
	sessions["test-session-id"] = Session{
		UserID:    1,
		Timestamp: time.Now(),
	}

	userID, err := getUserID("test-session-id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if userID != 1 {
		t.Errorf("Expected user ID 1, got %v", userID)
	}

	delete(sessions, "test-session-id")
}
