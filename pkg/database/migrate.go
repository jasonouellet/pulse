package database

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// schemaMigration pairs a migrations directory with the name of the
// tracking table golang-migrate should use for it. Each schema needs its
// own tracking table — otherwise golang-migrate shares a single
// schema_migrations table across all of them and, after the first schema
// is migrated, believes every other schema's migration "1" is already
// applied, silently skipping it.
type schemaMigration struct {
	path            string
	migrationsTable string
}

// RunAllMigrations applies every schema's migrations in dependency order.
// scheduling and core are shared kernels (ADR-003 §2.B, ADR-008 §4) that
// tournament depends on via real foreign keys, so they must be migrated
// first, in this exact order.
func RunAllMigrations(dbURL string) error {
	schemas := []schemaMigration{
		{path: "migrations/scheduling", migrationsTable: "scheduling_schema_migrations"},
		{path: "migrations/core", migrationsTable: "core_schema_migrations"},
		{path: "migrations/tournament", migrationsTable: "tournament_schema_migrations"},
	}

	for _, s := range schemas {
		if err := RunMigrations(s.path, dbURL, s.migrationsTable); err != nil {
			return fmt.Errorf("schema %q: %w", s.path, err)
		}
	}
	return nil
}

// RunMigrations executes all .up.sql files for a given module/schema,
// tracked in its own migrationsTable so multiple schemas can be migrated
// against the same database without colliding on migration state.
func RunMigrations(migrationPath string, dbURL string, migrationsTable string) error {
	slog.Info("Running database migrations...", "path", migrationPath, "table", migrationsTable)

	scopedURL, err := withMigrationsTable(dbURL, migrationsTable)
	if err != nil {
		return fmt.Errorf("failed to build scoped migration URL: %w", err)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationPath),
		scopedURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("Database migrations applied successfully!", "path", migrationPath)
	return nil
}

// withMigrationsTable adds/overrides the x-migrations-table query parameter
// on dbURL, using net/url so it composes correctly regardless of which
// other query parameters (sslmode, etc.) are already present.
func withMigrationsTable(dbURL string, table string) (string, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return "", fmt.Errorf("invalid database URL: %w", err)
	}
	q := u.Query()
	q.Set("x-migrations-table", table)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
