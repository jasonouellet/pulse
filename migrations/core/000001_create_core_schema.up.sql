CREATE SCHEMA IF NOT EXISTS core;

CREATE TABLE core.pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE core.pools IS 'Age pools grouping athletes of the same category (e.g. U9-U10 Summer 2026).';
COMMENT ON COLUMN core.pools.sport_id IS 'Reference to the relevant sport (e.g. Soccer, Hockey).';
COMMENT ON COLUMN core.pools.code IS 'Short identifier for display purposes (e.g. POOL-U10F).';