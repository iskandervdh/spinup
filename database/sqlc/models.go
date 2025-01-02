// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql"
)

type Command struct {
	ID      int64
	Name    string
	Command string
}

type DomainAlias struct {
	ID        int64
	Value     string
	ProjectID int64
}

type Project struct {
	ID     int64
	Name   string
	Domain string
	Port   int64
	Dir    sql.NullString
}

type ProjectCommand struct {
	ProjectID int64
	CommandID int64
}

type Variable struct {
	ID        int64
	Name      string
	Value     string
	ProjectID int64
}
