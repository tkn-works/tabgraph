package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	ID           string    `json:"id"`
	ConnectionID string    `json:"connectionId"`
	TableName    string    `json:"tableName"`
	ColumnName   *string   `json:"columnName"`
	Description  string    `json:"description"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type MetadataInput struct {
	ConnectionID string  `json:"connectionId"`
	TableName    string  `json:"tableName"`
	ColumnName   *string `json:"columnName"` // nil = table-level description
	Description  string  `json:"description"`
}

type SearchResult struct {
	TableName   string  `json:"tableName"`
	ColumnName  *string `json:"columnName"`
	Description string  `json:"description"`
	MatchType   string  `json:"matchType"`
}

// colKey converts ColumnName pointer to sentinel string for storage.
// nil (table description) → ""
// non-nil → column name
func colKey(columnName *string) string {
	if columnName == nil {
		return ""
	}
	return *columnName
}

func UpsertMetadata(db *sql.DB, input MetadataInput) (*Metadata, error) {
	id := uuid.New().String()
	key := colKey(input.ColumnName)

	_, err := db.Exec(`
		INSERT INTO metadata (id, connection_id, table_name, column_name, description, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(connection_id, table_name, column_name)
		DO UPDATE SET description = excluded.description, updated_at = CURRENT_TIMESTAMP
	`, id, input.ConnectionID, input.TableName, key, input.Description)
	if err != nil {
		return nil, err
	}

	var m Metadata
	var colName string
	err = db.QueryRow(`
		SELECT id, connection_id, table_name, column_name, description, updated_at
		FROM metadata
		WHERE connection_id = ? AND table_name = ? AND column_name = ?
	`, input.ConnectionID, input.TableName, key).Scan(
		&m.ID, &m.ConnectionID, &m.TableName, &colName, &m.Description, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if colName != "" {
		m.ColumnName = &colName
	}
	return &m, nil
}

// SeedMetadata inserts empty metadata rows for all tables and columns after sync,
// so that table/column names are always searchable even without descriptions.
// Uses INSERT OR IGNORE to preserve existing descriptions.
func SeedMetadata(sqlDB *sql.DB, connectionID string, cacheRows []SchemaCacheRow) error {
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tables := make(map[string]bool)
	for _, row := range cacheRows {
		tables[row.TableName] = true
	}
	for tableName := range tables {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO metadata (id, connection_id, table_name, column_name, description, updated_at)
			VALUES (?, ?, ?, '', '', CURRENT_TIMESTAMP)
		`, uuid.New().String(), connectionID, tableName)
		if err != nil {
			return err
		}
	}
	for _, row := range cacheRows {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO metadata (id, connection_id, table_name, column_name, description, updated_at)
			VALUES (?, ?, ?, ?, '', CURRENT_TIMESTAMP)
		`, uuid.New().String(), connectionID, row.TableName, row.ColumnName)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// fts5Phrase wraps a user query in double-quotes for FTS5 phrase matching,
// escaping any embedded quotes by doubling them. This prevents FTS5 syntax
// errors from special characters like *, ", (, ), -, ^.
func fts5Phrase(q string) string {
	return `"` + strings.ReplaceAll(q, `"`, `""`) + `"`
}

func SearchMetadata(db *sql.DB, connectionID, query string) ([]SearchResult, error) {
	rows, err := db.Query(`
		SELECT table_name, column_name, description
		FROM search_index
		WHERE connection_id = ? AND (search_index MATCH ? OR table_name LIKE ? OR column_name LIKE ?)
		ORDER BY rank
		LIMIT 50
	`, connectionID, fts5Phrase(query), "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		var colName string
		if err := rows.Scan(&r.TableName, &colName, &r.Description); err != nil {
			return nil, err
		}
		if colName != "" {
			r.ColumnName = &colName
			r.MatchType = "column"
		} else {
			r.MatchType = "table"
		}
		results = append(results, r)
	}
	if results == nil {
		results = []SearchResult{}
	}
	return results, rows.Err()
}
