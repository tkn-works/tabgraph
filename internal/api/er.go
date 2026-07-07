package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/koki/tabgraph/internal/db"
)

type erResponse struct {
	Diagram string `json:"diagram"`
}

func erHandler(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		tables, err := db.ListTables(sqlDB, id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tableColumns := make(map[string][]db.ColumnInfo)
		for _, t := range tables {
			cols, err := db.GetTableColumns(sqlDB, id, t.Name)
			if err != nil {
				jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tableColumns[t.Name] = cols
		}

		fks, err := db.GetForeignKeys(sqlDB, id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tableSet := make(map[string]bool)
		for _, t := range tables {
			tableSet[t.Name] = true
		}

		type relation struct{ from, to string }
		relSet := make(map[string]relation)

		for _, fk := range fks {
			key := fk.TableName + "." + fk.ColumnName
			relSet[key] = relation{from: fk.TableName, to: fk.ForeignTable}
		}

		// Heuristic: xxx_id → xxx / xxxs / xxxes
		for tName, cols := range tableColumns {
			for _, col := range cols {
				if col.IsForeignKey {
					continue
				}
				if strings.HasSuffix(col.Name, "_id") {
					prefix := strings.TrimSuffix(col.Name, "_id")
					for _, cand := range []string{prefix + "s", prefix + "es", prefix} {
						if tableSet[cand] {
							relSet[tName+"."+col.Name] = relation{from: tName, to: cand}
							break
						}
					}
				}
			}
		}

		var sb strings.Builder
		sb.WriteString("erDiagram\n")

		for _, t := range tables {
			sb.WriteString(fmt.Sprintf("    %s {\n", mermaidName(t.Name)))
			for _, col := range tableColumns[t.Name] {
				suffix := ""
				if col.IsPrimaryKey {
					suffix = " PK"
				} else if col.IsForeignKey {
					suffix = " FK"
				}
				sb.WriteString(fmt.Sprintf("        %s %s%s\n",
					mermaidType(col.DataType), mermaidName(col.Name), suffix))
			}
			sb.WriteString("    }\n")
		}

		seen := make(map[string]bool)
		for _, rel := range relSet {
			key := rel.to + "→" + rel.from
			if seen[key] {
				continue
			}
			seen[key] = true
			sb.WriteString(fmt.Sprintf("    %s ||--o{ %s : \"\"\n",
				mermaidName(rel.to), mermaidName(rel.from)))
		}

		jsonOK(w, erResponse{Diagram: sb.String()})
	}
}

func mermaidName(s string) string {
	return strings.NewReplacer("-", "_", " ", "_", ".", "_").Replace(s)
}

func mermaidType(dataType string) string {
	t := strings.ToLower(dataType)
	switch {
	case strings.Contains(t, "int"):
		return "int"
	case strings.Contains(t, "bool"):
		return "boolean"
	case strings.Contains(t, "float"), strings.Contains(t, "decimal"),
		strings.Contains(t, "numeric"), strings.Contains(t, "double"):
		return "float"
	case strings.Contains(t, "timestamp"), strings.Contains(t, "datetime"),
		strings.Contains(t, "date"):
		return "datetime"
	default:
		return "string"
	}
}
