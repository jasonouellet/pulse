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

// Test helper to initialize miniredis and the session store
func setupTestStore(t *testing.T) (*session.Store, *miniredis.Miniredis) {
	t.Helper()

	// 1. Start the local in-memory Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	// 2. Create the go-redis client bound to miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Automatic cleanup at the end of the test
	t.Cleanup(func() {
		_ = rdb.Close()
		mr.Close()
	})

	store := session.NewStore(rdb, 15*time.Minute)
	return store, mr
}

func TestSessionMiddleware_RequireAuth(t *testing.T) {
	store, mr := setupTestStore(t)
	mw := coreHTTP.NewSessionMiddleware(store)

	// Dummy handler protected by middleware to validate context
	protectedHandler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := coreHTTP.GetSessionFromContext(r.Context())
		if !ok {
			t.Error("Missing session in request context")
			http.Error(w, "no session", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("welcome " + sess.UserID))
	}))

	t.Run("Reject when no token is provided (401)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Allow access with valid Cookie (200)", func(t *testing.T) {
		ctx := reqContext()
		token := "sess_valid_cookie_123"

		// Prepare session in Redis via miniredis
		err := store.Create(ctx, token, session.SessionData{
			UserID:     "user_42",
			ActiveRole: "COACH",
			ActiveClub: "club_rimouski",
		})
		if err != nil {
			t.Fatalf("Error creating session: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		req.AddCookie(&http.Cookie{Name: "session_token", Value: token})
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
		if rec.Body.String() != "welcome user_42" {
			t.Errorf("Expected body 'welcome user_42', got '%s'", rec.Body.String())
		}
	})

	t.Run("Allow access with Bearer Authorization Header (200)", func(t *testing.T) {
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
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Reject after TTL expiration in Redis (401)", func(t *testing.T) {
		ctx := reqContext()
		token := "sess_expiring_789"

		_ = store.Create(ctx, token, session.SessionData{
			UserID: "user_expired",
		})

		// Fast forward miniredis simulated time by 16 minutes (TTL set to 15 min)
		mr.FastForward(16 * time.Minute)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 after TTL expiration, got %d", rec.Code)
		}
	})
}

func reqContext() context.Context {
	return context.Background()
}
