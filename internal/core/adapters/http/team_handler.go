package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"pulse/internal/core/ports"
)

type CreateTeamRequest struct {
	ClubID     string   `json:"club_id" format:"uuid"`
	SportID    string   `json:"sport_id" format:"uuid"`
	WindowID   string   `json:"window_id" format:"uuid" doc:"scheduling.date_windows UUID for this team's season range"`
	Type       string   `json:"type" enum:"TRAINING_GROUP,SEASON_TEAM"`
	Name       string   `json:"name"`
	SeasonYear int32    `json:"season_year"`
	PoolIDs    []string `json:"pool_ids" doc:"One pool for TRAINING_GROUP; one or more for SEASON_TEAM combining age groups"`
}

type CreateTeamInput struct {
	Body CreateTeamRequest
}

type GetTeamInput struct {
	ID string `path:"id" format:"uuid" doc:"Team UUID"`
}

type ListTeamsByClubInput struct {
	ClubID string `query:"club_id" format:"uuid" doc:"Club UUID to list teams for"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset int    `query:"offset" default:"0" minimum:"0"`
}

type AddPlayerRequest struct {
	PlayerID string `json:"player_id" format:"uuid"`
}

type AddPlayerInput struct {
	ID   string `path:"id" format:"uuid"`
	Body AddPlayerRequest
}

type RemovePlayerInput struct {
	ID       string `path:"id" format:"uuid"`
	PlayerID string `path:"playerId" format:"uuid"`
}

type ListPlayersInput struct {
	ID string `path:"id" format:"uuid"`
}

type TeamOutput struct {
	Body ports.TeamDTO
}

type TeamsOutput struct {
	Body []ports.TeamDTO
}

type PlayerIDsOutput struct {
	Body []uuid.UUID
}

// EmptyOutput is used for endpoints that return no body (204 No Content).
type EmptyOutput struct{}

type TeamHandler struct {
	repo ports.TeamRepository
}

func NewTeamHandler(repo ports.TeamRepository) *TeamHandler {
	return &TeamHandler{repo: repo}
}

func (h *TeamHandler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-core-team",
		Method:        http.MethodPost,
		Path:          "/api/v1/core/teams/",
		Summary:       "Create a team",
		Description:   "Creates a persistent training group or season team (ADR-008 §1) — independent of any tournament/event.",
		Tags:          []string{"Core teams"},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, h.CreateTeam)

	huma.Register(api, huma.Operation{
		OperationID: "get-core-team",
		Method:      http.MethodGet,
		Path:        "/api/v1/core/teams/{id}",
		Summary:     "Get a team",
		Tags:        []string{"Core teams"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, h.GetTeamByID)

	huma.Register(api, huma.Operation{
		OperationID: "list-core-teams-by-club",
		Method:      http.MethodGet,
		Path:        "/api/v1/core/teams/",
		Summary:     "List teams for a club",
		Tags:        []string{"Core teams"},
		Errors:      []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, h.ListTeamsByClub)

	huma.Register(api, huma.Operation{
		OperationID:   "add-core-team-player",
		Method:        http.MethodPost,
		Path:          "/api/v1/core/teams/{id}/players/",
		Summary:       "Add a player to a team",
		Tags:          []string{"Core teams"},
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, h.AddPlayerToTeam)

	huma.Register(api, huma.Operation{
		OperationID:   "remove-core-team-player",
		Method:        http.MethodDelete,
		Path:          "/api/v1/core/teams/{id}/players/{playerId}",
		Summary:       "Remove a player from a team",
		Tags:          []string{"Core teams"},
		DefaultStatus: http.StatusNoContent,
		Errors:        []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, h.RemovePlayerFromTeam)

	huma.Register(api, huma.Operation{
		OperationID: "list-core-team-players",
		Method:      http.MethodGet,
		Path:        "/api/v1/core/teams/{id}/players/",
		Summary:     "List a team's players",
		Tags:        []string{"Core teams"},
		Errors:      []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, h.ListPlayers)
}

func (h *TeamHandler) CreateTeam(ctx context.Context, input *CreateTeamInput) (*TeamOutput, error) {
	clubID, err := uuid.Parse(input.Body.ClubID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid club_id")
	}
	sportID, err := uuid.Parse(input.Body.SportID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid sport_id")
	}
	windowID, err := uuid.Parse(input.Body.WindowID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid window_id")
	}
	if input.Body.Name == "" {
		return nil, huma.Error400BadRequest("Name is required")
	}

	poolIDs := make([]uuid.UUID, 0, len(input.Body.PoolIDs))
	for _, raw := range input.Body.PoolIDs {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid pool_id: " + raw)
		}
		poolIDs = append(poolIDs, id)
	}

	team, err := h.repo.CreateTeam(ctx, ports.CreateTeamParams{
		ClubID:     clubID,
		SportID:    sportID,
		WindowID:   windowID,
		Type:       input.Body.Type,
		Name:       input.Body.Name,
		SeasonYear: input.Body.SeasonYear,
		PoolIDs:    poolIDs,
	})
	if err != nil {
		slog.Error("Failed to create team", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &TeamOutput{Body: *team}, nil
}

func (h *TeamHandler) GetTeamByID(ctx context.Context, input *GetTeamInput) (*TeamOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid UUID format")
	}

	team, err := h.repo.GetTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrTeamNotFound) {
			return nil, huma.Error404NotFound("Team not found")
		}
		slog.Error("Failed to get team", "error", err, "team_id", id)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &TeamOutput{Body: *team}, nil
}

func (h *TeamHandler) ListTeamsByClub(ctx context.Context, input *ListTeamsByClubInput) (*TeamsOutput, error) {
	clubID, err := uuid.Parse(input.ClubID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid club_id")
	}

	teams, err := h.repo.ListTeamsByClub(ctx, clubID, int32(input.Limit), int32(input.Offset))
	if err != nil {
		slog.Error("Failed to list teams", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &TeamsOutput{Body: teams}, nil
}

func (h *TeamHandler) AddPlayerToTeam(ctx context.Context, input *AddPlayerInput) (*EmptyOutput, error) {
	teamID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}
	playerID, err := uuid.Parse(input.Body.PlayerID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid player_id")
	}

	if err := h.repo.AddPlayerToTeam(ctx, teamID, playerID); err != nil {
		slog.Error("Failed to add player to team", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &EmptyOutput{}, nil
}

func (h *TeamHandler) RemovePlayerFromTeam(ctx context.Context, input *RemovePlayerInput) (*EmptyOutput, error) {
	teamID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}
	playerID, err := uuid.Parse(input.PlayerID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid player id")
	}

	if err := h.repo.RemovePlayerFromTeam(ctx, teamID, playerID); err != nil {
		slog.Error("Failed to remove player from team", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &EmptyOutput{}, nil
}

func (h *TeamHandler) ListPlayers(ctx context.Context, input *ListPlayersInput) (*PlayerIDsOutput, error) {
	teamID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid team id")
	}

	ids, err := h.repo.ListTeamPlayerIDs(ctx, teamID)
	if err != nil {
		slog.Error("Failed to list team players", "error", err)
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &PlayerIDsOutput{Body: ids}, nil
}
