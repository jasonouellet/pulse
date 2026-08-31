package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PgxPool is the subset of *pgxpool.Pool's methods used by repositories.
// Depending on this interface instead of the concrete *pgxpool.Pool type
// is what makes repositories testable with pgxmock — no real database
// needed in unit tests. *pgxpool.Pool satisfies this interface already,
// so production wiring (main.go) is unaffected.
type PgxPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
