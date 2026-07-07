package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TableInfo struct {
	Name        string `json:"name"`
	ColumnCount int    `json:"columnCount"`
	Description string `json:"description"`
}

type ColumnInfo struct {
	Name            string   `json:"name"`
	DataType        string   `json:"dataType"`
	IsNullable      bool     `json:"isNullable"`
	Default         *string  `json:"default"`
	IsPrimaryKey    bool     `json:"isPrimaryKey"`
	IsForeignKey    bool     `json:"isForeignKey"`
	ForeignTable    *string  `json:"foreignTable"`
	ForeignColumn   *string  `json:"foreignColumn"`
	Description     string   `json:"description"`
	OrdinalPosition int      `json:"ordinalPosition"`
}

type SchemaCacheRow struct {
	ConnectionID    string
	TableName       string
	ColumnName      string
	DataType        string
	IsNullable      bool
	ColumnDefault   *string
	CharMaxLength   *int
	IsPrimaryKey    bool
	IsForeignKey    bool
	FKTable         *string
	FKColumn        *string
	OrdinalPosition int
	SyncedAt        time.Time
}

func ReplaceSchemaCache(db *sql.DB, connectionID string, rows []SchemaCacheRow) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM schema_cache WHERE connection_id = ?`, connectionID); err != nil {
		return err
	}

	for _, row := range rows {
		id := uuid.New().String()
		if _, err := tx.Exec(`
			INSERT INTO schema_cache (
				id, connection_id, table_name, column_name, data_type,
				is_nullable, column_default, char_max_length,
				is_primary_key, is_foreign_key, fk_table, fk_column, ordinal_position
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, id, connectionID, row.TableName, row.ColumnName, row.DataType,
			row.IsNullable, row.ColumnDefault, row.CharMaxLength,
			row.IsPrimaryKey, row.IsForeignKey, row.FKTable, row.FKColumn, row.OrdinalPosition,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func ListTables(db *sql.DB, connectionID string) ([]TableInfo, error) {
	rows, err := db.Query(`
		SELECT sc.table_name, COUNT(*) AS col_count, COALESCE(m.description, '') AS description
		FROM schema_cache sc
		LEFT JOIN metadata m
			ON m.connection_id = sc.connection_id
			AND m.table_name = sc.table_name
			AND m.column_name = ''
		WHERE sc.connection_id = ?
		GROUP BY sc.table_name, m.description
		ORDER BY sc.table_name
	`, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.Name, &t.ColumnCount, &t.Description); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if tables == nil {
		tables = []TableInfo{}
	}
	return tables, rows.Err()
}

func GetTableColumns(db *sql.DB, connectionID, tableName string) ([]ColumnInfo, error) {
	rows, err := db.Query(`
		SELECT
			sc.column_name, sc.data_type, sc.is_nullable, sc.column_default,
			sc.is_primary_key, sc.is_foreign_key, sc.fk_table, sc.fk_column,
			sc.ordinal_position, COALESCE(m.description, '') AS description
		FROM schema_cache sc
		LEFT JOIN metadata m
			ON m.connection_id = sc.connection_id
			AND m.table_name = sc.table_name
			AND m.column_name = sc.column_name
		WHERE sc.connection_id = ? AND sc.table_name = ?
		ORDER BY sc.ordinal_position
	`, connectionID, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		if err := rows.Scan(
			&c.Name, &c.DataType, &c.IsNullable, &c.Default,
			&c.IsPrimaryKey, &c.IsForeignKey, &c.ForeignTable, &c.ForeignColumn,
			&c.OrdinalPosition, &c.Description,
		); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	if cols == nil {
		cols = []ColumnInfo{}
	}
	return cols, rows.Err()
}

type FKRow struct {
	TableName     string
	ColumnName    string
	ForeignTable  string
	ForeignColumn string
}

func GetForeignKeys(db *sql.DB, connectionID string) ([]FKRow, error) {
	rows, err := db.Query(`
		SELECT table_name, column_name, fk_table, fk_column
		FROM schema_cache
		WHERE connection_id = ? AND is_foreign_key = 1
		ORDER BY table_name, column_name
	`, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []FKRow
	for rows.Next() {
		var fk FKRow
		if err := rows.Scan(&fk.TableName, &fk.ColumnName, &fk.ForeignTable, &fk.ForeignColumn); err != nil {
			return nil, err
		}
		fks = append(fks, fk)
	}
	return fks, rows.Err()
}
