-- name: GetAllRoles :many
SELECT id, name, created_at, updated_at FROM roles ORDER BY name;

-- name: GetRole :one
SELECT id, name, compose_content, created_at, updated_at FROM roles WHERE id = ?;

-- name: CreateRole :one
INSERT INTO roles (id, name, compose_content)
VALUES (?, ?, ?)
RETURNING id, name, compose_content, created_at, updated_at;

-- name: UpdateRole :one
UPDATE roles SET name = ?, compose_content = ?, updated_at = datetime('now')
WHERE id = ?
RETURNING id, name, compose_content, created_at, updated_at;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = ?;
