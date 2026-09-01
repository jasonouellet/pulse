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

type TeamRepository struct {
	pool database.PgxPool
	q    *db.Queries
}

func NewTeamRepository(pool database.PgxPool) *TeamRepository {
	return &TeamRepository{
		pool: pool,
		q:    db.New(pool),
	}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, p ports.CreateTeamParams) (*ports.TeamDTO, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.q.WithTx(tx)

	teamRow, err := qtx.CreateTeam(ctx, db.CreateTeamParams{
		ClubID:     p.ClubID,
		SportID:    p.SportID,
		WindowID:   p.WindowID,
		Type:       db.CoreTeamType(p.Type),
		Name:       p.Name,
		SeasonYear: p.SeasonYear,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	for _, poolID := range p.PoolIDs {
		err := qtx.LinkTeamToPool(ctx, db.LinkTeamToPoolParams{
			TeamID: teamRow.ID,
			PoolID: poolID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to link pool %s to team: %w", poolID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &ports.TeamDTO{
		ID:         teamRow.ID,
		ClubID:     teamRow.ClubID,
		SportID:    teamRow.SportID,
		WindowID:   teamRow.WindowID,
		Type:       string(teamRow.Type),
		Name:       teamRow.Name,
		SeasonYear: teamRow.SeasonYear,
		PoolIDs:    p.PoolIDs,
	}, nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id uuid.UUID) (*ports.TeamDTO, error) {
	teamRow, err := r.q.GetTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team by id: %w", err)
	}

	poolIDs, err := r.q.ListTeamPoolIDs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list team pool ids: %w", err)
	}

	return &ports.TeamDTO{
		ID:         teamRow.ID,
		ClubID:     teamRow.ClubID,
		SportID:    teamRow.SportID,
		WindowID:   teamRow.WindowID,
		Type:       string(teamRow.Type),
		Name:       teamRow.Name,
		SeasonYear: teamRow.SeasonYear,
		PoolIDs:    poolIDs,
	}, nil
}

func (r *TeamRepository) ListTeamsByClub(ctx context.Context, clubID uuid.UUID, limit, offset int32) ([]ports.TeamDTO, error) {
	rows, err := r.q.ListTeamsByClub(ctx, db.ListTeamsByClubParams{
		ClubID: clubID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list teams by club: %w", err)
	}

	teams := make([]ports.TeamDTO, 0, len(rows))
	for _, row := range rows {
		poolIDs, err := r.q.ListTeamPoolIDs(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list team pool ids: %w", err)
		}

		teams = append(teams, ports.TeamDTO{
			ID:         row.ID,
			ClubID:     row.ClubID,
			SportID:    row.SportID,
			WindowID:   row.WindowID,
			Type:       string(row.Type),
			Name:       row.Name,
			SeasonYear: row.SeasonYear,
			PoolIDs:    poolIDs,
		})
	}

	return teams, nil
}

func (r *TeamRepository) AddPlayerToTeam(ctx context.Context, teamID, playerID uuid.UUID) error {
	err := r.q.AddPlayerToTeam(ctx, db.AddPlayerToTeamParams{
		TeamID:   teamID,
		PlayerID: playerID,
	})
	if err != nil {
		return fmt.Errorf("failed to add player to team: %w", err)
	}
	return nil
}

func (r *TeamRepository) RemovePlayerFromTeam(ctx context.Context, teamID, playerID uuid.UUID) error {
	err := r.q.RemovePlayerFromTeam(ctx, db.RemovePlayerFromTeamParams{
		TeamID:   teamID,
		PlayerID: playerID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove player from team: %w", err)
	}
	return nil
}

func (r *TeamRepository) ListTeamPlayerIDs(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	playerIDs, err := r.q.ListTeamPlayerIDs(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team player ids: %w", err)
	}
	return playerIDs, nil
}
