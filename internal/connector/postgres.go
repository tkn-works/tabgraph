package connector

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type postgresConnector struct {
	cfg ConnectionConfig
	db  *sql.DB
}

func newPostgres(cfg ConnectionConfig) Connector {
	return &postgresConnector{cfg: cfg}
}

func (c *postgresConnector) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.cfg.Host, c.cfg.Port, c.cfg.Database, c.cfg.Username, c.cfg.Password, c.cfg.SSLMode)
	db, err := sql.Open("postgres", dsn)
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

func (c *postgresConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *postgresConnector) FetchTables(ctx context.Context) ([]RawTable, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		  AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`)
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

func (c *postgresConnector) FetchColumns(ctx context.Context, table string) ([]RawColumn, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable = 'YES',
			c.column_default,
			c.character_maximum_length,
			c.ordinal_position,
			EXISTS (
				SELECT 1 FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu
					ON tc.constraint_name = kcu.constraint_name
					AND tc.table_schema = kcu.table_schema
				WHERE tc.constraint_type = 'PRIMARY KEY'
				  AND tc.table_name = c.table_name
				  AND kcu.column_name = c.column_name
			) AS is_primary_key
		FROM information_schema.columns c
		WHERE c.table_schema = 'public'
		  AND c.table_name = $1
		ORDER BY c.ordinal_position
	`, table)
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

func (c *postgresConnector) FetchForeignKeys(ctx context.Context) ([]RawForeignKey, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			kcu.table_name,
			kcu.column_name,
			ccu.table_name AS foreign_table,
			ccu.column_name AS foreign_column
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
		  AND tc.table_schema = 'public'
		ORDER BY kcu.table_name, kcu.column_name
	`)
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
