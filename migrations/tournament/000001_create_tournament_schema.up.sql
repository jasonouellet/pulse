-- ============================================================================
-- SCHEMA: tournament
-- Manages events and the event-specific alignments (rosters) formed for
-- them. Training groups and season teams live in core (ADR-008 §1) — they
-- are not tournament concepts.
--
-- Per ADR-003 (§2.B) and ADR-008 (§4): core and scheduling are shared
-- kernels. tournament references both via real FK. It must never reference
-- another peer/business schema (finance, evaluation) directly.
--
-- Depends on: core, scheduling (both must be migrated first).
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS tournament;

-- ----------------------------------------------------------------------------
-- TABLE: tournament.events
-- A tournament/event, organized by one club. window_id replaces raw
-- start_date/end_date (ADR-008 §4) — the coarse date range for the whole
-- event; scheduling.activities carries the per-match/field detail.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organizing_club_id uuid NOT NULL REFERENCES core.clubs (id) ON DELETE RESTRICT,
    sport_id uuid NOT NULL REFERENCES core.sports (id) ON DELETE RESTRICT,
    window_id uuid NOT NULL REFERENCES scheduling.date_windows (id) ON DELETE RESTRICT,
    name varchar(150) NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.event_eligibility
-- Who can register for an event, expressed as one or more age/gender
-- brackets. Multiple rows per event allow non-contiguous eligibility (e.g.
-- U4-U8 AND U13-M-and-up on the same event) — see ADR-008 §2.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.event_eligibility (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL REFERENCES tournament.events (id) ON DELETE CASCADE,
    min_age int NOT NULL,
    max_age int,
    gender core.gender_category NOT NULL DEFAULT 'MIXED',
    label varchar(150),
    CONSTRAINT ck_eligibility_age_range CHECK (max_age IS NULL OR max_age >= min_age)
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.team_entries
-- A club's decision to enter a team into an event, for one of its own
-- pools. This is NOT yet a list of players — that's the roster below.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.team_entries (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL REFERENCES tournament.events (id) ON DELETE CASCADE,
    club_id uuid NOT NULL REFERENCES core.clubs (id) ON DELETE CASCADE,
    pool_id uuid NOT NULL REFERENCES core.pools (id) ON DELETE RESTRICT,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    CONSTRAINT uk_team_entry_event_club_pool UNIQUE (event_id, club_id, pool_id)
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.rosters
-- The alignment (French: alignement) for one team entry — the specific
-- players, decided closer to the event date. 1:1 with its team entry.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.rosters (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    team_entry_id uuid NOT NULL UNIQUE REFERENCES tournament.team_entries (id) ON DELETE CASCADE,
    name varchar(150) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.roster_players
-- Membership. event_id is denormalized (sourced from
-- roster -> team_entry -> event at write time) specifically to enforce
-- "one player, one roster per event" with a plain partial unique index.
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.roster_players (
    roster_id uuid NOT NULL REFERENCES tournament.rosters (id) ON DELETE CASCADE,
    player_id uuid NOT NULL REFERENCES core.player_profiles (id) ON DELETE RESTRICT,
    event_id uuid NOT NULL REFERENCES tournament.events (id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (roster_id, player_id)
);

CREATE UNIQUE INDEX uk_one_roster_per_player_per_event
    ON tournament.roster_players (event_id, player_id);

-- ----------------------------------------------------------------------------
-- TABLE: tournament.roster_coaches
-- Per-event coach assignment — can differ from the player's regular
-- core.teams coach. A roster can have more than one (assistant coaches).
-- ----------------------------------------------------------------------------
CREATE TABLE tournament.roster_coaches (
    roster_id uuid NOT NULL REFERENCES tournament.rosters (id) ON DELETE CASCADE,
    coach_user_id uuid NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
    assigned_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (roster_id, coach_user_id)
);

-- ----------------------------------------------------------------------------
-- INDEXES
-- ----------------------------------------------------------------------------
CREATE INDEX idx_events_organizing_club ON tournament.events (organizing_club_id);
CREATE INDEX idx_events_sport ON tournament.events (sport_id);
CREATE INDEX idx_events_window ON tournament.events (window_id);
CREATE INDEX idx_event_eligibility_event ON tournament.event_eligibility (event_id);
CREATE INDEX idx_team_entries_event ON tournament.team_entries (event_id);
CREATE INDEX idx_team_entries_club ON tournament.team_entries (club_id);
CREATE INDEX idx_roster_players_player ON tournament.roster_players (player_id);
CREATE INDEX idx_roster_coaches_coach ON tournament.roster_coaches (coach_user_id);

-- ----------------------------------------------------------------------------
-- COMMENTS
-- ----------------------------------------------------------------------------
COMMENT ON SCHEMA tournament IS 'Events and the event-specific alignments (rosters) formed for them. Training groups and season teams live in core (ADR-008 §1).';
COMMENT ON TABLE tournament.events IS 'A tournament/event organized by one club. Other clubs enter teams via team_entries.';
COMMENT ON TABLE tournament.event_eligibility IS 'Age/gender brackets eligible for an event. Multiple rows allow non-contiguous eligibility on the same event.';
COMMENT ON TABLE tournament.team_entries IS 'A club''s decision to enter one of its pools into an event — distinct from the roster (specific players), decided separately and later.';
COMMENT ON TABLE tournament.rosters IS 'The alignment (alignement) for one team entry: the specific players, 1:1 with its team_entry.';
COMMENT ON TABLE tournament.roster_players IS 'Roster membership. event_id enables the one-roster-per-event-per-player constraint.';
COMMENT ON TABLE tournament.roster_coaches IS 'Per-event coach assignment, independent of the player''s regular core.teams coach.';
