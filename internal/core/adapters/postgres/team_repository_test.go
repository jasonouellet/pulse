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

func newMockTeamRepo(t *testing.T) (*TeamRepository, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(mock.Close)
	return NewTeamRepository(mock), mock
}

func TestTeamRepository_CreateTeam_WithPools(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year"}).
		AddRow(teamID, clubID, sportID, windowID, "TRAINING_GROUP", "U10-M Division 1", int32(2026))

	// A committed transaction: pgx.Tx.Rollback() is still called via defer
	// after a successful Commit — pgxmock requires an explicit
	// ExpectRollback() for that trailing no-op call, or
	// ExpectationsWereMet() fails.
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO core.teams").
		WithArgs(clubID, sportID, windowID, "TRAINING_GROUP", "U10-M Division 1", int32(2026)).
		WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO core.team_pools").
		WithArgs(teamID, poolID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	mock.ExpectRollback()

	got, err := repo.CreateTeam(context.Background(), ports.CreateTeamParams{
		ClubID:     clubID,
		SportID:    sportID,
		WindowID:   windowID,
		Type:       "TRAINING_GROUP",
		Name:       "U10-M Division 1",
		SeasonYear: 2026,
		PoolIDs:    []uuid.UUID{poolID},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != teamID {
		t.Errorf("expected id %s, got %s", teamID, got.ID)
	}
	if len(got.PoolIDs) != 1 || got.PoolIDs[0] != poolID {
		t.Errorf("expected pool ids [%s], got %v", poolID, got.PoolIDs)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_CreateTeam_PoolLinkFails_Rollback(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year"}).
		AddRow(teamID, clubID, sportID, windowID, "SEASON_TEAM", "Faucons U9-U10", int32(2026))

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO core.teams").
		WithArgs(clubID, sportID, windowID, "SEASON_TEAM", "Faucons U9-U10", int32(2026)).
		WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO core.team_pools").
		WithArgs(teamID, poolID).
		WillReturnError(errors.New("foreign_key_violation: pool_id does not exist"))
	mock.ExpectRollback()

	_, err := repo.CreateTeam(context.Background(), ports.CreateTeamParams{
		ClubID:     clubID,
		SportID:    sportID,
		WindowID:   windowID,
		Type:       "SEASON_TEAM",
		Name:       "Faucons U9-U10",
		SeasonYear: 2026,
		PoolIDs:    []uuid.UUID{poolID},
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_GetTeamByID_Found(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID1, poolID2 := uuid.New(), uuid.New()

	teamRows := pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year"}).
		AddRow(teamID, clubID, sportID, windowID, "SEASON_TEAM", "Faucons U9-U10", int32(2026))
	poolRows := pgxmock.NewRows([]string{"pool_id"}).
		AddRow(poolID1).
		AddRow(poolID2)

	mock.ExpectQuery("SELECT id, club_id, sport_id, window_id, type").
		WithArgs(teamID).
		WillReturnRows(teamRows)
	mock.ExpectQuery("SELECT pool_id FROM core.team_pools").
		WithArgs(teamID).
		WillReturnRows(poolRows)

	got, err := repo.GetTeamByID(context.Background(), teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.PoolIDs) != 2 {
		t.Errorf("expected 2 pool ids, got %d", len(got.PoolIDs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_GetTeamByID_NotFound(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID := uuid.New()

	mock.ExpectQuery("SELECT id, club_id, sport_id, window_id, type").
		WithArgs(teamID).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetTeamByID(context.Background(), teamID)
	if !errors.Is(err, ports.ErrTeamNotFound) {
		t.Errorf("expected ErrTeamNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_ListTeams(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	clubID := uuid.New()
	team1, team2 := uuid.New(), uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()

	teamRows := pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year"}).
		AddRow(team1, clubID, sportID, windowID, "TRAINING_GROUP", "U10-M Division 1", int32(2026)).
		AddRow(team2, clubID, sportID, windowID, "SEASON_TEAM", "Faucons U9-U10", int32(2026))

	mock.ExpectQuery("SELECT id, club_id, sport_id, window_id, type").
		WithArgs(clubID, int32(20), int32(0)).
		WillReturnRows(teamRows)

	// GetTeamByID's listPoolIDs helper runs once per team returned.
	mock.ExpectQuery("SELECT pool_id FROM core.team_pools").
		WithArgs(team1).
		WillReturnRows(pgxmock.NewRows([]string{"pool_id"}))
	mock.ExpectQuery("SELECT pool_id FROM core.team_pools").
		WithArgs(team2).
		WillReturnRows(pgxmock.NewRows([]string{"pool_id"}))

	got, err := repo.ListTeams(context.Background(), clubID, 20, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_AddPlayer(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID, playerID := uuid.New(), uuid.New()

	mock.ExpectExec("INSERT INTO core.team_players").
		WithArgs(teamID, playerID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err := repo.AddPlayer(context.Background(), teamID, playerID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_RemovePlayer(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID, playerID := uuid.New(), uuid.New()

	mock.ExpectExec("DELETE FROM core.team_players").
		WithArgs(teamID, playerID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := repo.RemovePlayer(context.Background(), teamID, playerID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestTeamRepository_ListTeamPlayerIDs(t *testing.T) {
	repo, mock := newMockTeamRepo(t)
	teamID := uuid.New()
	p1, p2 := uuid.New(), uuid.New()

	rows := pgxmock.NewRows([]string{"player_id"}).
		AddRow(p1).
		AddRow(p2)

	mock.ExpectQuery("SELECT player_id FROM core.team_players").
		WithArgs(teamID).
		WillReturnRows(rows)

	got, err := repo.ListTeamPlayerIDs(context.Background(), teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 player ids, got %d", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
