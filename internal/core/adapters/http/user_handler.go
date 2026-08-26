package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pulse/internal/core/ports"
)

// maxCreateUserBodyBytes borne la taille du corps accepté pour éviter
// qu'un payload démesuré ne consomme la mémoire du serveur.
const maxCreateUserBodyBytes = 1 << 20 // 1 MiB

type UserHandler struct {
	repo ports.UserRepository
}

func NewUserHandler(repo ports.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/core/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/", h.ListUsers)
		r.Get("/{id}", h.GetUserByID)
	})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxCreateUserBodyBytes)

	var req ports.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.PasswordHash = string(hash)

	user, err := h.repo.CreateUser(r.Context(), req)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get user", "error", err, "user_id", id)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := h.repo.ListUsers(r.Context(), int32(limit), int32(offset))
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}
