package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Driver       string     `json:"driver"`
	Host         string     `json:"host"`
	Port         int        `json:"port"`
	Database     string     `json:"database"`
	Username     string     `json:"username"`
	SSLMode      string     `json:"sslMode"`
	CreatedAt    time.Time  `json:"createdAt"`
	LastSyncedAt *time.Time `json:"lastSyncedAt"`
}

type ConnectionInput struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"sslMode"`
}

func ListConnections(db *sql.DB) ([]Connection, error) {
	rows, err := db.Query(`
		SELECT id, name, driver, host, port, database, username, ssl_mode, created_at, last_synced_at
		FROM connections ORDER BY created_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conns []Connection
	for rows.Next() {
		var c Connection
		if err := rows.Scan(&c.ID, &c.Name, &c.Driver, &c.Host, &c.Port, &c.Database,
			&c.Username, &c.SSLMode, &c.CreatedAt, &c.LastSyncedAt); err != nil {
			return nil, err
		}
		conns = append(conns, c)
	}
	if conns == nil {
		conns = []Connection{}
	}
	return conns, rows.Err()
}

func CreateConnection(db *sql.DB, input ConnectionInput) (*Connection, error) {
	id := uuid.New().String()
	sslMode := input.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	_, err := db.Exec(`
		INSERT INTO connections (id, name, driver, host, port, database, username, password, ssl_mode)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, input.Name, input.Driver, input.Host, input.Port, input.Database,
		input.Username, input.Password, sslMode)
	if err != nil {
		return nil, err
	}
	return GetConnection(db, id)
}

func GetConnection(db *sql.DB, id string) (*Connection, error) {
	var c Connection
	err := db.QueryRow(`
		SELECT id, name, driver, host, port, database, username, ssl_mode, created_at, last_synced_at
		FROM connections WHERE id = ?
	`, id).Scan(&c.ID, &c.Name, &c.Driver, &c.Host, &c.Port, &c.Database,
		&c.Username, &c.SSLMode, &c.CreatedAt, &c.LastSyncedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func GetConnectionWithPassword(db *sql.DB, id string) (*ConnectionInput, error) {
	var c ConnectionInput
	var connID string
	err := db.QueryRow(`
		SELECT id, name, driver, host, port, database, username, password, ssl_mode
		FROM connections WHERE id = ?
	`, id).Scan(&connID, &c.Name, &c.Driver, &c.Host, &c.Port, &c.Database,
		&c.Username, &c.Password, &c.SSLMode)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func DeleteConnection(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM connections WHERE id = ?`, id)
	return err
}

func UpdateLastSynced(db *sql.DB, id string) error {
	_, err := db.Exec(`UPDATE connections SET last_synced_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}
