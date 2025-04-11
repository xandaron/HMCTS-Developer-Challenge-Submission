package session

import (
	"crypto/rand"
	"HMCTS-Developer-Challenge/errors"
	"net/http"
	"time"
)

const sessionTimeout = 0*time.Second + 5*time.Minute + 0*time.Hour

type Session struct {
	UserID    uint
	Timestamp time.Time
}

var sessions = make(map[string]Session)

var errSessionExpired = errors.NewBaseError("session-expired")
var errSessionNotFound = errors.NewBaseError("session-not-found")

func CreateUserSessionCookie(w http.ResponseWriter, userID uint) {
	sessionID, sessionTimout := createUserSession(userID)
	SetCookie(w, "session_id", sessionID, sessionTimout)
}

func GetUserIDFromSession(w http.ResponseWriter, r *http.Request) (uint, error) {
	sessionID, err := getSessionID(w, r)
	if err != nil {
		return 0, errors.New(err, "session.go: GetUserIDFromSession - getSessionID")
	}

	userID, err := getUserID(sessionID)
	if err != nil {
		return 0, errors.New(err, "session.go: GetUserIDFromSession - getUserID")
	}

	return userID, nil
}

func DeleteUserSessionCookie(w http.ResponseWriter, r *http.Request) error {
	sessionID, err := getSessionID(w, r)
	if err == nil {
		delete(sessions, sessionID)
		SetCookie(w, "session_id", "", time.Time{})
	}
	return err
}

func createUserSession(userID uint) (string, time.Time) {
	sessionID := rand.Text()
	timeStamp := time.Now()

	sessions[sessionID] = Session{
		UserID:    userID,
		Timestamp: timeStamp,
	}

	return sessionID, timeStamp.Add(sessionTimeout)
}

func getSessionID(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", errors.New(err, "session.go: getSessionID - Cookie")
	}

	sessionID := cookie.Value
	if time.Since(sessions[sessionID].Timestamp) > sessionTimeout {
		delete(sessions, sessionID)
		SetCookie(w, "session_id", "", time.Time{})
		return "", errSessionExpired
	}

	sessions[sessionID] = Session{
		UserID:    sessions[sessionID].UserID,
		Timestamp: time.Now(),
	}

	return sessionID, nil
}

func getUserID(sessionID string) (uint, error) {
	session, exists := sessions[sessionID]

	if !exists {
		return 0, errSessionNotFound
	}

	return session.UserID, nil
}
