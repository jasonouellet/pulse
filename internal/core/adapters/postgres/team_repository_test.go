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

func TestTeamRepository_CreateTeam_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?i)INSERT INTO core\.teams`).
		WithArgs(clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year", "created_at", "updated_at"}).
				AddRow(teamID, clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026), now, now),
		)
	mock.ExpectExec(`(?i)INSERT INTO core\.team_pools`).
		WithArgs(teamID, poolID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	team, err := repo.CreateTeam(context.Background(), ports.CreateTeamParams{
		ClubID:     clubID,
		SportID:    sportID,
		WindowID:   windowID,
		Type:       "SEASON_TEAM",
		Name:       "U10F D2",
		SeasonYear: 2026,
		PoolIDs:    []uuid.UUID{poolID},
	})

	require.NoError(t, err)
	assert.Equal(t, teamID, team.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_CreateTeam_BeginTxError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	mock.ExpectBegin().WillReturnError(errors.New("tx error"))

	_, err = repo.CreateTeam(context.Background(), ports.CreateTeamParams{})
	assert.ErrorContains(t, err, "failed to begin tx")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetTeamByID_Found(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at FROM core\.teams`).
		WithArgs(teamID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year", "created_at", "updated_at"}).
				AddRow(teamID, clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026), now, now),
		)

	mock.ExpectQuery(`(?i)SELECT pool_id FROM core\.team_pools`).
		WithArgs(teamID).
		WillReturnRows(
			pgxmock.NewRows([]string{"pool_id"}).AddRow(poolID),
		)

	team, err := repo.GetTeamByID(context.Background(), teamID)
	require.NoError(t, err)
	assert.Equal(t, teamID, team.ID)
	assert.Len(t, team.PoolIDs, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetTeamByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()

	mock.ExpectQuery(`(?i)SELECT id, club_id`).
		WithArgs(teamID).
		WillReturnError(pgx.ErrNoRows)

	_, err = repo.GetTeamByID(context.Background(), teamID)
	assert.ErrorIs(t, err, ports.ErrTeamNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_AddPlayerToTeam(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	playerID := uuid.New()

	mock.ExpectExec(`(?i)INSERT INTO core\.team_players`).
		WithArgs(teamID, playerID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.AddPlayerToTeam(context.Background(), teamID, playerID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_RemovePlayerFromTeam(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	playerID := uuid.New()

	mock.ExpectExec(`(?i)DELETE FROM core\.team_players`).
		WithArgs(teamID, playerID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.RemovePlayerFromTeam(context.Background(), teamID, playerID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_ListTeamPlayerIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	playerID := uuid.New()

	mock.ExpectQuery(`(?i)SELECT player_id FROM core\.team_players`).
		WithArgs(teamID).
		WillReturnRows(
			pgxmock.NewRows([]string{"player_id"}).AddRow(playerID),
		)

	ids, err := repo.ListTeamPlayerIDs(context.Background(), teamID)
	require.NoError(t, err)
	assert.Len(t, ids, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestTeamRepository_CreateTeam_LinkPoolError_Rollback(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?i)INSERT INTO core\.teams`).
		WithArgs(clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year", "created_at", "updated_at"}).
				AddRow(teamID, clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026), now, now),
		)
	mock.ExpectExec(`(?i)INSERT INTO core\.team_pools`).
		WithArgs(teamID, poolID).
		WillReturnError(errors.New("db error link pool"))
	mock.ExpectRollback()

	_, err = repo.CreateTeam(context.Background(), ports.CreateTeamParams{
		ClubID:     clubID,
		SportID:    sportID,
		WindowID:   windowID,
		Type:       "SEASON_TEAM",
		Name:       "U10F D2",
		SeasonYear: 2026,
		PoolIDs:    []uuid.UUID{poolID},
	})

	assert.ErrorContains(t, err, "failed to link pool")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_CreateTeam_CommitError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`(?i)INSERT INTO core\.teams`).
		WithArgs(clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year", "created_at", "updated_at"}).
				AddRow(teamID, clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026), now, now),
		)
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	_, err = repo.CreateTeam(context.Background(), ports.CreateTeamParams{
		ClubID:     clubID,
		SportID:    sportID,
		WindowID:   windowID,
		Type:       "SEASON_TEAM",
		Name:       "U10F D2",
		SeasonYear: 2026,
	})

	assert.ErrorContains(t, err, "failed to commit tx")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetTeamByID_ListPoolsError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at FROM core\.teams`).
		WithArgs(teamID).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year", "created_at", "updated_at"}).
				AddRow(teamID, clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026), now, now),
		)

	mock.ExpectQuery(`(?i)SELECT pool_id FROM core\.team_pools`).
		WithArgs(teamID).
		WillReturnError(errors.New("pool fetch error"))

	_, err = repo.GetTeamByID(context.Background(), teamID)
	assert.ErrorContains(t, err, "failed to list team pool ids")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_ListTeamsByClub_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	clubID := uuid.New()
	sportID := uuid.New()
	windowID := uuid.New()
	poolID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`(?i)SELECT id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at FROM core\.teams`).
		WithArgs(clubID, int32(10), int32(0)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "club_id", "sport_id", "window_id", "type", "name", "season_year", "created_at", "updated_at"}).
				AddRow(teamID, clubID, sportID, windowID, db.CoreTeamType("SEASON_TEAM"), "U10F D2", int32(2026), now, now),
		)

	mock.ExpectQuery(`(?i)SELECT pool_id FROM core\.team_pools`).
		WithArgs(teamID).
		WillReturnRows(
			pgxmock.NewRows([]string{"pool_id"}).AddRow(poolID),
		)

	teams, err := repo.ListTeamsByClub(context.Background(), clubID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, teams, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_ListTeamsByClub_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	clubID := uuid.New()

	mock.ExpectQuery(`(?i)SELECT id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at FROM core\.teams`).
		WithArgs(clubID, int32(10), int32(0)).
		WillReturnError(errors.New("db error list"))

	_, err = repo.ListTeamsByClub(context.Background(), clubID, 10, 0)
	assert.ErrorContains(t, err, "failed to list teams by club")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_AddPlayerToTeam_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	playerID := uuid.New()

	mock.ExpectExec(`(?i)INSERT INTO core\.team_players`).
		WithArgs(teamID, playerID).
		WillReturnError(errors.New("db insert error"))

	err = repo.AddPlayerToTeam(context.Background(), teamID, playerID)
	assert.ErrorContains(t, err, "failed to add player to team")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_RemovePlayerFromTeam_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()
	playerID := uuid.New()

	mock.ExpectExec(`(?i)DELETE FROM core\.team_players`).
		WithArgs(teamID, playerID).
		WillReturnError(errors.New("db delete error"))

	err = repo.RemovePlayerFromTeam(context.Background(), teamID, playerID)
	assert.ErrorContains(t, err, "failed to remove player from team")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_ListTeamPlayerIDs_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewTeamRepository(mock)
	teamID := uuid.New()

	mock.ExpectQuery(`(?i)SELECT player_id FROM core\.team_players`).
		WithArgs(teamID).
		WillReturnError(errors.New("db query error"))

	_, err = repo.ListTeamPlayerIDs(context.Background(), teamID)
	assert.ErrorContains(t, err, "failed to list team player ids")
	assert.NoError(t, mock.ExpectationsWereMet())
}
