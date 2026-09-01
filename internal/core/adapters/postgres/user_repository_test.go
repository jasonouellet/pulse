package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pulse/internal/core/adapters/postgres/db"
	"pulse/internal/core/ports"
)

func TestUserRepository_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	userID := uuid.New()
	phone := "555-0199"
	now := time.Now()

	mock.ExpectQuery(`(?i)INSERT INTO core\.users`).
		WithArgs("jean.dupont@pulse.local", "hashed-secret", "Jean", "Dupont", &phone, db.CoreUserRole("COACH")).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active", "last_login_at", "created_at", "updated_at"}).
				AddRow(userID, "jean.dupont@pulse.local", "Jean", "Dupont", &phone, db.CoreUserRole("COACH"), true, nil, now, now),
		)

	user, err := repo.CreateUser(context.Background(), ports.CreateUserParams{
		Email:        "jean.dupont@pulse.local",
		PasswordHash: "hashed-secret",
		FirstName:    "Jean",
		LastName:     "Dupont",
		Phone:        &phone,
		Role:         "COACH",
	})

	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "COACH", user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_DBError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	phone := "555-0199"

	mock.ExpectQuery(`(?i)INSERT INTO core\.users`).
		WithArgs("dup@pulse.local", "hashed-secret", "Jean", "Dupont", &phone, db.CoreUserRole("COACH")).
		WillReturnError(errors.New("db error"))

	_, err = repo.CreateUser(context.Background(), ports.CreateUserParams{
		Email:        "dup@pulse.local",
		PasswordHash: "hashed-secret",
		FirstName:    "Jean",
		LastName:     "Dupont",
		Phone:        &phone,
		Role:         "COACH",
	})

	assert.ErrorContains(t, err, "failed to create user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_Found(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	userID := uuid.New()
	phone := "555-0199"
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at FROM core\.users`).
		WithArgs(userID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "password_hash", "first_name", "last_name", "phone", "role", "is_active", "last_login_at", "created_at", "updated_at"}).
				AddRow(userID, "jean.dupont@pulse.local", "hashed", "Jean", "Dupont", &phone, db.CoreUserRole("COACH"), true, nil, now, now),
		)

	user, err := repo.GetUserByID(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	userID := uuid.New()

	mock.ExpectQuery(`(?i)SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at FROM core\.users`).
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	_, err = repo.GetUserByID(context.Background(), userID)
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_DBError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	userID := uuid.New()

	mock.ExpectQuery(`(?i)SELECT id, email, password_hash`).
		WithArgs(userID).
		WillReturnError(errors.New("db connection lost"))

	_, err = repo.GetUserByID(context.Background(), userID)
	assert.ErrorContains(t, err, "failed to get user by id")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail_Found(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	userID := uuid.New()
	phone := "555-0199"
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at FROM core\.users`).
		WithArgs("jean.dupont@pulse.local").
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "password_hash", "first_name", "last_name", "phone", "role", "is_active", "last_login_at", "created_at", "updated_at"}).
				AddRow(userID, "jean.dupont@pulse.local", "hashed", "Jean", "Dupont", &phone, db.CoreUserRole("COACH"), true, nil, now, now),
		)

	user, err := repo.GetUserByEmail(context.Background(), "jean.dupont@pulse.local")
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	mock.ExpectQuery(`(?i)SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at FROM core\.users`).
		WithArgs("ghost@pulse.local").
		WillReturnError(pgx.ErrNoRows)

	_, err = repo.GetUserByEmail(context.Background(), "ghost@pulse.local")
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail_DBError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	mock.ExpectQuery(`(?i)SELECT id, email, password_hash`).
		WithArgs("test@pulse.local").
		WillReturnError(errors.New("db connection lost"))

	_, err = repo.GetUserByEmail(context.Background(), "test@pulse.local")
	assert.ErrorContains(t, err, "failed to get user by email")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ListUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)
	userID := uuid.New()
	phone := "555-0199"
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id, email, first_name, last_name, phone, role, is_active, created_at FROM core\.users`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active", "created_at"}).
				AddRow(userID, "jean.dupont@pulse.local", "Jean", "Dupont", &phone, db.CoreUserRole("COACH"), true, now),
		)

	users, err := repo.ListUsers(context.Background(), 10, 0)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ListUsers_DBError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	mock.ExpectQuery(`(?i)SELECT id, email, first_name`).
		WithArgs(int32(10), int32(0)).
		WillReturnError(errors.New("db timeout"))

	_, err = repo.ListUsers(context.Background(), 10, 0)
	assert.ErrorContains(t, err, "failed to list users")
	assert.NoError(t, mock.ExpectationsWereMet())
}
