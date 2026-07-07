package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/koki/tabgraph/internal/connector"
	"github.com/koki/tabgraph/internal/db"
)

func listConnections(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conns, err := db.ListConnections(sqlDB)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOK(w, conns)
	}
}

func createConnection(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input db.ConnectionInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		conn, err := db.CreateConnection(sqlDB, input)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOK(w, conn)
	}
}

func deleteConnection(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := db.DeleteConnection(sqlDB, id); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type syncResult struct {
	TablesCount  int    `json:"tablesCount"`
	ColumnsCount int    `json:"columnsCount"`
	SyncedAt     string `json:"syncedAt"`
}

func syncConnection(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		connInfo, err := db.GetConnectionWithPassword(sqlDB, id)
		if connInfo == nil || err != nil {
			jsonError(w, "connection not found", http.StatusNotFound)
			return
		}

		conn, err := connector.New(connector.ConnectionConfig{
			Driver:   connInfo.Driver,
			Host:     connInfo.Host,
			Port:     connInfo.Port,
			Database: connInfo.Database,
			Username: connInfo.Username,
			Password: connInfo.Password,
			SSLMode:  connInfo.SSLMode,
		})
		if err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		if err := conn.Connect(ctx); err != nil {
			jsonError(w, "connection failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer conn.Close()

		tables, err := conn.FetchTables(ctx)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fks, err := conn.FetchForeignKeys(ctx)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fkMap := make(map[string]map[string]connector.RawForeignKey)
		for _, fk := range fks {
			if fkMap[fk.TableName] == nil {
				fkMap[fk.TableName] = make(map[string]connector.RawForeignKey)
			}
			fkMap[fk.TableName][fk.ColumnName] = fk
		}

		var cacheRows []db.SchemaCacheRow
		for _, t := range tables {
			cols, err := conn.FetchColumns(ctx, t.Name)
			if err != nil {
				jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, col := range cols {
				row := db.SchemaCacheRow{
					ConnectionID:    id,
					TableName:       t.Name,
					ColumnName:      col.Name,
					DataType:        col.DataType,
					IsNullable:      col.IsNullable,
					ColumnDefault:   col.Default,
					CharMaxLength:   col.CharMaxLength,
					IsPrimaryKey:    col.IsPrimaryKey,
					OrdinalPosition: col.OrdinalPosition,
				}
				if fk, ok := fkMap[t.Name][col.Name]; ok {
					row.IsForeignKey = true
					row.FKTable = &fk.ForeignTable
					row.FKColumn = &fk.ForeignColumn
				}
				cacheRows = append(cacheRows, row)
			}
		}

		if err := db.ReplaceSchemaCache(sqlDB, id, cacheRows); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := db.SeedMetadata(sqlDB, id, cacheRows); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := db.UpdateLastSynced(sqlDB, id); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonOK(w, syncResult{
			TablesCount:  len(tables),
			ColumnsCount: len(cacheRows),
		})
	}
}
