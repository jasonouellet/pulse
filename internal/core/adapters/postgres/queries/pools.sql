-- CreatePool insère un nouveau bassin d'âge dans la structure du club.
-- name: CreatePool :one
INSERT INTO core.pools (sport_id, code, name)
VALUES ($1, $2, $3)
RETURNING *;

-- GetPoolByID récupère les détails d'un bassin d'âge par son UUID.
-- name: GetPoolByID :one
SELECT * FROM core.pools
WHERE id = $1;