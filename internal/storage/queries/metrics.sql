-- name: UpsertDailyMetrics :exec
INSERT INTO daily_metrics (source, metric_date, metrics)
VALUES ($1, $2, $3)
ON CONFLICT (source, metric_date)
DO UPDATE SET metrics = EXCLUDED.metrics;

-- name: GetRollingMetrics :many
SELECT metric_date, metrics
FROM daily_metrics
WHERE source = $1
    AND metric_date >= CURRENT_DATE - ($2::int * INTERVAL '1 day')
ORDER BY metric_date DESC;