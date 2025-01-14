CREATE TABLE commands (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT NOT NULL,
  command       TEXT NOT NULL
);

CREATE TABLE projects (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT NOT NULL,
  port          INT NOT NULL,
  dir           TEXT
);

CREATE TABLE variables (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT NOT NULL,
  value         TEXT NOT NULL,

  project_id    INTEGER NOT NULL,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE domain_aliases (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  value         TEXT NOT NULL,

  project_id    INTEGER NOT NULL,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE project_commands (
  project_id    INTEGER NOT NULL,
  command_id    INTEGER NOT NULL,

  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,

  PRIMARY KEY (command_id, project_id)
);
