DROP INDEX IF EXISTS scheduling.idx_attendances_player;
DROP INDEX IF EXISTS scheduling.idx_activities_field_time;
DROP INDEX IF EXISTS scheduling.idx_activities_owner;
DROP INDEX IF EXISTS scheduling.idx_sub_fields_field;
DROP INDEX IF EXISTS scheduling.idx_fields_club;

DROP TABLE IF EXISTS scheduling.attendances;
DROP TABLE IF EXISTS scheduling.activities;
DROP TABLE IF EXISTS scheduling.sub_fields;
DROP TABLE IF EXISTS scheduling.fields;
DROP TABLE IF EXISTS scheduling.date_windows;

DROP TYPE IF EXISTS scheduling.activity_owner_type;
DROP TYPE IF EXISTS scheduling.activity_type;
DROP TYPE IF EXISTS scheduling.window_type;

DROP SCHEMA IF EXISTS scheduling;
