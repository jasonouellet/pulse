-- ============================================================================
-- PROJECT PULSE — Core Module SQL Queries (sqlc)
-- Schema: core
-- Entities: users, player_profiles, parents_children, pools
-- ============================================================================

-- ----------------------------------------------------------------------------
-- USERS QUERIES
-- ----------------------------------------------------------------------------

-- name: CreateUser :one
INSERT INTO core.users (
    email,
    password_hash,
    first_name,
    last_name,
    phone,
    role
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, email, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at
FROM core.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, last_login_at, created_at, updated_at
FROM core.users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, first_name, last_name, phone, role, is_active, created_at
FROM core.users
ORDER BY last_name ASC, first_name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE core.users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    phone = COALESCE($4, phone),
    role = COALESCE($5, role),
    is_active = COALESCE($6, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, email, first_name, last_name, phone, role, is_active, updated_at;

-- name: UpdateUserLastLogin :exec
UPDATE core.users
SET last_login_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM core.users
WHERE id = $1;

-- ----------------------------------------------------------------------------
-- PLAYER PROFILES QUERIES
-- ----------------------------------------------------------------------------

-- name: CreatePlayerProfile :one
INSERT INTO core.player_profiles (
    user_id,
    first_name,
    last_name,
    date_of_birth,
    gender,
    medical_notes,
    emergency_contact_name,
    emergency_contact_phone
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, user_id, first_name, last_name, date_of_birth, gender, medical_notes, emergency_contact_name, emergency_contact_phone, created_at, updated_at;

-- name: GetPlayerProfileByID :one
SELECT id, user_id, first_name, last_name, date_of_birth, gender, medical_notes, emergency_contact_name, emergency_contact_phone, created_at, updated_at
FROM core.player_profiles
WHERE id = $1;

-- name: ListChildrenByParentID :many
SELECT p.id, p.user_id, p.first_name, p.last_name, p.date_of_birth, p.gender, pc.relationship, pc.is_primary_contact
FROM core.player_profiles p
JOIN core.parents_children pc ON p.id = pc.child_id
WHERE pc.parent_id = $1
ORDER BY p.date_of_birth DESC;

-- ----------------------------------------------------------------------------
-- PARENTS_CHILDREN RELATIONSHIP QUERIES
-- ----------------------------------------------------------------------------

-- name: LinkParentToChild :exec
INSERT INTO core.parents_children (
    parent_id,
    child_id,
    relationship,
    is_primary_contact
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (parent_id, child_id) DO UPDATE
SET relationship = EXCLUDED.relationship,
    is_primary_contact = EXCLUDED.is_primary_contact;

-- name: UnlinkParentFromChild :exec
DELETE FROM core.parents_children
WHERE parent_id = $1 AND child_id = $2;

-- ----------------------------------------------------------------------------
-- POOLS (AGE CATEGORIES) QUERIES
-- ----------------------------------------------------------------------------

-- name: CreatePool :one
INSERT INTO core.pools (
    sport_id,
    name,
    code,
    min_age,
    max_age,
    gender,
    season_year
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, sport_id, name, code, min_age, max_age, gender, season_year, is_active, created_at;

-- name: ListPoolsBySportAndSeason :many
SELECT id, sport_id, name, code, min_age, max_age, gender, season_year, is_active
FROM core.pools
WHERE sport_id = $1 AND season_year = $2 AND is_active = true
ORDER BY min_age ASC, code ASC;
