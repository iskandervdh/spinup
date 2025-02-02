// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: projects.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (
  name, port
) VALUES (
  ?, ?
)
RETURNING id, name, port, dir
`

type CreateProjectParams struct {
	Name string
	Port int64
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject, arg.Name, arg.Port)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Port,
		&i.Dir,
	)
	return i, err
}

const createProjectCommand = `-- name: CreateProjectCommand :exec
INSERT INTO project_commands (
  command_id, project_id
) VALUES (
  ?, ?
)
`

type CreateProjectCommandParams struct {
	CommandID int64
	ProjectID int64
}

func (q *Queries) CreateProjectCommand(ctx context.Context, arg CreateProjectCommandParams) error {
	_, err := q.db.ExecContext(ctx, createProjectCommand, arg.CommandID, arg.ProjectID)
	return err
}

const deleteCommandsProjects = `-- name: DeleteCommandsProjects :exec

DELETE FROM project_commands
WHERE project_id = ? AND command_id = ?
`

type DeleteCommandsProjectsParams struct {
	ProjectID int64
	CommandID int64
}

// -- name: CreateCommandsProjects :exec
// INSERT INTO project_commands (
//
//	command_id, project_id
//
// ) VALUES (
//
//	VALUES(?),
//	VALUES(?)
//
// );
func (q *Queries) DeleteCommandsProjects(ctx context.Context, arg DeleteCommandsProjectsParams) error {
	_, err := q.db.ExecContext(ctx, deleteCommandsProjects, arg.ProjectID, arg.CommandID)
	return err
}

const deleteProject = `-- name: DeleteProject :exec
DELETE FROM projects
WHERE name = ?
`

func (q *Queries) DeleteProject(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, deleteProject, name)
	return err
}

const deleteProjectById = `-- name: DeleteProjectById :exec
DELETE FROM projects
WHERE id = ?
`

func (q *Queries) DeleteProjectById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProjectById, id)
	return err
}

const getProject = `-- name: GetProject :one
SELECT id, name, port, dir
FROM projects
WHERE name = ? LIMIT 1
`

func (q *Queries) GetProject(ctx context.Context, name string) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProject, name)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Port,
		&i.Dir,
	)
	return i, err
}

const getProjectCommands = `-- name: GetProjectCommands :many
SELECT c.id, c.name, c.command
FROM commands c
JOIN project_commands cp ON c.id = cp.command_id
WHERE cp.project_id = ?
`

func (q *Queries) GetProjectCommands(ctx context.Context, projectID int64) ([]Command, error) {
	rows, err := q.db.QueryContext(ctx, getProjectCommands, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Command
	for rows.Next() {
		var i Command
		if err := rows.Scan(&i.ID, &i.Name, &i.Command); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProjectDomainAliases = `-- name: GetProjectDomainAliases :many
SELECT id, value, project_id
FROM domain_aliases
WHERE project_id = ?
`

func (q *Queries) GetProjectDomainAliases(ctx context.Context, projectID int64) ([]DomainAlias, error) {
	rows, err := q.db.QueryContext(ctx, getProjectDomainAliases, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DomainAlias
	for rows.Next() {
		var i DomainAlias
		if err := rows.Scan(&i.ID, &i.Value, &i.ProjectID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProjectVariables = `-- name: GetProjectVariables :many
SELECT id, name, value, project_id
FROM variables
WHERE project_id = ?
`

func (q *Queries) GetProjectVariables(ctx context.Context, projectID int64) ([]Variable, error) {
	rows, err := q.db.QueryContext(ctx, getProjectVariables, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Variable
	for rows.Next() {
		var i Variable
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Value,
			&i.ProjectID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProjects = `-- name: GetProjects :many
SELECT id, name, port, dir
FROM projects
`

func (q *Queries) GetProjects(ctx context.Context) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, getProjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Port,
			&i.Dir,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const renameProject = `-- name: RenameProject :exec
UPDATE projects
SET name = ?
WHERE name = ?
`

type RenameProjectParams struct {
	Name   string
	Name_2 string
}

func (q *Queries) RenameProject(ctx context.Context, arg RenameProjectParams) error {
	_, err := q.db.ExecContext(ctx, renameProject, arg.Name, arg.Name_2)
	return err
}

const setProjectDir = `-- name: SetProjectDir :exec
UPDATE projects
SET dir = ?
WHERE id = ?
`

type SetProjectDirParams struct {
	Dir sql.NullString
	ID  int64
}

func (q *Queries) SetProjectDir(ctx context.Context, arg SetProjectDirParams) error {
	_, err := q.db.ExecContext(ctx, setProjectDir, arg.Dir, arg.ID)
	return err
}

const updateProject = `-- name: UpdateProject :exec
UPDATE projects
SET port = ?, dir = ?
WHERE name = ?
`

type UpdateProjectParams struct {
	Port int64
	Dir  sql.NullString
	Name string
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) error {
	_, err := q.db.ExecContext(ctx, updateProject, arg.Port, arg.Dir, arg.Name)
	return err
}

const updateProjectById = `-- name: UpdateProjectById :exec
UPDATE projects
SET name = ?, port = ?, dir = ?
WHERE id = ?
`

type UpdateProjectByIdParams struct {
	Name string
	Port int64
	Dir  sql.NullString
	ID   int64
}

func (q *Queries) UpdateProjectById(ctx context.Context, arg UpdateProjectByIdParams) error {
	_, err := q.db.ExecContext(ctx, updateProjectById,
		arg.Name,
		arg.Port,
		arg.Dir,
		arg.ID,
	)
	return err
}
