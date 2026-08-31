DROP INDEX IF EXISTS tournament.idx_roster_coaches_coach;
DROP INDEX IF EXISTS tournament.idx_roster_players_player;
DROP INDEX IF EXISTS tournament.uk_one_roster_per_player_per_event;
DROP INDEX IF EXISTS tournament.idx_team_entries_club;
DROP INDEX IF EXISTS tournament.idx_team_entries_event;
DROP INDEX IF EXISTS tournament.idx_event_eligibility_event;
DROP INDEX IF EXISTS tournament.idx_events_window;
DROP INDEX IF EXISTS tournament.idx_events_sport;
DROP INDEX IF EXISTS tournament.idx_events_organizing_club;

DROP TABLE IF EXISTS tournament.roster_coaches;
DROP TABLE IF EXISTS tournament.roster_players;
DROP TABLE IF EXISTS tournament.rosters;
DROP TABLE IF EXISTS tournament.team_entries;
DROP TABLE IF EXISTS tournament.event_eligibility;
DROP TABLE IF EXISTS tournament.events;

DROP SCHEMA IF EXISTS tournament;
