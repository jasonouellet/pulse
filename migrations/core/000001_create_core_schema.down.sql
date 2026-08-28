DROP INDEX IF EXISTS core.idx_player_ratings_sport;
DROP INDEX IF EXISTS core.idx_player_position_prefs_player;
DROP INDEX IF EXISTS core.idx_positions_sport;
DROP INDEX IF EXISTS core.idx_pool_divisions_pool;
DROP INDEX IF EXISTS core.idx_pools_sport_season;
DROP INDEX IF EXISTS core.idx_parents_children_child;
DROP INDEX IF EXISTS core.idx_player_profiles_dob;
DROP INDEX IF EXISTS core.idx_users_role;
DROP INDEX IF EXISTS core.idx_users_email;
DROP INDEX IF EXISTS core.idx_sports_code;

DROP TABLE IF EXISTS core.pool_divisions;
DROP TABLE IF EXISTS core.pools;
DROP TABLE IF EXISTS core.parents_children;
DROP TABLE IF EXISTS core.player_ratings;
DROP TABLE IF EXISTS core.player_position_preferences;
DROP TABLE IF EXISTS core.positions;
DROP TABLE IF EXISTS core.player_profiles;
DROP TABLE IF EXISTS core.users;
DROP TABLE IF EXISTS core.sports;

DROP TYPE IF EXISTS core.relationship_type;
DROP TYPE IF EXISTS core.gender_category;
DROP TYPE IF EXISTS core.user_role;

DROP SCHEMA IF EXISTS core;
