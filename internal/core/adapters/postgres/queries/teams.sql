-- name: CreateTeam :one
INSERT INTO core.teams (
    club_id,
    sport_id,
    window_id,
    type,
    name,
    season_year
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at;

-- name: LinkTeamToPool :exec
INSERT INTO core.team_pools (
    team_id,
    pool_id
) VALUES (
    $1, $2
);

-- name: GetTeamByID :one
SELECT id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at
FROM core.teams
WHERE id = $1;

-- name: ListTeamsByClub :many
SELECT id, club_id, sport_id, window_id, type, name, season_year, created_at, updated_at
FROM core.teams
WHERE club_id = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: AddPlayerToTeam :exec
INSERT INTO core.team_players (
    team_id,
    player_id
) VALUES (
    $1, $2
);

-- name: RemovePlayerFromTeam :exec
DELETE FROM core.team_players
WHERE team_id = $1 AND player_id = $2;

-- name: ListTeamPlayerIDs :many
SELECT player_id
FROM core.team_players
WHERE team_id = $1;

-- name: ListTeamPoolIDs :many
SELECT pool_id
FROM core.team_pools
WHERE team_id = $1;
