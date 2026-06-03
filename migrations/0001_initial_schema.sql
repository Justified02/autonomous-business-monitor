-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- SQL to create your tables goes here
CREATE TABLE snapshots (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source      TEXT NOT NULL,
    fetched_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data        JSONB NOT NULL,
    period_start TIMESTAMPTZ,
    period_end  TIMESTAMPTZ
);

CREATE TABLE daily_metrics (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source      TEXT NOT NULL,
    metric_date DATE NOT NULL,
    revenue     NUMERIC(10,2) NOT NULL,
    failed_payments INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, metric_date)   
);

CREATE TABLE digests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content     TEXT NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    has_critical_alerts BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
-- SQL to reverse the migration goes here (drop the tables)
DROP TABLE IF EXISTS digests;
DROP TABLE IF EXISTS daily_metrics;
DROP TABLE IF EXISTS snapshots;