package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/tbirddv/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	tokenSecret    string
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) writeMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", c.fileserverHits.Load())))
}

func (c *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	if c.platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	c.fileserverHits.Store(0)
	c.dbQueries.DeleteUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))
}
