package api

import (
	"HMCTS-Developer-Challenge/session"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = session.DeleteUserSessionCookie(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}