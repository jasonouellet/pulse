package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	coreHTTP "pulse/internal/core/adapters/http"
	"pulse/internal/core/ports"
)

// MockTeamRepository implements ports.TeamRepository in-memory for fast unit testing.
type MockTeamRepository struct {
	teams   map[uuid.UUID]ports.TeamDTO
	players map[uuid.UUID][]uuid.UUID
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{
		teams:   make(map[uuid.UUID]ports.TeamDTO),
		players: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (m *MockTeamRepository) CreateTeam(ctx context.Context, p ports.CreateTeamParams) (*ports.TeamDTO, error) {
	id := uuid.New()
	t := ports.TeamDTO{
		ID:         id,
		ClubID:     p.ClubID,
		SportID:    p.SportID,
		WindowID:   p.WindowID,
		Type:       p.Type,
		Name:       p.Name,
		SeasonYear: p.SeasonYear,
		PoolIDs:    p.PoolIDs,
	}
	m.teams[id] = t
	return &t, nil
}

func (m *MockTeamRepository) GetTeamByID(ctx context.Context, id uuid.UUID) (*ports.TeamDTO, error) {
	t, exists := m.teams[id]
	if !exists {
		return nil, ports.ErrTeamNotFound
	}
	return &t, nil
}

func (m *MockTeamRepository) ListTeamsByClub(ctx context.Context, clubID uuid.UUID, limit, offset int32) ([]ports.TeamDTO, error) {
	var list []ports.TeamDTO
	for _, t := range m.teams {
		if t.ClubID == clubID {
			list = append(list, t)
		}
	}
	return list, nil
}

func (m *MockTeamRepository) AddPlayerToTeam(ctx context.Context, teamID, playerID uuid.UUID) error {
	m.players[teamID] = append(m.players[teamID], playerID)
	return nil
}

func (m *MockTeamRepository) RemovePlayerFromTeam(ctx context.Context, teamID, playerID uuid.UUID) error {
	ids := m.players[teamID]
	for i, id := range ids {
		if id == playerID {
			m.players[teamID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

func (m *MockTeamRepository) ListTeamPlayerIDs(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	return m.players[teamID], nil
}

func newTeamTestAPI(repo ports.TeamRepository) (*chi.Mux, huma.API) {
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("PULSE test API", "1.0.0"))
	coreHTTP.NewTeamHandler(repo).RegisterRoutes(api)
	return r, api
}

func TestCreateTeam_Success(t *testing.T) {
	repo := NewMockTeamRepository()
	r, _ := newTeamTestAPI(repo)

	payload := map[string]any{
		"club_id":     uuid.New().String(),
		"sport_id":    uuid.New().String(),
		"window_id":   uuid.New().String(),
		"type":        "TRAINING_GROUP",
		"name":        "U10-M Division 1",
		"season_year": 2026,
		"pool_ids":    []string{uuid.New().String()},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/core/teams/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var created ports.TeamDTO
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.Name != "U10-M Division 1" {
		t.Errorf("expected name %q, got %q", "U10-M Division 1", created.Name)
	}
	if len(created.PoolIDs) != 1 {
		t.Errorf("expected 1 pool id, got %d", len(created.PoolIDs))
	}
}

func TestCreateTeam_InvalidClubID(t *testing.T) {
	repo := NewMockTeamRepository()
	r, _ := newTeamTestAPI(repo)

	payload := map[string]any{
		"club_id":     "not-a-uuid",
		"sport_id":    uuid.New().String(),
		"window_id":   uuid.New().String(),
		"type":        "TRAINING_GROUP",
		"name":        "U10-M Division 1",
		"season_year": 2026,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/core/teams/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}

func TestCreateTeam_MissingName(t *testing.T) {
	repo := NewMockTeamRepository()
	r, _ := newTeamTestAPI(repo)

	payload := map[string]any{
		"club_id":     uuid.New().String(),
		"sport_id":    uuid.New().String(),
		"window_id":   uuid.New().String(),
		"type":        "TRAINING_GROUP",
		"season_year": 2026,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/core/teams/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}

func TestGetTeamByID_Found(t *testing.T) {
	repo := NewMockTeamRepository()
	teamID := uuid.New()
	repo.teams[teamID] = ports.TeamDTO{ID: teamID, Name: "U10-M Division 1"}
	r, _ := newTeamTestAPI(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/core/teams/"+teamID.String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestGetTeamByID_NotFound(t *testing.T) {
	repo := NewMockTeamRepository()
	r, _ := newTeamTestAPI(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/core/teams/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}

func TestGetTeamByID_InvalidUUID(t *testing.T) {
	repo := NewMockTeamRepository()
	r, _ := newTeamTestAPI(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/core/teams/not-a-uuid", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}

func TestListTeamsByClub_Success(t *testing.T) {
	repo := NewMockTeamRepository()
	clubID := uuid.New()
	repo.teams[uuid.New()] = ports.TeamDTO{ID: uuid.New(), ClubID: clubID, Name: "Team A"}
	r, _ := newTeamTestAPI(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/core/teams/?club_id="+clubID.String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestAddPlayerToTeam_Success(t *testing.T) {
	repo := NewMockTeamRepository()
	teamID := uuid.New()
	playerID := uuid.New()
	r, _ := newTeamTestAPI(repo)

	payload := map[string]string{"player_id": playerID.String()}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/core/teams/"+teamID.String()+"/players/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, rec.Code, rec.Body.String())
	}
	if len(repo.players[teamID]) != 1 || repo.players[teamID][0] != playerID {
		t.Errorf("expected player %s to be added to team %s", playerID, teamID)
	}
}

func TestRemovePlayerFromTeam_Success(t *testing.T) {
	repo := NewMockTeamRepository()
	teamID := uuid.New()
	playerID := uuid.New()
	repo.players[teamID] = []uuid.UUID{playerID}
	r, _ := newTeamTestAPI(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/core/teams/"+teamID.String()+"/players/"+playerID.String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, rec.Code, rec.Body.String())
	}
	if len(repo.players[teamID]) != 0 {
		t.Errorf("expected player to be removed, got %v", repo.players[teamID])
	}
}

func TestListPlayers_Success(t *testing.T) {
	repo := NewMockTeamRepository()
	teamID := uuid.New()
	playerID := uuid.New()
	repo.players[teamID] = []uuid.UUID{playerID}
	r, _ := newTeamTestAPI(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/core/teams/"+teamID.String()+"/players/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var ids []uuid.UUID
	if err := json.NewDecoder(rec.Body).Decode(&ids); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(ids) != 1 || ids[0] != playerID {
		t.Errorf("expected [%s], got %v", playerID, ids)
	}
}
