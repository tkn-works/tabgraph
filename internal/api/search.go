package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/koki/tabgraph/internal/db"
)

func searchHandler(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		q := r.URL.Query().Get("q")
		if q == "" {
			jsonOK(w, []db.SearchResult{})
			return
		}
		results, err := db.SearchMetadata(sqlDB, id, q)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOK(w, results)
	}
}
