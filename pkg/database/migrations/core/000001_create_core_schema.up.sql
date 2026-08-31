-- ============================================================================
-- PROJECT PULSE — Migration UP: Core Schema
-- Version: 000001
-- Description: Creates core schema, enum types, multi-sport reference,
--              users, player profiles, parent-child links, age pools, and
--              persistent teams (training groups / season teams).
--
-- Depends on: scheduling (core.pools.window_id references
-- scheduling.date_windows) — the scheduling migration must run first.
-- ============================================================================
SET search_path = core, public;

CREATE SCHEMA IF NOT EXISTS core;

-- ----------------------------------------------------------------------------
-- ENUM TYPES
-- ----------------------------------------------------------------------------
CREATE TYPE core.user_role AS ENUM (
    'SYSTEM_ADMIN',
    'SUPER_ADMIN',
    'CLUB_ADMIN',
    'TECHNICAL_DIRECTOR',
    'COACH',
    'GUARDIAN',
    'PLAYER'
);

CREATE TYPE core.gender_category AS ENUM (
    'MASCULINE',
    'FEMININE',
    'MIXED'
);

CREATE TYPE core.relationship_type AS ENUM (
    'MOTHER',
    'FATHER',
    'LEGAL_GUARDIAN',
    'OTHER'
);

CREATE TYPE core.registration_status AS ENUM (
    'PENDING',
    'ASSIGNED'
);

-- ADR-008 §6: a club can be a standalone club, or belong to a regional
-- association or a federation. org_type distinguishes the level; nesting
-- depth itself is unconstrained (see core.clubs.parent_club_id).
CREATE TYPE core.organization_type AS ENUM (
    'CLUB',
    'ASSOCIATION',
    'FEDERATION'
);

-- ADR-008 §1: TRAINING_GROUP and SEASON_TEAM are persistent club structure,
-- independent of any tournament/event — they live in core, not tournament.
CREATE TYPE core.team_type AS ENUM (
    'TRAINING_GROUP',
    'SEASON_TEAM'
);

-- ----------------------------------------------------------------------------
-- TABLE: core.clubs
-- The platform is multi-club (and multi-sport within each club). Sports
-- (below) stay global reference data — clubs don't redefine "soccer" — but
-- pools, rosters, and role grants are all club-scoped.
-- ----------------------------------------------------------------------------
CREATE TABLE core.clubs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_club_id uuid REFERENCES core.clubs (id) ON DELETE SET NULL,
    org_type core.organization_type NOT NULL DEFAULT 'CLUB',
    name varchar(150) NOT NULL,
    slug varchar(150) NOT NULL UNIQUE,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: core.sports
-- ----------------------------------------------------------------------------
CREATE TABLE core.sports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code varchar(32) NOT NULL UNIQUE,
    name varchar(100) NOT NULL,
    slug varchar(100) NOT NULL UNIQUE,
    default_periods int NOT NULL DEFAULT 2,
    default_period_duration_minutes int NOT NULL DEFAULT 25,
    rules_config jsonb NOT NULL DEFAULT '{}'::jsonb,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: core.users
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS core.users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    first_name varchar(100) NOT NULL,
    last_name varchar(100) NOT NULL,
    phone varchar(50),
    role core.user_role NOT NULL DEFAULT 'GUARDIAN',
    is_active boolean NOT NULL DEFAULT true,
    last_login_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- ----------------------------------------------------------------------------
-- TABLE: core.user_roles
-- A person can hold several roles at once (e.g., a GUARDIAN who is also a
-- COACH). is_primary picks the default context shown at login; the UI role
-- switcher lets them move between all granted roles. Authorization checks
-- must always test against the full granted set, never just is_primary.
-- ----------------------------------------------------------------------------
CREATE TABLE core.user_roles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
    club_id uuid REFERENCES core.clubs (id) ON DELETE CASCADE,
    role core.user_role NOT NULL,
    is_primary boolean NOT NULL DEFAULT false,
    granted_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    granted_by uuid REFERENCES core.users (id) ON DELETE SET NULL,
    CONSTRAINT ck_super_admin_is_platform_scoped CHECK (
        (role = 'SUPER_ADMIN' AND club_id IS null)
        OR (role <> 'SUPER_ADMIN' AND club_id IS NOT null)
    )
);

-- A club-scoped role can only be granted once per (user, role, club).
CREATE UNIQUE INDEX uk_user_role_club
ON core.user_roles (user_id, role, club_id)
WHERE club_id IS NOT null;

-- SUPER_ADMIN has no club_id to include in a composite unique key, so it
-- needs its own partial index to stay a one-time grant per user.
CREATE UNIQUE INDEX uk_user_super_admin
ON core.user_roles (user_id, role)
WHERE club_id IS null;

CREATE UNIQUE INDEX uk_one_primary_role_per_user
ON core.user_roles (user_id)
WHERE is_primary;

-- ----------------------------------------------------------------------------
-- TABLE: core.player_profiles
-- ----------------------------------------------------------------------------
CREATE TABLE core.player_profiles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid UNIQUE REFERENCES core.users (id) ON DELETE SET NULL,
    first_name varchar(100) NOT NULL,
    last_name varchar(100) NOT NULL,
    date_of_birth date NOT NULL,
    gender core.gender_category NOT NULL,
    medical_notes text,
    emergency_contact_name varchar(200),
    emergency_contact_phone varchar(30),
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: core.positions
-- Positions are sport-specific (ADR-004: zero schema refactor per sport) —
-- e.g., Attaquant/Défense/Demi/Gardien for soccer, different set for hockey.
-- ----------------------------------------------------------------------------
CREATE TABLE core.positions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE CASCADE,
    code varchar(32) NOT NULL,
    name varchar(100) NOT NULL,
    display_order int NOT NULL DEFAULT 0,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT uk_position_sport_code UNIQUE (sport_id, code)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.player_position_preferences
-- A player can have several preferred positions, ranked (1 = most preferred).
-- ----------------------------------------------------------------------------
CREATE TABLE core.player_position_preferences (
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE CASCADE,
    position_id uuid NOT NULL REFERENCES core.positions (id) ON DELETE CASCADE,
    preference_rank int NOT NULL DEFAULT 1,
    PRIMARY KEY (player_id, position_id),
    CONSTRAINT ck_preference_rank_positive CHECK (preference_rank > 0)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.player_ratings
-- A manually-set balance score per player, per sport, per season — scoped
-- by season_year because a player's rating evolves year over year and
-- overwriting it would lose that history (same principle as
-- pool_registrations: identity is permanent, the rating tied to it is not).
-- 0-100 scale is an assumption — flag for confirmation. May later be
-- superseded/fed by the `evaluation` schema once it exists; kept
-- lightweight and directly editable for now so roster balancing isn't
-- blocked on the full evaluation module. "Current" score = the row for the
-- active season_year, not just the latest by updated_at.
-- ----------------------------------------------------------------------------
CREATE TABLE core.player_ratings (
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE CASCADE,
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE CASCADE,
    season_year int NOT NULL,
    score smallint NOT NULL,
    updated_by uuid REFERENCES core.users (id) ON DELETE SET NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (player_id, sport_id, season_year),
    CONSTRAINT ck_score_range CHECK (score BETWEEN 0 AND 100)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.parents_children
-- ----------------------------------------------------------------------------
CREATE TABLE core.parents_children (
    parent_id uuid NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
    child_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE CASCADE,
    relationship core.relationship_type NOT NULL DEFAULT 'LEGAL_GUARDIAN',
    is_primary_contact boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (parent_id, child_id)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.pools
-- ----------------------------------------------------------------------------
CREATE TABLE core.pools (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id uuid NOT NULL REFERENCES core.clubs (id) ON DELETE CASCADE,
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE RESTRICT,
    window_id uuid NOT NULL REFERENCES scheduling.date_windows (id) ON DELETE RESTRICT,
    name varchar(100) NOT NULL,
    code varchar(50) NOT NULL,
    min_age int NOT NULL,
    max_age int NOT NULL,
    gender core.gender_category NOT NULL DEFAULT 'MIXED',
    season_year int NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT uk_pool_club_sport_season UNIQUE (club_id, sport_id, code, season_year)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.pool_divisions
-- ----------------------------------------------------------------------------
CREATE TABLE core.pool_divisions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    pool_id uuid NOT NULL REFERENCES core.pools (id) ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    code varchar(50) NOT NULL,
    display_order int NOT NULL DEFAULT 0,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT uk_division_pool_code UNIQUE (pool_id, code)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.pool_registrations
-- Explicit registration of a player into a pool for a given season — not
-- derived automatically from date_of_birth/gender, since age shifts every
-- year and registration is a deliberate per-season action.
-- ----------------------------------------------------------------------------
CREATE TABLE core.pool_registrations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE CASCADE,
    pool_id uuid NOT NULL REFERENCES core.pools (id) ON DELETE RESTRICT,
    division_id uuid REFERENCES core.pool_divisions (id) ON DELETE SET NULL,
    status core.registration_status NOT NULL DEFAULT 'PENDING',
    registered_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT uk_player_pool_season UNIQUE (player_id, pool_id)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.teams
-- Persistent club structure, independent of any event (ADR-008 §1):
--   TRAINING_GROUP — 1:1 with a single pool, organizes practices
--   SEASON_TEAM    — can span multiple pools (e.g., U9 + U10)
-- ----------------------------------------------------------------------------
CREATE TABLE core.teams (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id uuid NOT NULL REFERENCES core.clubs (id) ON DELETE CASCADE,
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE RESTRICT,
    window_id uuid NOT NULL REFERENCES scheduling.date_windows (id) ON DELETE RESTRICT,
    type core.team_type NOT NULL,
    name varchar(150) NOT NULL,
    season_year int NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: core.team_pools
-- Which pool(s) a team draws players from. TRAINING_GROUP has exactly one
-- row (1:1 with its pool); SEASON_TEAM can have several (combining ages).
-- ----------------------------------------------------------------------------
CREATE TABLE core.team_pools (
    team_id uuid NOT NULL REFERENCES core.teams (id) ON DELETE CASCADE,
    pool_id uuid NOT NULL REFERENCES core.pools (id) ON DELETE RESTRICT,
    PRIMARY KEY (team_id, pool_id)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.team_players
-- Team membership. No event scoping here — that only applies to
-- tournament.roster_players (event-specific alignments).
-- ----------------------------------------------------------------------------
CREATE TABLE core.team_players (
    team_id uuid NOT NULL REFERENCES core.teams (id) ON DELETE CASCADE,
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (team_id, player_id)
);

-- ----------------------------------------------------------------------------
-- INDEXES
-- ----------------------------------------------------------------------------
CREATE INDEX idx_sports_code ON core.sports (code);
CREATE INDEX idx_clubs_slug ON core.clubs (slug);
CREATE INDEX idx_clubs_parent ON core.clubs (parent_club_id);
CREATE INDEX idx_users_email ON core.users (email);
CREATE INDEX idx_user_roles_club ON core.user_roles (club_id);
CREATE INDEX idx_user_roles_role ON core.user_roles (role);
CREATE INDEX idx_player_profiles_dob ON core.player_profiles (date_of_birth);
CREATE INDEX idx_parents_children_child ON core.parents_children (child_id);
CREATE INDEX idx_pools_club_sport_season ON core.pools (club_id, sport_id, season_year);
CREATE INDEX idx_pools_window ON core.pools (window_id);
CREATE INDEX idx_pool_divisions_pool ON core.pool_divisions (pool_id);
CREATE INDEX idx_pool_registrations_pool ON core.pool_registrations (pool_id);
CREATE INDEX idx_pool_registrations_player ON core.pool_registrations (player_id);
CREATE INDEX idx_positions_sport ON core.positions (sport_id);
CREATE INDEX idx_player_position_prefs_player ON core.player_position_preferences (player_id);
CREATE INDEX idx_player_ratings_sport_season ON core.player_ratings (sport_id, season_year);
CREATE INDEX idx_teams_club_sport ON core.teams (club_id, sport_id);
CREATE INDEX idx_teams_window ON core.teams (window_id);
CREATE INDEX idx_team_pools_pool ON core.team_pools (pool_id);
CREATE INDEX idx_team_players_player ON core.team_players (player_id);

-- ----------------------------------------------------------------------------
-- DOCUMENTATION (COMMENT ON)
-- ----------------------------------------------------------------------------
COMMENT ON SCHEMA core IS 'Central schema managing identities, user accounts, player profiles, parent-child links, and age pools.';

COMMENT ON TABLE core.sports IS 'Reference table for multi-sport support abstraction (ADR-004).';
COMMENT ON COLUMN core.sports.id IS 'Primary key UUID generated via gen_random_uuid().';
COMMENT ON COLUMN core.sports.code IS 'Short code identifier for the sport (e.g., SOCCER, HOCKEY).';
COMMENT ON COLUMN core.sports.rules_config IS 'Flexible JSONB storing sport-specific configurations such as roster size rules or match periods.';

COMMENT ON TABLE core.clubs IS 'A club/organization. Platform is multi-club — most other data (pools, rosters, role grants) is scoped to one club.';

COMMENT ON TABLE core.users IS 'Platform accounts including administrators, technical directors, coaches, parents, and players. A single account can hold roles across multiple clubs — see core.user_roles.';
COMMENT ON TABLE core.user_roles IS 'Multi-role grants per user. A user can be GUARDIAN and COACH simultaneously; authorization checks the full set, not a single value.';
COMMENT ON COLUMN core.user_roles.is_primary IS 'Default role context shown at login/UI — not an authorization boundary by itself.';

COMMENT ON TABLE core.player_profiles IS 'Identity records for individual players.';
COMMENT ON COLUMN core.player_profiles.user_id IS 'Optional link to a dedicated user account if the player logs into the platform directly.';

COMMENT ON TABLE core.parents_children IS 'Junction table mapping parental links and primary contact designations.';
COMMENT ON COLUMN core.parents_children.is_primary_contact IS 'Flag indicating if this parent should receive primary communications for the child.';

COMMENT ON TABLE core.pools IS 'Age group pools and categories for player registration (e.g., U10F, U12M), scoped to one club. window_id references scheduling.date_windows for actual dates; season_year stays a human label.';
COMMENT ON COLUMN core.pools.code IS 'Division short code (e.g., U10F_D1, U12M_LOCAL).';

COMMENT ON TABLE core.pool_divisions IS 'Named skill/competitive levels within a pool (e.g., Division 1, Recreational). A roster is formed within one division.';
COMMENT ON COLUMN core.pool_divisions.display_order IS 'Controls the order divisions are listed in (e.g., Division 1 before Division 2).';

COMMENT ON TABLE core.pool_registrations IS 'Explicit, per-season registration of a player into a pool (and optionally a specific division).';
COMMENT ON COLUMN core.pool_registrations.status IS 'PENDING at registration; the application sets it to ASSIGNED when the player is added to an active TRAINING_GROUP team in core.team_players. No DB trigger — this is an application-layer responsibility.';

COMMENT ON TABLE core.clubs IS 'A club, association, or federation — see org_type and parent_club_id for the nesting hierarchy (ADR-008 §6).';
COMMENT ON COLUMN core.clubs.parent_club_id IS 'Nullable, self-referencing. A club can belong directly to a federation, or through an intermediate association — nesting depth is unconstrained.';

COMMENT ON TABLE core.teams IS 'Persistent club structure (training group or season team), independent of any tournament/event (ADR-008 §1).';
COMMENT ON COLUMN core.teams.window_id IS 'References scheduling.date_windows — teams no longer store their own start_date/end_date (ADR-008 §4).';
COMMENT ON TABLE core.team_pools IS 'Pool(s) a team draws players from. Multiple rows is how a season team combines age groups (e.g., U9 + U10).';
COMMENT ON TABLE core.team_players IS 'Team membership — no event scoping (that is tournament.roster_players).';
