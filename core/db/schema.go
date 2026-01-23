package db

const schema = `
PRAGMA foreign_keys = ON;

-- Content-addressed storage
CREATE TABLE IF NOT EXISTS objects (
	hash TEXT PRIMARY KEY,
	content BLOB NOT NULL,
	created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Commits (snapshots)
CREATE TABLE IF NOT EXISTS snapshot (
	hash TEXT PRIMARY KEY,
	parent_hash TEXT,
	message TEXT NOT NULL,
	author TEXT,
	timestamp INTEGER DEFAULT (strftime('%s', 'now')),
	FOREIGN KEY(parent_hash) REFERENCES snapshot(hash) ON DELETE SET NULL
);

-- Files inside a snapshot
CREATE TABLE IF NOT EXISTS snapshot_files (
	snapshot_hash TEXT NOT NULL,
	path TEXT NOT NULL CHECK (path != ''),
	object_hash TEXT NOT NULL,
	PRIMARY KEY (snapshot_hash, path),
	FOREIGN KEY(snapshot_hash) REFERENCES snapshot(hash) ON DELETE CASCADE,
	FOREIGN KEY(object_hash) REFERENCES objects(hash) ON DELETE CASCADE
);

-- refs
CREATE TABLE IF NOT EXISTS refs (
	name TEXT PRIMARY KEY CHECK (name != ''),
	snapshot_hash TEXT,
	updated_at INTEGER DEFAULT (strftime('%s', 'now')),
	FOREIGN KEY(snapshot_hash) REFERENCES snapshot(hash) ON DELETE SET NULL
);

-- Config
CREATE TABLE IF NOT EXISTS config (
	key TEXT PRIMARY KEY CHECK (key != ''),
	value TEXT NOT NULL
);

-- Staging = intent only
CREATE TABLE IF NOT EXISTS staged_files (
	path TEXT PRIMARY KEY CHECK (path != ''),
	staged_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Defaults
INSERT OR IGNORE INTO config (key, value) VALUES ('head', 'main');
`
