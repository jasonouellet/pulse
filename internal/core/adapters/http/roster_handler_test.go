package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	coreHTTP "pulse/internal/core/adapters/http"
	"pulse/internal/platform/session"
)

func TestGetRostersHandler_AccessControl(t *testing.T) {
	tests := []struct {
		name           string
		activeRole     string
		expectedStatus int
	}{
		{
			name:           "Allow COACH role",
			activeRole:     "COACH",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Allow CLUB_ADMIN role",
			activeRole:     "CLUB_ADMIN",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Forbidden for GUARDIAN role",
			activeRole:     "GUARDIAN",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Forbidden for PLAYER role",
			activeRole:     "PLAYER",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(coreHTTP.GetRostersHandler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)

			// Injecte la session dans le contexte HTTP
			sessData := &session.SessionData{
				UserID:     "user_test_123",
				ActiveRole: tt.activeRole,
				ActiveClub: "club_rimouski",
			}
			ctx := context.WithValue(req.Context(), coreHTTP.SessionContextKey, sessData)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetRostersHandler_MissingSession(t *testing.T) {
	handler := http.HandlerFunc(coreHTTP.GetRostersHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rosters", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d for missing session, got %d", http.StatusInternalServerError, rec.Code)
	}
}
