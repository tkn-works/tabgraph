package api

import (
	"database/sql"
	"io"
	"io/fs"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tabui "github.com/koki/tabgraph/ui"
)

func NewServer(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Get("/connections", listConnections(db))
		r.Post("/connections", createConnection(db))
		r.Delete("/connections/{id}", deleteConnection(db))
		r.Post("/connections/{id}/sync", syncConnection(db))

		r.Get("/connections/{id}/tables", listTables(db))
		r.Get("/connections/{id}/tables/{table}", getTableDetail(db))

		r.Put("/metadata", upsertMetadata(db))

		r.Get("/connections/{id}/search", searchHandler(db))
		r.Get("/connections/{id}/er", erHandler(db))
	})

	sub, err := fs.Sub(tabui.StaticFiles, "out")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(sub))
	r.Get("/*", spaHandler(sub, fileServer))

	return r
}

// Dynamic route segments to replace with "_" for static file lookup.
// RSC payload files and HTML are generated only for the placeholder "_".
var (
	reConnTable = regexp.MustCompile(`^(/connections/)[^/]+(/tables/)[^/]+((?:/.*)?)$`)
	reConn      = regexp.MustCompile(`^(/connections/)[^/]+((?:/.*)?)$`)
)

// pathToPlaceholder maps a real dynamic path to the pre-built "_" placeholder path.
//
//	/connections/real-id/tables/orders/__next._head.txt
//	→ /connections/_/tables/_/__next._head.txt
func pathToPlaceholder(path string) string {
	if m := reConnTable.FindStringSubmatch(path); m != nil {
		return m[1] + "_" + m[2] + "_" + m[3]
	}
	if m := reConn.FindStringSubmatch(path); m != nil {
		return m[1] + "_" + m[2]
	}
	return path
}

// spaHandler serves static files and falls back to the "_" placeholder for dynamic
// routes so that Next.js RSC payloads are served correctly without a runtime server.
func spaHandler(sub fs.FS, fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Serve exact file if it exists.
		if tryServe(sub, fileServer, w, r, path) {
			return
		}

		// Map dynamic segments to placeholder and try again.
		ph := pathToPlaceholder(path)
		if ph != path && tryServe(sub, fileServer, w, r, ph) {
			return
		}

		// Final SPA fallback: root index.html.
		index, err := sub.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer index.Close()
		fi, _ := index.Stat()
		http.ServeContent(w, r, "index.html", fi.ModTime(), index.(io.ReadSeeker))
	}
}

// tryServe attempts to serve path as a file or as a directory index.
func tryServe(sub fs.FS, fileServer http.Handler, w http.ResponseWriter, r *http.Request, path string) bool {
	rel := strings.TrimPrefix(path, "/")
	if rel == "" {
		rel = "."
	}

	if f, err := sub.Open(rel); err == nil {
		fi, _ := f.Stat()
		f.Close()
		if !fi.IsDir() {
			r2 := requestWithPath(r, path)
			fileServer.ServeHTTP(w, r2)
			return true
		}
	}

	// Try directory index.
	if f, err := sub.Open(rel + "/index.html"); err == nil {
		f.Close()
		r2 := requestWithPath(r, path+"/index.html")
		fileServer.ServeHTTP(w, r2)
		return true
	}

	return false
}

func requestWithPath(r *http.Request, path string) *http.Request {
	r2 := r.Clone(r.Context())
	r2.URL = r.URL.ResolveReference(r.URL)
	r2.URL.Path = path
	return r2
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
