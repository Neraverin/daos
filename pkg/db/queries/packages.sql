-- name: GetAllPackages :many
SELECT id, name, created_at, updated_at FROM packages ORDER BY name;

-- name: GetPackage :one
SELECT id, name, compose_content, created_at, updated_at FROM packages WHERE id = ?;

-- name: CreatePackage :one
INSERT INTO packages (name, compose_content)
VALUES (?, ?)
RETURNING id, name, compose_content, created_at, updated_at;

-- name: UpdatePackage :one
UPDATE packages SET name = ?, compose_content = ?, updated_at = datetime('now')
WHERE id = ?
RETURNING id, name, compose_content, created_at, updated_at;

-- name: DeletePackage :exec
DELETE FROM packages WHERE id = ?;
