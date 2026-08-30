package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pulse/internal/core/ports"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, p ports.CreateTeamParams) (*ports.TeamDTO, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var t ports.TeamDTO
	query := `
		INSERT INTO core.teams (club_id, sport_id, window_id, type, name, season_year)
		VALUES ($1, $2, $3, $4::core.team_type, $5, $6)
		RETURNING id, club_id, sport_id, window_id, type::text, name, season_year;
	`
	err = tx.QueryRow(ctx, query, p.ClubID, p.SportID, p.WindowID, p.Type, p.Name, p.SeasonYear).Scan(
		&t.ID, &t.ClubID, &t.SportID, &t.WindowID, &t.Type, &t.Name, &t.SeasonYear,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	for _, poolID := range p.PoolIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO core.team_pools (team_id, pool_id) VALUES ($1, $2);`, t.ID, poolID); err != nil {
			return nil, fmt.Errorf("failed to link pool %s: %w", poolID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit team creation: %w", err)
	}

	t.PoolIDs = p.PoolIDs
	return &t, nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id uuid.UUID) (*ports.TeamDTO, error) {
	var t ports.TeamDTO
	query := `
		SELECT id, club_id, sport_id, window_id, type::text, name, season_year
		FROM core.teams
		WHERE id = $1;
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.ClubID, &t.SportID, &t.WindowID, &t.Type, &t.Name, &t.SeasonYear,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team by id: %w", err)
	}

	poolIDs, err := r.listPoolIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	t.PoolIDs = poolIDs
	return &t, nil
}

func (r *TeamRepository) ListTeams(ctx context.Context, clubID uuid.UUID, limit, offset int32) ([]ports.TeamDTO, error) {
	query := `
		SELECT id, club_id, sport_id, window_id, type::text, name, season_year
		FROM core.teams
		WHERE club_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	rows, err := r.db.Query(ctx, query, clubID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	defer rows.Close()

	var teams []ports.TeamDTO
	for rows.Next() {
		var t ports.TeamDTO
		if err := rows.Scan(&t.ID, &t.ClubID, &t.SportID, &t.WindowID, &t.Type, &t.Name, &t.SeasonYear); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}

	for i := range teams {
		poolIDs, err := r.listPoolIDs(ctx, teams[i].ID)
		if err != nil {
			return nil, err
		}
		teams[i].PoolIDs = poolIDs
	}

	return teams, nil
}

func (r *TeamRepository) AddPlayer(ctx context.Context, teamID, playerID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO core.team_players (team_id, player_id)
		VALUES ($1, $2)
		ON CONFLICT (team_id, player_id) DO NOTHING;
	`, teamID, playerID)
	if err != nil {
		return fmt.Errorf("failed to add player to team: %w", err)
	}
	return nil
}

func (r *TeamRepository) RemovePlayer(ctx context.Context, teamID, playerID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM core.team_players WHERE team_id = $1 AND player_id = $2;
	`, teamID, playerID)
	if err != nil {
		return fmt.Errorf("failed to remove player from team: %w", err)
	}
	return nil
}

func (r *TeamRepository) ListTeamPlayerIDs(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT player_id FROM core.team_players WHERE team_id = $1;`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team players: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *TeamRepository) listPoolIDs(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT pool_id FROM core.team_pools WHERE team_id = $1;`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team pools: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}