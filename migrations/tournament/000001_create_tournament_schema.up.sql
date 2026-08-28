-- ============================================================================
-- SCHEMA: tournament
-- Manages training groups, season teams, and tournament rosters.
--
-- Per ADR-003 (§2.B): core is a shared kernel — real FKs to core.* are
-- expected and preferred over app-level UUID validation. What's still
-- forbidden is a direct dependency between peer modules (tournament must
-- never reference scheduling/finance/evaluation, or vice versa).
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS tournament;

-- ----------------------------------------------------------------------------
-- ENUM: tournament.roster_type
-- ----------------------------------------------------------------------------
CREATE TYPE tournament.roster_type AS ENUM (
    'TRAINING_GROUP',
    'SEASON_TEAM',
    'EVENT_ROSTER'
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.events
-- Minimal scaffold for now — a specific tournament/event that EVENT_ROSTER
-- rosters are formed for. Registration and bracket generation are out of
-- scope for this migration.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE RESTRICT,
    name varchar(150) NOT NULL,
    season_year int NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT ck_event_dates CHECK (end_date >= start_date)
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.rosters
-- A roster is one of three types:
--   TRAINING_GROUP — 1:1 with a single pool, season-long, organizes practices
--   SEASON_TEAM    — can span multiple pools (e.g., U9 + U10), season-long
--   EVENT_ROSTER   — formed per tournament/event, can also span pools
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.rosters (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE RESTRICT,
    type tournament.roster_type NOT NULL,
    name varchar(150) NOT NULL,
    season_year int NOT NULL,
    start_date date,
    end_date date,
    event_id uuid REFERENCES tournament.events (id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT ck_event_roster_has_event CHECK (
        (type = 'EVENT_ROSTER' AND event_id IS NOT null)
        OR (type <> 'EVENT_ROSTER' AND event_id IS null)
    ),
    CONSTRAINT ck_roster_dates CHECK (
        (type = 'EVENT_ROSTER' AND start_date IS null AND end_date IS null)
        OR (
            type <> 'EVENT_ROSTER'
            AND start_date IS NOT null
            AND end_date IS NOT null
            AND end_date >= start_date
        )
    )
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.roster_pools
-- Links a roster to the pool(s) it draws players from. A TRAINING_GROUP has
-- exactly one row here (1:1 with its pool); SEASON_TEAM and EVENT_ROSTER can
-- have several (e.g., a season team combining U9 and U10).
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.roster_pools (
    roster_id uuid NOT NULL REFERENCES tournament.rosters (id) ON DELETE CASCADE,
    pool_id uuid NOT NULL REFERENCES core.pools (id) ON DELETE RESTRICT,
    PRIMARY KEY (roster_id, pool_id)
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.roster_players
-- Membership. event_id is denormalized from the parent roster (NULL for
-- non-event rosters) specifically to enforce "one player, one roster per
-- event" with a plain partial unique index — no trigger required.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.roster_players (
    roster_id uuid NOT NULL REFERENCES tournament.rosters (id) ON DELETE CASCADE,
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE RESTRICT,
    event_id uuid REFERENCES tournament.events (id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (roster_id, player_id)
);

CREATE UNIQUE INDEX uk_one_roster_per_player_per_event
ON tournament.roster_players (event_id, player_id)
WHERE event_id IS NOT null;

-- ----------------------------------------------------------------------------
-- INDEXES
-- ----------------------------------------------------------------------------
CREATE INDEX idx_rosters_type_season ON tournament.rosters (type, season_year);
CREATE INDEX idx_rosters_event ON tournament.rosters (event_id);
CREATE INDEX idx_roster_pools_pool ON tournament.roster_pools (pool_id);
CREATE INDEX idx_roster_players_player ON tournament.roster_players (player_id);
CREATE INDEX idx_events_sport_season ON tournament.events (sport_id, season_year);

-- ----------------------------------------------------------------------------
-- COMMENTS
-- ----------------------------------------------------------------------------
COMMENT ON SCHEMA tournament IS 'Training groups, season teams, and tournament rosters. Isolated per ADR-003 — no cross-schema JOINs or FKs.';
COMMENT ON TABLE tournament.rosters IS 'A group of players: a season-long training group, a season-long competitive team, or a one-off tournament roster.';
COMMENT ON COLUMN tournament.rosters.type IS 'TRAINING_GROUP: 1:1 with a pool. SEASON_TEAM: can span multiple pools. EVENT_ROSTER: formed per event, can also span pools.';
COMMENT ON COLUMN tournament.rosters.season_year IS 'Human label for filtering (e.g., 2026), independent of the actual start_date/end_date range.';
COMMENT ON COLUMN tournament.rosters.start_date IS 'Required for TRAINING_GROUP/SEASON_TEAM. Null for EVENT_ROSTER, which inherits its dates from tournament.events. May span two calendar years (e.g., a fall-to-summer club season).';
COMMENT ON TABLE tournament.roster_pools IS 'Pool(s) a roster draws its players from. Multiple rows is how a season team combines age groups (e.g., U9 + U10).';
COMMENT ON TABLE tournament.roster_players IS 'Roster membership. event_id is set only for EVENT_ROSTER members, enabling the one-roster-per-event-per-player constraint.';
