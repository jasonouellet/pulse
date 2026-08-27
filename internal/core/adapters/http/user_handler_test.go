package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	coreHTTP "pulse/internal/core/adapters/http"
	"pulse/internal/core/ports"
)

// MockUserRepository implements ports.UserRepository in-memory for fast unit testing.
type MockUserRepository struct {
	users map[uuid.UUID]ports.UserDTO
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[uuid.UUID]ports.UserDTO)}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, p ports.CreateUserParams) (*ports.UserDTO, error) {
	id := uuid.New()
	u := ports.UserDTO{
		ID:        id,
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Role:      p.Role,
		IsActive:  true,
	}
	m.users[id] = u
	return &u, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*ports.UserDTO, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, ports.ErrUserNotFound
	}
	return &u, nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*ports.UserDTO, error) {
	return nil, ports.ErrUserNotFound
}

func (m *MockUserRepository) ListUsers(ctx context.Context, limit, offset int32) ([]ports.UserDTO, error) {
	var list []ports.UserDTO
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

func TestCreateUser_Success(t *testing.T) {
	repo := NewMockUserRepository()
	handler := coreHTTP.NewUserHandler(repo)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// Reflète le DTO HTTP réel (createUserRequest côté handler), pas
	// ports.CreateUserParams qui n'a pas de champ pour le mot de passe
	// en clair.
	payload := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	}{
		Email:     "coach@pulse.local",
		Password:  "S3cur3P@ssw0rd!",
		FirstName: "Jean",
		LastName:  "Dupont",
		Role:      "COACH",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/core/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	var created ports.UserDTO
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.Email != payload.Email {
		t.Errorf("expected email %s, got %s", payload.Email, created.Email)
	}
}
