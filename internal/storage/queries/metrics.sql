-- name: GetLastSevenDays :many
SELECT * 
FROM daily_metrics
WHERE source = $1
AND metric_date >= NOW() - INTERVAL '7 days'
ORDER BY metric_date DESC;