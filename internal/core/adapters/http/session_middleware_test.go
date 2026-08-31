package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	coreHTTP "pulse/internal/core/adapters/http"
	"pulse/internal/platform/session"
)

// Helper de test pour initialiser miniredis et le magasin de session
func setupTestStore(t *testing.T) (*session.Store, *miniredis.Miniredis) {
	t.Helper()

	// 1. Démarrer le serveur Redis in-memory local
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Échec du démarrage de miniredis: %v", err)
	}

	// 2. Créer le client go-redis lié à miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Nettoyage automatique en fin de test
	t.Cleanup(func() {
		rdb.Close()
		mr.Close()
	})

	store := session.NewStore(rdb, 15*time.Minute)
	return store, mr
}

func TestSessionMiddleware_RequireAuth(t *testing.T) {
	store, mr := setupTestStore(t)
	mw := coreHTTP.NewSessionMiddleware(store)

	// Handler factice protégé par le middleware pour valider le contexte
	protectedHandler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := coreHTTP.GetSessionFromContext(r.Context())
		if !ok {
			t.Error("Session manquante dans le contexte de la requête")
			http.Error(w, "no session", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("bienvenue " + sess.UserID))
	}))

	t.Run("Rejet si aucun token fourni (401)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Statut attendu 401, obtenu %d", rec.Code)
		}
	})

	t.Run("Accès autorisé avec Cookie valide (200)", func(t *testing.T) {
		ctx := reqContext()
		token := "sess_valid_cookie_123"

		// Préparation de la session dans Redis via miniredis
		err := store.Create(ctx, token, session.SessionData{
			UserID:     "user_42",
			ActiveRole: "COACH",
			ActiveClub: "club_rimouski",
		})
		if err != nil {
			t.Fatalf("Erreur lors de la création de la session: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Statut attendu 200, obtenu %d", rec.Code)
		}
		if rec.Body.String() != "bienvenue user_42" {
			t.Errorf("Corps attendu 'bienvenue user_42', obtenu '%s'", rec.Body.String())
		}
	})

	t.Run("Accès autorisé avec En-tête Authorization Bearer (200)", func(t *testing.T) {
		ctx := reqContext()
		token := "sess_valid_bearer_456"

		_ = store.Create(ctx, token, session.SessionData{
			UserID:     "user_99",
			ActiveRole: "CLUB_ADMIN",
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Statut attendu 200, obtenu %d", rec.Code)
		}
	})

	t.Run("Rejet après expiration du TTL dans Redis (401)", func(t *testing.T) {
		ctx := reqContext()
		token := "sess_expiring_789"

		_ = store.Create(ctx, token, session.SessionData{
			UserID: "user_expired",
		})

		// Avancer le temps simulé de miniredis de 16 minutes (TTL configuré à 15 min)
		mr.FastForward(16 * time.Minute)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Statut attendu 401 après expiration du TTL, obtenu %d", rec.Code)
		}
	})
}

func reqContext() context.Context {
	return context.Background()
}
