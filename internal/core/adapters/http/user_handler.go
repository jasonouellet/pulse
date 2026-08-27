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

// CreateUserRequest is the HTTP transport DTO for user creation. The plain
// text password remains in the handler and must not be passed to
// ports.CreateUserParams, which carries only the persistence hash.
type CreateUserRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     *string `json:"phone,omitempty"`
	Role      string  `json:"role"`
}

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

// CreateUser creates a user in the Core module.
//
// @Summary      Create a user
// @Description  Creates a user and stores a bcrypt hash of the supplied password.
// @Tags         Core users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User to create"
// @Success      201   {object}  ports.UserDTO
// @Failure      400   {string}  string  "Invalid request"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /api/v1/core/users/ [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxCreateUserBodyBytes)

	var req CreateUserRequest
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

	params := ports.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         req.Role,
	}

	user, err := h.repo.CreateUser(r.Context(), params)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// GetUserByID returns a user by its UUID.
//
// @Summary      Get a user
// @Description  Returns a Core user for the supplied UUID.
// @Tags         Core users
// @Produce      json
// @Param        id   path      string         true  "User UUID"  format(uuid)
// @Success      200  {object}  ports.UserDTO
// @Failure      400  {string}  string  "Invalid UUID format"
// @Failure      404  {string}  string  "User not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/v1/core/users/{id} [get]
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

// ListUsers returns a paginated list of users.
//
// @Summary      List users
// @Description  Returns Core users using offset pagination. The default limit is 20 and the maximum is 100.
// @Tags         Core users
// @Produce      json
// @Param        limit   query     integer  false  "Maximum number of users to return"  default(20)  minimum(1)  maximum(100)
// @Param        offset  query     integer  false  "Number of users to skip"             default(0)   minimum(0)
// @Success      200     {array}   ports.UserDTO
// @Failure      500     {string}  string  "Internal server error"
// @Router       /api/v1/core/users/ [get]
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
