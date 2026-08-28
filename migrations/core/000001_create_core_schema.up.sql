-- ============================================================================
-- PROJECT PULSE — Migration UP: Core Schema
-- Version: 000001
-- Description: Creates core schema, enum types, multi-sport reference,
--              users, player profiles, parent-child links, and age pools.
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS core;

-- ----------------------------------------------------------------------------
-- ENUM TYPES
-- ----------------------------------------------------------------------------
CREATE TYPE core.user_role AS ENUM (
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
CREATE TABLE core.users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    first_name varchar(100) NOT NULL,
    last_name varchar(100) NOT NULL,
    phone varchar(30),
    role core.user_role NOT NULL DEFAULT 'GUARDIAN',
    is_active boolean NOT NULL DEFAULT true,
    last_login_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

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
-- A manually-set balance score per player, per sport (a player's soccer
-- rating is meaningless for hockey). 0-100 scale is an assumption — flag
-- for confirmation. May later be superseded/fed by the `evaluation` schema
-- once it exists; kept lightweight and directly editable for now so roster
-- balancing isn't blocked on the full evaluation module.
-- ----------------------------------------------------------------------------
CREATE TABLE core.player_ratings (
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE CASCADE,
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE CASCADE,
    score smallint NOT NULL,
    updated_by uuid REFERENCES core.users (id) ON DELETE SET NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (player_id, sport_id),
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
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE RESTRICT,
    name varchar(100) NOT NULL,
    code varchar(50) NOT NULL,
    min_age int NOT NULL,
    max_age int NOT NULL,
    gender core.gender_category NOT NULL DEFAULT 'MIXED',
    season_year int NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT uk_pool_sport_season UNIQUE (sport_id, code, season_year)
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
-- INDEXES
-- ----------------------------------------------------------------------------
CREATE INDEX idx_sports_code ON core.sports (code);
CREATE INDEX idx_users_email ON core.users (email);
CREATE INDEX idx_users_role ON core.users (role);
CREATE INDEX idx_player_profiles_dob ON core.player_profiles (date_of_birth);
CREATE INDEX idx_parents_children_child ON core.parents_children (child_id);
CREATE INDEX idx_pools_sport_season ON core.pools (sport_id, season_year);
CREATE INDEX idx_pool_divisions_pool ON core.pool_divisions (pool_id);
CREATE INDEX idx_positions_sport ON core.positions (sport_id);
CREATE INDEX idx_player_position_prefs_player ON core.player_position_preferences (player_id);
CREATE INDEX idx_player_ratings_sport ON core.player_ratings (sport_id);

-- ----------------------------------------------------------------------------
-- DOCUMENTATION (COMMENT ON)
-- ----------------------------------------------------------------------------
COMMENT ON SCHEMA core IS 'Central schema managing identities, user accounts, player profiles, parent-child links, and age pools.';

COMMENT ON TABLE core.sports IS 'Reference table for multi-sport support abstraction (ADR-004).';
COMMENT ON COLUMN core.sports.id IS 'Primary key UUID generated via gen_random_uuid().';
COMMENT ON COLUMN core.sports.code IS 'Short code identifier for the sport (e.g., SOCCER, HOCKEY).';
COMMENT ON COLUMN core.sports.rules_config IS 'Flexible JSONB storing sport-specific configurations such as roster size rules or match periods.';

COMMENT ON TABLE core.users IS 'Platform accounts including administrators, technical directors, coaches, parents, and players.';
COMMENT ON COLUMN core.users.role IS 'Global RBAC role of the user.';

COMMENT ON TABLE core.player_profiles IS 'Identity records for individual players.';
COMMENT ON COLUMN core.player_profiles.user_id IS 'Optional link to a dedicated user account if the player logs into the platform directly.';

COMMENT ON TABLE core.parents_children IS 'Junction table mapping parental links and primary contact designations.';
COMMENT ON COLUMN core.parents_children.is_primary_contact IS 'Flag indicating if this parent should receive primary communications for the child.';

COMMENT ON TABLE core.pools IS 'Age group pools and categories for player registration (e.g., U10F, U12M).';
COMMENT ON COLUMN core.pools.code IS 'Division short code (e.g., U10F_D1, U12M_LOCAL).';

COMMENT ON TABLE core.pool_divisions IS 'Named skill/competitive levels within a pool (e.g., Division 1, Recreational). A roster is formed within one division.';
COMMENT ON COLUMN core.pool_divisions.display_order IS 'Controls the order divisions are listed in (e.g., Division 1 before Division 2).';
