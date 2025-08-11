package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) writeMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d\n", c.fileserverHits.Load())))
}

func (c *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))
}
