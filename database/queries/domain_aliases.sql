-- name: GetDomainAliases :many
SELECT value
FROM domain_aliases
WHERE project_id = ?;

-- name: CreateDomainAlias :exec
INSERT INTO domain_aliases (
  value, project_id
) VALUES (
  ?, ?
);

-- name: UpdateDomainAlias :exec
UPDATE domain_aliases
SET value = ?
WHERE value = ? AND project_id = ?;

-- name: DeleteDomainAlias :exec
DELETE FROM domain_aliases
WHERE value = ? AND project_id = ?;
