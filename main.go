package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tbirddv/chirpy/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	config := &apiConfig{dbQueries: database.New(db), platform: platform}

	handler := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	handler.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	handler.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	handler.HandleFunc("GET /admin/metrics", config.writeMetrics)
	handler.HandleFunc("POST /admin/reset", config.resetMetrics)
	handler.HandleFunc("POST /api/validate_chirp", validateChirp)
	handler.HandleFunc("POST /api/users", config.createUser)

	log.Fatal(server.ListenAndServe())
}
