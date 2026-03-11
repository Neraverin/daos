-- name: GetAllDeployments :many
SELECT 
    d.id, 
    d.host_id, 
    d.role_id, 
    d.status, 
    d.created_at, 
    d.updated_at,
    h.name as host_name,
    h.hostname as host_hostname,
    r.name as role_name
FROM deployments d
JOIN hosts h ON d.host_id = h.id
JOIN roles r ON d.role_id = r.id
ORDER BY d.created_at DESC;

-- name: GetDeployment :one
SELECT 
    d.id, 
    d.host_id, 
    d.role_id, 
    d.status, 
    d.created_at, 
    d.updated_at,
    h.name as host_name,
    h.hostname as host_hostname,
    r.name as role_name
FROM deployments d
JOIN hosts h ON d.host_id = h.id
JOIN roles r ON d.role_id = r.id
WHERE d.id = ?;

-- name: CreateDeployment :one
INSERT INTO deployments (id, host_id, role_id, status)
VALUES (?, ?, ?, 'pending')
RETURNING id, host_id, role_id, status, created_at, updated_at;

-- name: UpdateDeploymentStatus :one
UPDATE deployments SET status = ?, updated_at = datetime('now')
WHERE id = ?
RETURNING id, host_id, role_id, status, created_at, updated_at;

-- name: DeleteDeployment :exec
DELETE FROM deployments WHERE id = ?;

-- name: GetDeploymentDetails :one
SELECT 
    d.host_id,
    d.role_id,
    h.hostname,
    h.username,
    h.ssh_key_path,
    r.compose_content
FROM deployments d
JOIN hosts h ON d.host_id = h.id
JOIN roles r ON d.role_id = r.id
WHERE d.id = ?;
