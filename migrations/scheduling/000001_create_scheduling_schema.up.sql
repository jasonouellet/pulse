-- ============================================================================
-- SCHEMA: scheduling
-- Second shared kernel (ADR-008 §4) — owns every temporal window on the
-- platform (season, pool, event) plus the fine-grained calendar (fields,
-- practices, matches). core and tournament reference scheduling via real
-- FK for their date_windows; scheduling never references core/tournament
-- back, to stay a stable dependency other modules can build on.
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS scheduling;

-- ----------------------------------------------------------------------------
-- ENUM: scheduling.window_type
-- ----------------------------------------------------------------------------
CREATE TYPE scheduling.window_type AS ENUM (
    'SEASON',
    'POOL',
    'EVENT'
);

-- ----------------------------------------------------------------------------
-- TABLE: scheduling.date_windows
-- The canonical start/end date range for a season, a pool's registration
-- window, or an event. core.pools and tournament.events reference a row
-- here instead of storing their own start_date/end_date.
-- ----------------------------------------------------------------------------
CREATE TABLE scheduling.date_windows (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type scheduling.window_type NOT NULL,
    label varchar(150) NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT ck_window_dates CHECK (end_date >= start_date)
);

-- ----------------------------------------------------------------------------
-- TABLE: scheduling.fields
-- club_id is a loose UUID reference on purpose, not a FK — scheduling must
-- stay dependency-free (a true kernel other modules build on), so it never
-- references core, even though core is itself a shared kernel. core.clubs
-- validates this at the application layer.
-- ----------------------------------------------------------------------------
CREATE TABLE scheduling.fields (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id uuid NOT NULL,
    name varchar(150) NOT NULL,
    address varchar(255),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: scheduling.sub_fields
-- A field can be split into playable sub-surfaces (e.g. one 11v11 field
-- divided into two 9v9s or four 7v7s) so several activities can run on it
-- concurrently.
-- ----------------------------------------------------------------------------
CREATE TABLE scheduling.sub_fields (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    field_id uuid NOT NULL REFERENCES scheduling.fields (id) ON DELETE CASCADE,
    name varchar(150) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT uk_sub_field_name UNIQUE (field_id, name)
);

-- ----------------------------------------------------------------------------
-- ENUM: scheduling.activity_type
-- ----------------------------------------------------------------------------
CREATE TYPE scheduling.activity_type AS ENUM (
    'PRACTICE',
    'MATCH'
);

-- ----------------------------------------------------------------------------
-- TABLE: scheduling.activities
-- A single practice or match, with a place and a time. owner_id is a loose
-- UUID reference on purpose — it can point at core.teams or
-- tournament.rosters depending on owner_type, and scheduling must never FK
-- into a peer/business schema (that would break its role as a stable
-- kernel). Validated at the application layer.
-- ----------------------------------------------------------------------------
CREATE TYPE scheduling.activity_owner_type AS ENUM (
    'TEAM',
    'ROSTER'
);

CREATE TABLE scheduling.activities (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type scheduling.activity_type NOT NULL,
    owner_type scheduling.activity_owner_type NOT NULL,
    owner_id uuid NOT NULL,
    field_id uuid REFERENCES scheduling.fields (id) ON DELETE SET NULL,
    sub_field_id uuid REFERENCES scheduling.sub_fields (id) ON DELETE SET NULL,
    starts_at timestamp with time zone NOT NULL,
    ends_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT ck_activity_times CHECK (ends_at > starts_at)
);

-- ----------------------------------------------------------------------------
-- TABLE: scheduling.attendances
-- player_id is a loose UUID reference, same reasoning as fields.club_id.
-- ----------------------------------------------------------------------------
CREATE TABLE scheduling.attendances (
    activity_id uuid NOT NULL REFERENCES scheduling.activities (id) ON DELETE CASCADE,
    player_id uuid NOT NULL,
    is_present boolean,
    recorded_at timestamp with time zone,
    PRIMARY KEY (activity_id, player_id)
);

-- ----------------------------------------------------------------------------
-- INDEXES
-- ----------------------------------------------------------------------------
CREATE INDEX idx_fields_club ON scheduling.fields (club_id);
CREATE INDEX idx_sub_fields_field ON scheduling.sub_fields (field_id);
CREATE INDEX idx_activities_owner ON scheduling.activities (owner_type, owner_id);
CREATE INDEX idx_activities_field_time ON scheduling.activities (field_id, starts_at);
CREATE INDEX idx_attendances_player ON scheduling.attendances (player_id);

-- ----------------------------------------------------------------------------
-- COMMENTS
-- ----------------------------------------------------------------------------
COMMENT ON SCHEMA scheduling IS 'Second shared kernel (ADR-008 §4): owns every temporal window and the fine-grained calendar. Other schemas reference it; it never references them back.';
COMMENT ON TABLE scheduling.date_windows IS 'Canonical start/end date range for a season, a pool, or an event.';
COMMENT ON TABLE scheduling.activities IS 'A single practice or match. owner_id is a loose reference (no FK) to core.teams or tournament.rosters, chosen via owner_type — kept loose so scheduling stays a dependency-free kernel.';
