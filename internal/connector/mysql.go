package connector

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlConnector struct {
	cfg ConnectionConfig
	db  *sql.DB
}

func newMySQL(cfg ConnectionConfig) Connector {
	return &mysqlConnector{cfg: cfg}
}

func (c *mysqlConnector) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.cfg.Username, c.cfg.Password, c.cfg.Host, c.cfg.Port, c.cfg.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return err
	}
	c.db = db
	return nil
}

func (c *mysqlConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *mysqlConnector) FetchTables(ctx context.Context) ([]RawTable, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = ?
		  AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`, c.cfg.Database)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []RawTable
	for rows.Next() {
		var t RawTable
		if err := rows.Scan(&t.Name); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (c *mysqlConnector) FetchColumns(ctx context.Context, table string) ([]RawColumn, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			column_name,
			data_type,
			is_nullable = 'YES',
			column_default,
			character_maximum_length,
			ordinal_position,
			column_key = 'PRI'
		FROM information_schema.columns
		WHERE table_schema = ?
		  AND table_name = ?
		ORDER BY ordinal_position
	`, c.cfg.Database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []RawColumn
	for rows.Next() {
		var col RawColumn
		col.TableName = table
		if err := rows.Scan(
			&col.Name, &col.DataType, &col.IsNullable,
			&col.Default, &col.CharMaxLength, &col.OrdinalPosition, &col.IsPrimaryKey,
		); err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, rows.Err()
}

func (c *mysqlConnector) FetchForeignKeys(ctx context.Context) ([]RawForeignKey, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			table_name,
			column_name,
			referenced_table_name,
			referenced_column_name
		FROM information_schema.key_column_usage
		WHERE table_schema = ?
		  AND referenced_table_name IS NOT NULL
		ORDER BY table_name, column_name
	`, c.cfg.Database)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []RawForeignKey
	for rows.Next() {
		var fk RawForeignKey
		if err := rows.Scan(&fk.TableName, &fk.ColumnName, &fk.ForeignTable, &fk.ForeignColumn); err != nil {
			return nil, err
		}
		fks = append(fks, fk)
	}
	return fks, rows.Err()
}
