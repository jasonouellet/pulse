package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pulse/internal/core/ports"
)

const maxCreateUserBodyBytes = 1 << 20

// CreateUserRequest is the HTTP transport DTO for user creation.
type CreateUserRequest struct {
	Email     string  `json:"email" doc:"User email address"`
	Password  string  `json:"password" minLength:"8" doc:"Plain text password; never returned"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     *string `json:"phone,omitempty"`
	Role      string  `json:"role"`
}

type CreateUserInput struct {
	Body CreateUserRequest
}

type GetUserInput struct {
	ID string `path:"id" format:"uuid" doc:"User UUID"`
}

type ListUsersInput struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Maximum number of users to return"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"Number of users to skip"`
}

type UserOutput struct {
	Body ports.UserDTO
}

type UsersOutput struct {
	Body []ports.UserDTO
}

type UserHandler struct {
	repo ports.UserRepository
}

func NewUserHandler(repo ports.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-core-user",
		Method:        http.MethodPost,
		Path:          "/api/v1/core/users/",
		Summary:       "Create a user",
		Description:   "Creates a user and stores a bcrypt hash of the supplied password.",
		Tags:          []string{"Core users"},
		DefaultStatus: http.StatusCreated,
		MaxBodyBytes:  maxCreateUserBodyBytes,
		Errors:        []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, h.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "list-core-users",
		Method:      http.MethodGet,
		Path:        "/api/v1/core/users/",
		Summary:     "List users",
		Description: "Returns Core users using offset pagination.",
		Tags:        []string{"Core users"},
		Errors:      []int{http.StatusInternalServerError},
	}, h.ListUsers)

	huma.Register(api, huma.Operation{
		OperationID: "get-core-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/core/users/{id}",
		Summary:     "Get a user",
		Description: "Returns a Core user for the supplied UUID.",
		Tags:        []string{"Core users"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, h.GetUserByID)
}

func (h *UserHandler) CreateUser(ctx context.Context, input *CreateUserInput) (*UserOutput, error) {
	input.Body.Email = strings.TrimSpace(input.Body.Email)
	if input.Body.Email == "" {
		return nil, huma.Error400BadRequest("Email is required")
	}
	if input.Body.Password == "" {
		return nil, huma.Error400BadRequest("Password is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	user, err := h.repo.CreateUser(ctx, ports.CreateUserParams{
		Email:        input.Body.Email,
		PasswordHash: string(hash),
		FirstName:    input.Body.FirstName,
		LastName:     input.Body.LastName,
		Phone:        input.Body.Phone,
		Role:         input.Body.Role,
	})
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &UserOutput{Body: *user}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, input *GetUserInput) (*UserOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid UUID format")
	}

	user, err := h.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		}
		slog.Error("Failed to get user", "error", err, "user_id", id)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &UserOutput{Body: *user}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, input *ListUsersInput) (*UsersOutput, error) {
	users, err := h.repo.ListUsers(ctx, int32(input.Limit), int32(input.Offset))
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &UsersOutput{Body: users}, nil
}
