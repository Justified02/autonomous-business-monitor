-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

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
    metrics     JSONB NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, metric_date)   
);

CREATE TABLE digests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sources_count INT NOT NULL,
    alerts_count INT NOT NULL DEFAULT 0,
    content     TEXT NOT NULL,
    payload     JSONB NOT NULL,
    delivered   BOOLEAN NOT NULL DEFAULT FALSE,
    delivered_at TIMESTAMPTZ
);

CREATE INDEX idx_snapshots_source_fetched  ON snapshots(source, fetched_at DESC);
CREATE INDEX idx_daily_metrics_source_date ON daily_metrics(source, metric_date DESC);
CREATE INDEX idx_digests_generated         ON digests(generated_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_digests_generated;
DROP INDEX IF EXISTS idx_daily_metrics_source_date;
DROP INDEX IF EXISTS idx_snapshots_source_fetched;
DROP TABLE IF EXISTS digests;
DROP TABLE IF EXISTS daily_metrics;
DROP TABLE IF EXISTS snapshots;
