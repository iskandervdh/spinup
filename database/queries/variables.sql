-- name: GetVariables :many
SELECT name, value
FROM variables
WHERE project_id = ?;

-- name: CreateVariable :exec
INSERT INTO variables (
  name, value, project_id
) VALUES (
  ?, ?, ?
);

-- name: UpdateVariable :exec
UPDATE variables
SET value = ?
WHERE name = ? AND project_id = ?;

-- name: RenameVariable :exec
UPDATE variables
SET name = ?
WHERE name = ? AND project_id = ?;

-- name: DeleteVariable :exec
DELETE FROM variables
WHERE name = ? AND project_id = ?;
