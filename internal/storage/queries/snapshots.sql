-- name: SaveSnapshot :exec
INSERT INTO snapshots (source, data, period_start, period_end)
VALUES ($1, $2, $3, $4);

-- name: GetRecentSnapshots :many
SELECT id, source, fetched_at, data, period_start, period_end
FROM snapshots
WHERE source = $1
    AND fetched_at >= NOW() - ($2::text || ' days')::INTERVAL
ORDER BY fetched_at DESC;