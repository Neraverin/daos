-- name: GetAllHosts :many
SELECT id, name, hostname, port, username, ssh_key_path, created_at, updated_at FROM hosts ORDER BY name;

-- name: GetHost :one
SELECT id, name, hostname, port, username, ssh_key_path, created_at, updated_at FROM hosts WHERE id = ?;

-- name: CreateHost :one
INSERT INTO hosts (name, hostname, port, username, ssh_key_path)
VALUES (?, ?, ?, ?, ?)
RETURNING id, name, hostname, port, username, ssh_key_path, created_at, updated_at;

-- name: UpdateHost :one
UPDATE hosts SET name = ?, hostname = ?, port = ?, username = ?, ssh_key_path = ?, updated_at = datetime('now')
WHERE id = ?
RETURNING id, name, hostname, port, username, ssh_key_path, created_at, updated_at;

-- name: DeleteHost :exec
DELETE FROM hosts WHERE id = ?;
