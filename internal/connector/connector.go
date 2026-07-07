package connector

import (
	"context"
	"fmt"
)

type Connector interface {
	Connect(ctx context.Context) error
	Close() error
	FetchTables(ctx context.Context) ([]RawTable, error)
	FetchColumns(ctx context.Context, table string) ([]RawColumn, error)
	FetchForeignKeys(ctx context.Context) ([]RawForeignKey, error)
}

type RawTable struct {
	Name string
}

type RawColumn struct {
	TableName       string
	Name            string
	DataType        string
	IsNullable      bool
	Default         *string
	CharMaxLength   *int
	IsPrimaryKey    bool
	OrdinalPosition int
}

type RawForeignKey struct {
	TableName     string
	ColumnName    string
	ForeignTable  string
	ForeignColumn string
}

type ConnectionConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
}

func New(cfg ConnectionConfig) (Connector, error) {
	switch cfg.Driver {
	case "postgres":
		return newPostgres(cfg), nil
	case "mysql":
		return newMySQL(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}
}
