-- name: GetLastSevenDays :many
SELECT * 
FROM daily_metrics
WHERE source = $1
AND metric_date >= NOW() - INTERVAL '7 days'
ORDER BY metric_date DESC;


-- name: GetPastDigests :many
SELECT *
FROM digests
ORDER BY generated_at DESC
LIMIT 30;


-- name: GetMetricsTrend :many
SELECT *
FROM daily_metrics
WHERE source = $1
ORDER BY metric_date DESC
LIMIT 30;