-- name: GetCommand :one
SELECT *
FROM commands
WHERE name = ? LIMIT 1;

-- name: GetCommands :many
SELECT *
FROM commands;

-- name: CreateCommand :exec
INSERT INTO commands (
  name, command
) VALUES (
  ?, ?
);

-- name: UpdateCommand :exec
UPDATE commands
SET command = ?
WHERE name = ?;

-- name: UpdateCommandById :exec
UPDATE commands
SET name = ?, command = ?
WHERE id = ?;

-- name: RenameCommand :exec
UPDATE commands
SET name = ?
WHERE name = ?;

-- name: DeleteCommand :exec
DELETE FROM commands
WHERE name = ?;

-- name: DeleteCommandById :exec
DELETE FROM commands
WHERE id = ?;
