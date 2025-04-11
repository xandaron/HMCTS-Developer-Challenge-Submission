package session

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, name string, sessionID string, sessionTimout time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    sessionID,
		Path:     "/",
		Expires:  sessionTimout,
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	})
}
