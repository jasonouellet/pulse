package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"pulse/internal/core/ports"
	"pulse/pkg/database"
)

type UserRepository struct {
	db database.PgxPool
}

func NewUserRepository(db database.PgxPool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, p ports.CreateUserParams) (*ports.UserDTO, error) {
	query := `
		INSERT INTO core.users (email, password_hash, first_name, last_name, phone, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, first_name, last_name, phone, role, is_active;
	`
	var u ports.UserDTO
	err := r.db.QueryRow(ctx, query, p.Email, p.PasswordHash, p.FirstName, p.LastName, p.Phone, p.Role).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Role, &u.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*ports.UserDTO, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, role, is_active
		FROM core.users
		WHERE id = $1;
	`
	var u ports.UserDTO
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Role, &u.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*ports.UserDTO, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, role, is_active
		FROM core.users
		WHERE email = $1;
	`
	var u ports.UserDTO
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Role, &u.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int32) ([]ports.UserDTO, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, role, is_active
		FROM core.users
		ORDER BY last_name ASC, first_name ASC
		LIMIT $1 OFFSET $2;
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []ports.UserDTO
	for rows.Next() {
		var u ports.UserDTO
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Role, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
