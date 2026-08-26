package ports

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrUserNotFound is returned when a requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

// UserDTO represents the user entity transferred across boundaries.
type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     *string   `json:"phone,omitempty"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
}

// CreateUserParams holds data needed to create a new user.
type CreateUserParams struct {
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Phone        *string `json:"phone,omitempty"`
	Role         string  `json:"role"`
}

// UserRepository defines the persistence port for the Core module.
type UserRepository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*UserDTO, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDTO, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]UserDTO, error)
}
