DROP INDEX IF EXISTS tournament.idx_events_sport_season;
DROP INDEX IF EXISTS tournament.idx_roster_players_player;
DROP INDEX IF EXISTS tournament.idx_roster_pools_pool;
DROP INDEX IF EXISTS tournament.idx_rosters_event;
DROP INDEX IF EXISTS tournament.idx_rosters_type_season;
DROP INDEX IF EXISTS tournament.uk_one_roster_per_player_per_event;

DROP TABLE IF EXISTS tournament.roster_players;
DROP TABLE IF EXISTS tournament.roster_pools;
DROP TABLE IF EXISTS tournament.rosters;
DROP TABLE IF EXISTS tournament.events;

DROP TYPE IF EXISTS tournament.roster_type;

DROP SCHEMA IF EXISTS tournament;
