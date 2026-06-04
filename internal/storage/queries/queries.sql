-- name: SaveSnapshot :one
INSERT INTO snapshots (source, data, period_start, period_end)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: SaveDailyMetrics :one
INSERT INTO daily_metrics (source, metric_date, revenue, failed_payments)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: SaveDigest :one
INSERT INTO digests (content, has_critical_alerts)
VALUES ($1, $2)
RETURNING *;
