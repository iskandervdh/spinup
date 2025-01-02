-- name: GetProject :one
SELECT id, name, domain, port, dir
FROM projects
WHERE name = ? LIMIT 1;

-- name: GetProjects :many
SELECT *
FROM projects;

-- name: GetProjectCommands :many
SELECT c.*
FROM commands c
JOIN project_commands cp ON c.id = cp.command_id
WHERE cp.project_id = ?;

-- name: GetProjectVariables :many
SELECT *
FROM variables
WHERE project_id = ?;

-- name: GetProjectDomainAliases :many
SELECT *
FROM domain_aliases
WHERE project_id = ?;

-- name: CreateProject :one
INSERT INTO projects (
  name, domain, port
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: SetProjectDir :exec
UPDATE projects
SET dir = ?
WHERE id = ?;

-- name: CreateProjectCommand :exec
INSERT INTO project_commands (
  command_id, project_id
) VALUES (
  ?, ?
);

-- -- name: CreateCommandsProjects :exec
-- INSERT INTO project_commands (
--   command_id, project_id
-- ) VALUES (
--   VALUES(?),
--   VALUES(?)
-- );

-- name: DeleteCommandsProjects :exec
DELETE FROM project_commands
WHERE project_id = ? AND command_id = ?;

-- name: UpdateProject :exec
UPDATE projects
SET domain = ?, port = ?, dir = ?
WHERE name = ?;

-- name: RenameProject :exec
UPDATE projects
SET name = ?
WHERE name = ?;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE name = ?;
