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
    'PARENT',
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    default_periods INT NOT NULL DEFAULT 2,
    default_period_duration_minutes INT NOT NULL DEFAULT 25,
    rules_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ----------------------------------------------------------------------------
-- TABLE: core.users
-- ----------------------------------------------------------------------------
CREATE TABLE core.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(30),
    role core.user_role NOT NULL DEFAULT 'PARENT',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ----------------------------------------------------------------------------
-- TABLE: core.player_profiles
-- ----------------------------------------------------------------------------
CREATE TABLE core.player_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES core.users(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender core.gender_category NOT NULL,
    medical_notes TEXT,
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(30),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ----------------------------------------------------------------------------
-- TABLE: core.parents_children
-- ----------------------------------------------------------------------------
CREATE TABLE core.parents_children (
    parent_id UUID NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,
    child_id UUID NOT NULL REFERENCES core.player_profiles(id) ON DELETE CASCADE,
    relationship core.relationship_type NOT NULL DEFAULT 'LEGAL_GUARDIAN',
    is_primary_contact BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (parent_id, child_id)
);

-- ----------------------------------------------------------------------------
-- TABLE: core.pools
-- ----------------------------------------------------------------------------
CREATE TABLE core.pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES core.sports(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    min_age INT NOT NULL,
    max_age INT NOT NULL,
    gender core.gender_category NOT NULL DEFAULT 'MIXED',
    season_year INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_pool_sport_season UNIQUE (sport_id, code, season_year)
);

-- ----------------------------------------------------------------------------
-- INDEXES
-- ----------------------------------------------------------------------------
CREATE INDEX idx_sports_code ON core.sports(code);
CREATE INDEX idx_users_email ON core.users(email);
CREATE INDEX idx_users_role ON core.users(role);
CREATE INDEX idx_player_profiles_dob ON core.player_profiles(date_of_birth);
CREATE INDEX idx_parents_children_child ON core.parents_children(child_id);
CREATE INDEX idx_pools_sport_season ON core.pools(sport_id, season_year);

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
