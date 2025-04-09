package session

import (
	"crypto/rand"
	"errors"
	"time"
)

const sessionTimeout = 0*time.Second + 5*time.Minute + 0*time.Hour

type Session struct {
	UserID    uint
	Timestamp time.Time
}

var sessions = make(map[string]Session)

func CreateUserSession(userID uint) (string, time.Time) {
	sessionID := rand.Text()
	timeStamp := time.Now()
	
	sessions[sessionID] = Session{
		UserID:    userID,
		Timestamp: timeStamp,
	}

	return sessionID, timeStamp.Add(sessionTimeout)
}

func GetUserID(sessionID string) (uint, error) {
	session, exists := sessions[sessionID]

	if !exists {
		return 0, errors.New("session not found")
	}

	if time.Since(session.Timestamp) > sessionTimeout {
		delete(sessions, sessionID)
		return 0, errors.New("session expired")
	}

	// Need to update the cookie timestamp
	session.Timestamp = time.Now()
	sessions[sessionID] = session

	return session.UserID, nil
}
