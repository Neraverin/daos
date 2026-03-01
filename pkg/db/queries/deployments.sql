-- name: GetAllDeployments :many
SELECT 
    d.id, 
    d.host_id, 
    d.package_id, 
    d.status, 
    d.created_at, 
    d.updated_at,
    h.name as host_name,
    h.hostname as host_hostname,
    p.name as package_name
FROM deployments d
JOIN hosts h ON d.host_id = h.id
JOIN packages p ON d.package_id = p.id
ORDER BY d.created_at DESC;

-- name: GetDeployment :one
SELECT 
    d.id, 
    d.host_id, 
    d.package_id, 
    d.status, 
    d.created_at, 
    d.updated_at,
    h.name as host_name,
    h.hostname as host_hostname,
    p.name as package_name
FROM deployments d
JOIN hosts h ON d.host_id = h.id
JOIN packages p ON d.package_id = p.id
WHERE d.id = ?;

-- name: CreateDeployment :one
INSERT INTO deployments (host_id, package_id, status)
VALUES (?, ?, 'pending')
RETURNING id, host_id, package_id, status, created_at, updated_at;

-- name: UpdateDeploymentStatus :one
UPDATE deployments SET status = ?, updated_at = datetime('now')
WHERE id = ?
RETURNING id, host_id, package_id, status, created_at, updated_at;

-- name: DeleteDeployment :exec
DELETE FROM deployments WHERE id = ?;
