package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"

	"pulse/internal/core/ports"
)

func newMockUserRepo(t *testing.T) (*UserRepository, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(mock.Close)
	return NewUserRepository(mock), mock
}

func TestUserRepository_CreateUser(t *testing.T) {
	repo, mock := newMockUserRepo(t)
	id := uuid.New()
	phone := "418-555-0192"

	rows := pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active"}).
		AddRow(id, "coach@pulse.local", "Jean", "Dupont", &phone, "COACH", true)

	mock.ExpectQuery("INSERT INTO core.users").
		WithArgs("coach@pulse.local", "hashed-secret", "Jean", "Dupont", &phone, "COACH").
		WillReturnRows(rows)

	got, err := repo.CreateUser(context.Background(), ports.CreateUserParams{
		Email:        "coach@pulse.local",
		PasswordHash: "hashed-secret",
		FirstName:    "Jean",
		LastName:     "Dupont",
		Phone:        &phone,
		Role:         "COACH",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected id %s, got %s", id, got.ID)
	}
	if got.Email != "coach@pulse.local" {
		t.Errorf("expected email %q, got %q", "coach@pulse.local", got.Email)
	}
	if !got.IsActive {
		t.Error("expected new user to be active")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_CreateUser_DBError(t *testing.T) {
	repo, mock := newMockUserRepo(t)
	phone := "418-555-0192"

	mock.ExpectQuery("INSERT INTO core.users").
		WithArgs("dup@pulse.local", "hashed-secret", "Jean", "Dupont", &phone, "COACH").
		WillReturnError(errors.New("unique_violation: email already exists"))

	_, err := repo.CreateUser(context.Background(), ports.CreateUserParams{
		Email:        "dup@pulse.local",
		PasswordHash: "hashed-secret",
		FirstName:    "Jean",
		LastName:     "Dupont",
		Phone:        &phone,
		Role:         "COACH",
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_GetUserByID_Found(t *testing.T) {
	repo, mock := newMockUserRepo(t)
	id := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active"}).
		AddRow(id, "leo@pulse.local", "Léo", "Tremblay", nil, "GUARDIAN", true)

	mock.ExpectQuery("SELECT id, email, first_name, last_name, phone, role").
		WithArgs(id).
		WillReturnRows(rows)

	got, err := repo.GetUserByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.FirstName != "Léo" {
		t.Errorf("expected first name %q, got %q", "Léo", got.FirstName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	repo, mock := newMockUserRepo(t)
	id := uuid.New()

	mock.ExpectQuery("SELECT id, email, first_name, last_name, phone, role").
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetUserByID(context.Background(), id)
	if !errors.Is(err, ports.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_GetUserByEmail_Found(t *testing.T) {
	repo, mock := newMockUserRepo(t)
	id := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active"}).
		AddRow(id, "zoe@pulse.local", "Zoé", "Tremblay", nil, "PLAYER", true)

	mock.ExpectQuery("SELECT id, email, first_name, last_name, phone, role").
		WithArgs("zoe@pulse.local").
		WillReturnRows(rows)

	got, err := repo.GetUserByEmail(context.Background(), "zoe@pulse.local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected id %s, got %s", id, got.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	repo, mock := newMockUserRepo(t)

	mock.ExpectQuery("SELECT id, email, first_name, last_name, phone, role").
		WithArgs("ghost@pulse.local").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetUserByEmail(context.Background(), "ghost@pulse.local")
	if !errors.Is(err, ports.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_ListUsers(t *testing.T) {
	repo, mock := newMockUserRepo(t)
	id1, id2 := uuid.New(), uuid.New()

	rows := pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active"}).
		AddRow(id1, "a@pulse.local", "A", "A", nil, "GUARDIAN", true).
		AddRow(id2, "b@pulse.local", "B", "B", nil, "COACH", true)

	mock.ExpectQuery("SELECT id, email, first_name, last_name, phone, role").
		WithArgs(int32(20), int32(0)).
		WillReturnRows(rows)

	got, err := repo.ListUsers(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserRepository_ListUsers_Empty(t *testing.T) {
	repo, mock := newMockUserRepo(t)

	rows := pgxmock.NewRows([]string{"id", "email", "first_name", "last_name", "phone", "role", "is_active"})

	mock.ExpectQuery("SELECT id, email, first_name, last_name, phone, role").
		WithArgs(int32(20), int32(0)).
		WillReturnRows(rows)

	got, err := repo.ListUsers(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 users, got %d", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
