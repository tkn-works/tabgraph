package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/koki/tabgraph/internal/db"
)

func listTables(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		tables, err := db.ListTables(sqlDB, id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOK(w, tables)
	}
}

type tableDetail struct {
	Table   db.TableInfo   `json:"table"`
	Columns []db.ColumnInfo `json:"columns"`
}

func getTableDetail(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		tableName := chi.URLParam(r, "table")

		tables, err := db.ListTables(sqlDB, id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var tableInfo db.TableInfo
		for _, t := range tables {
			if t.Name == tableName {
				tableInfo = t
				break
			}
		}

		cols, err := db.GetTableColumns(sqlDB, id, tableName)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonOK(w, tableDetail{Table: tableInfo, Columns: cols})
	}
}
