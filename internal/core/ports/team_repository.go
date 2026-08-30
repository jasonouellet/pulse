package ports

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrTeamNotFound is returned when a requested team does not exist.
var ErrTeamNotFound = errors.New("team not found")

// TeamDTO represents a persistent club team (training group or season
// team — ADR-008 §1) transferred across boundaries.
type TeamDTO struct {
	ID         uuid.UUID   `json:"id"`
	ClubID     uuid.UUID   `json:"club_id"`
	SportID    uuid.UUID   `json:"sport_id"`
	WindowID   uuid.UUID   `json:"window_id"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	SeasonYear int32       `json:"season_year"`
	PoolIDs    []uuid.UUID `json:"pool_ids"`
}

// CreateTeamParams holds data needed to create a new team. PoolIDs is one
// entry for a TRAINING_GROUP (1:1 with its pool) or several for a
// SEASON_TEAM combining age groups (e.g., U9 + U10).
type CreateTeamParams struct {
	ClubID     uuid.UUID
	SportID    uuid.UUID
	WindowID   uuid.UUID
	Type       string
	Name       string
	SeasonYear int32
	PoolIDs    []uuid.UUID
}

// TeamRepository defines the persistence port for core.teams and its
// player membership (core.team_players).
type TeamRepository interface {
	CreateTeam(ctx context.Context, params CreateTeamParams) (*TeamDTO, error)
	GetTeamByID(ctx context.Context, id uuid.UUID) (*TeamDTO, error)
	ListTeams(ctx context.Context, clubID uuid.UUID, limit, offset int32) ([]TeamDTO, error)
	AddPlayer(ctx context.Context, teamID, playerID uuid.UUID) error
	RemovePlayer(ctx context.Context, teamID, playerID uuid.UUID) error
	ListTeamPlayerIDs(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error)
}