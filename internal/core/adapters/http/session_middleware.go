package http

import (
	"context"
	"net/http"
	"strings"

	"pulse/internal/platform/session"
)

type contextKey string

const SessionContextKey = contextKey("user_session")

type SessionMiddleware struct {
	store *session.Store
}

func NewSessionMiddleware(store *session.Store) *SessionMiddleware {
	return &SessionMiddleware{store: store}
}

func (m *SessionMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, `{"error":"Non autorisé"}`, http.StatusUnauthorized)
			return
		}

		sessData, err := m.store.Get(r.Context(), token)
		if err != nil {
			http.Error(w, `{"error":"Session invalide ou expirée"}`, http.StatusUnauthorized)
			return
		}

		// Injecte la session dans le contexte de la requête
		ctx := context.WithValue(r.Context(), SessionContextKey, sessData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Extraire le token depuis un Cookie ou le header Authorization: Bearer <token>
func extractToken(r *http.Request) string {
	if cookie, err := r.Cookie("session_token"); err == nil {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

// Helper pour récupérer la session facilement depuis un handler HTTP
func GetSessionFromContext(ctx context.Context) (*session.SessionData, bool) {
	sess, ok := ctx.Value(SessionContextKey).(*session.SessionData)
	return sess, ok
}
