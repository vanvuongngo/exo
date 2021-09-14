package sqlite

const ddlStatements = `
CREATE TABLE IF NOT EXISTS workspace (
  id TEXT PRIMARY KEY,
	root TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS workspace_root ON workspace ( root );

CREATE TABLE IF NOT EXISTS component (
	id TEXT PRIMARY KEY,
	parent_id TEXT,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	definition_version INTEGER NOT NULL,
	state_version INTEGER NOT NULL,
	created INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS component_name ON component ( parent_id, name );

CREATE TABLE component_definition
	component_id TEXT NOT NULL,
	version INTEGER NOT NULL,
	parent_id TEXT NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	timestamp INTEGER NOT NULL,
	spec TEXT NOT NULL,
	meta TEXT NOT NULL,

	PRIMARY KEY (component_id, version)
);

CREATE TABLE IF NOT EXISTS component_state (
	component_id TEXT NOT NULL,
	version INTEGER NOT NULL,
	timestamp INTEGER NOT NULL,
	value TEXT NOT NULL,

	PRIMARY KEY (component_id, version)
);
`
