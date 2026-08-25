CREATE SCHEMA IF NOT EXISTS core;

-- Commentaire sur la table
CREATE TABLE core.pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE core.pools IS 'Bassins d''âges regroupant les sportifs d''une même catégorie (ex: U9-U10 Été 2026).';
COMMENT ON COLUMN core.pools.sport_id IS 'Référence au sport concerné (ex: Soccer, Hockey).';
COMMENT ON COLUMN core.pools.code IS 'Identifiant court pour l''affichage (ex: POOL-U10F).';