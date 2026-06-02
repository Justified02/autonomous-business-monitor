-- name: SaveDigest :one
INSERT INTO digests (sources_count, alerts_count, content, payload)
VALUES ($1, $2, $3, $4)
RETURNING id, generated_at;

-- name: GetRecentDigests :many
SELECT id, generated_at, sources_count, alerts_count, content, delivered
FROM digests
ORDER BY generated_at DESC
LIMIT $1;

-- name: MarkDigestDelivered :exec
UPDATE digests
SET delivered = true, delivered_at = NOW()
WHERE id = $1;