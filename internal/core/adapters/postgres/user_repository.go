package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"pulse/internal/core/adapters/postgres/db"
	"pulse/internal/core/ports"
	"pulse/pkg/database"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(pool database.PgxPool) *UserRepository {
	return &UserRepository{
		q: db.New(pool),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, p ports.CreateUserParams) (*ports.UserDTO, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Phone:        p.Phone,
		Role:         db.CoreUserRole(p.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &ports.UserDTO{
		ID:        row.ID,
		Email:     row.Email,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Phone:     row.Phone,
		Role:      string(row.Role),
		IsActive:  row.IsActive,
	}, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*ports.UserDTO, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &ports.UserDTO{
		ID:        row.ID,
		Email:     row.Email,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Phone:     row.Phone,
		Role:      string(row.Role),
		IsActive:  row.IsActive,
	}, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*ports.UserDTO, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &ports.UserDTO{
		ID:        row.ID,
		Email:     row.Email,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Phone:     row.Phone,
		Role:      string(row.Role),
		IsActive:  row.IsActive,
	}, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int32) ([]ports.UserDTO, error) {
	rows, err := r.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]ports.UserDTO, 0, len(rows))
	for _, row := range rows {
		users = append(users, ports.UserDTO{
			ID:        row.ID,
			Email:     row.Email,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Phone:     row.Phone,
			Role:      string(row.Role),
			IsActive:  row.IsActive,
		})
	}
	return users, nil
}
