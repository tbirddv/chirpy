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
	tokenSecret := os.Getenv("TOKENSECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	config := &apiConfig{dbQueries: database.New(db), platform: platform, tokenSecret: tokenSecret, polkaKey: polkaKey}

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
	handler.HandleFunc("POST /api/chirps", config.CreateChirp)
	handler.HandleFunc("GET /api/chirps", config.GetChirps)
	handler.HandleFunc("GET /api/chirps/{id}", config.GetChirpByID)
	handler.HandleFunc("POST /api/users", config.createUser)
	handler.HandleFunc("POST /api/login", config.HandleLogin)
	handler.HandleFunc("POST /api/refresh", config.HandleRefresh)
	handler.HandleFunc("POST /api/revoke", config.HandleRevoke)
	handler.HandleFunc("PUT /api/users", config.updateUser)
	handler.HandleFunc("DELETE /api/chirps/{id}", config.DeleteChirp)
	handler.HandleFunc("POST /api/polka/webhooks", config.GiveChirpyRed)

	log.Fatal(server.ListenAndServe())
}
