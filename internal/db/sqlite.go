package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite は並列書き込み不可
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

var migrations = []string{
	`PRAGMA journal_mode=WAL`,
	`PRAGMA foreign_keys=ON`,
	`CREATE TABLE IF NOT EXISTS connections (
		id             TEXT    PRIMARY KEY,
		name           TEXT    NOT NULL,
		driver         TEXT    NOT NULL,
		host           TEXT    NOT NULL,
		port           INTEGER NOT NULL,
		database       TEXT    NOT NULL,
		username       TEXT    NOT NULL,
		password       TEXT    NOT NULL,
		ssl_mode       TEXT    NOT NULL DEFAULT 'disable',
		created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_synced_at DATETIME
	)`,
	`CREATE TABLE IF NOT EXISTS schema_cache (
		id               TEXT    PRIMARY KEY,
		connection_id    TEXT    NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
		table_name       TEXT    NOT NULL,
		column_name      TEXT    NOT NULL,
		data_type        TEXT    NOT NULL,
		is_nullable      BOOLEAN NOT NULL DEFAULT 1,
		column_default   TEXT,
		char_max_length  INTEGER,
		is_primary_key   BOOLEAN NOT NULL DEFAULT 0,
		is_foreign_key   BOOLEAN NOT NULL DEFAULT 0,
		fk_table         TEXT,
		fk_column        TEXT,
		ordinal_position INTEGER NOT NULL,
		synced_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE INDEX IF NOT EXISTS idx_sc_conn  ON schema_cache(connection_id)`,
	`CREATE INDEX IF NOT EXISTS idx_sc_table ON schema_cache(connection_id, table_name)`,
	// column_name uses '' as sentinel for table-level descriptions (NULL breaks UNIQUE)
	`CREATE TABLE IF NOT EXISTS metadata (
		id            TEXT    PRIMARY KEY,
		connection_id TEXT    NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
		table_name    TEXT    NOT NULL,
		column_name   TEXT    NOT NULL DEFAULT '',
		description   TEXT    NOT NULL DEFAULT '',
		updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(connection_id, table_name, column_name)
	)`,
	`CREATE VIRTUAL TABLE IF NOT EXISTS search_index USING fts5(
		connection_id UNINDEXED,
		table_name,
		column_name,
		description,
		content='metadata',
		content_rowid='rowid',
		tokenize='trigram'
	)`,
	`CREATE TRIGGER IF NOT EXISTS metadata_ai AFTER INSERT ON metadata BEGIN
		INSERT INTO search_index(rowid, connection_id, table_name, column_name, description)
		VALUES (new.rowid, new.connection_id, new.table_name, new.column_name, new.description);
	END`,
	`CREATE TRIGGER IF NOT EXISTS metadata_ad AFTER DELETE ON metadata BEGIN
		INSERT INTO search_index(search_index, rowid, connection_id, table_name, column_name, description)
		VALUES ('delete', old.rowid, old.connection_id, old.table_name, old.column_name, old.description);
	END`,
	`CREATE TRIGGER IF NOT EXISTS metadata_au AFTER UPDATE ON metadata BEGIN
		INSERT INTO search_index(search_index, rowid, connection_id, table_name, column_name, description)
		VALUES ('delete', old.rowid, old.connection_id, old.table_name, old.column_name, old.description);
		INSERT INTO search_index(rowid, connection_id, table_name, column_name, description)
		VALUES (new.rowid, new.connection_id, new.table_name, new.column_name, new.description);
	END`,
}

func migrate(db *sql.DB) error {
	for _, stmt := range migrations {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
