package main

import (
	"log"
	"net/http"
)

func main() {
	config := &apiConfig{}
	handler := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	handler.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	handler.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	handler.HandleFunc("GET /metrics", config.writeMetrics)
	handler.HandleFunc("POST /reset", config.resetMetrics)

	log.Fatal(server.ListenAndServe())
}
