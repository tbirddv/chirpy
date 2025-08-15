package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tbirddv/chirpy/internal/auth"
	"github.com/tbirddv/chirpy/internal/database"
)

func validateLength(w http.ResponseWriter, chirp chirpParams) bool {
	chirp.Body = strings.TrimSpace(chirp.Body)
	if len(chirp.Body) == 0 {
		respondWithError(w, "Chirp cannot be empty", http.StatusBadRequest)
		return false
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, "Chirp is too long", http.StatusBadRequest)
		return false
	}
	return true
}

func (c *apiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {

	var chirpParams chirpParams
	if err := json.NewDecoder(r.Body).Decode(&chirpParams); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !validateLength(w, chirpParams) {
		return
	}
	badWordsSlice := []string{"kerfuffle", "sharbert", "fornax"}
	chirpParams.Body = cleanProfanity(chirpParams.Body, badWordsSlice)

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, c.tokenSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createParams := database.CreateChirpParams{
		Body:   chirpParams.Body,
		UserID: userID,
	}

	chirp, err := c.dbQueries.CreateChirp(r.Context(), createParams)
	if err != nil {
		respondWithError(w, "Failed to create chirp", http.StatusInternalServerError)
		return
	}
	JSONChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, JSONChirp, http.StatusCreated)
}

func (c *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, "Failed to retrieve chirps", http.StatusInternalServerError)
		return
	}

	var chirpList []Chirp
	for _, chirp := range chirps {
		chirpList = append(chirpList, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, chirpList, http.StatusOK)
}

func (c *apiConfig) GetChirpByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := c.dbQueries.GetChirpByID(r.Context(), id)
	if err == sql.ErrNoRows {
		respondWithError(w, "Chirp not found", http.StatusNotFound)
		return
	}
	if err != nil {
		respondWithError(w, "Failed to retrieve chirp", http.StatusInternalServerError)
		return
	}
	JSONChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, JSONChirp, http.StatusOK)
}
