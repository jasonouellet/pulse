-- CreatePool inserts a new age pool into the club structure.
-- name: CreatePool :one
INSERT INTO core.pools (sport_id, code, name)
VALUES ($1, $2, $3)
RETURNING *;

-- GetPoolByID retrieves the details of an age pool by its UUID.
-- name: GetPoolByID :one
SELECT * FROM core.pools
WHERE id = $1;
