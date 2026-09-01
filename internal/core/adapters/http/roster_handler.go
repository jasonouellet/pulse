package http

import (
	"net/http"
)

func GetRostersHandler(w http.ResponseWriter, r *http.Request) {
	// Appel direct à GetSessionFromContext (même package http)
	sess, ok := GetSessionFromContext(r.Context())
	if !ok {
		http.Error(w, "Session manquante", http.StatusInternalServerError)
		return
	}

	if sess.ActiveRole != "COACH" && sess.ActiveRole != "CLUB_ADMIN" {
		http.Error(w, "Accès refusé pour ce rôle", http.StatusForbidden)
		return
	}

	// Logique métier avec sess.ActiveClub, etc.
}
