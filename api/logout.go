package api

import (
	"HMCTS-Developer-Challenge/session"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session.DeleteUserSessionCookie(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}