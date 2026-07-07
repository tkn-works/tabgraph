package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/koki/tabgraph/internal/db"
)

func upsertMetadata(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input db.MetadataInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		m, err := db.UpsertMetadata(sqlDB, input)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOK(w, m)
	}
}
