-- name: GetLogsByDeployment :many
SELECT id, deployment_id, timestamp, message FROM logs WHERE deployment_id = ? ORDER BY timestamp;

-- name: CreateLog :one
INSERT INTO logs (deployment_id, message)
VALUES (?, ?)
RETURNING id, deployment_id, timestamp, message;

-- name: DeleteLogsByDeployment :exec
DELETE FROM logs WHERE deployment_id = ?;
