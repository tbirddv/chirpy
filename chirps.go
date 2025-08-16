package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
		respondWithError(w, "Failed to decode chirp params", http.StatusBadRequest)
		log.Printf("Error decoding chirp params: %v", err)
		return
	}

	if !validateLength(w, chirpParams) {
		return
	}
	badWordsSlice := []string{"kerfuffle", "sharbert", "fornax"}
	chirpParams.Body = cleanProfanity(chirpParams.Body, badWordsSlice)

	userID, err := c.getLoggedInUser(r)
	if err != nil {
		respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createParams := database.CreateChirpParams{
		Body:   chirpParams.Body,
		UserID: userID,
	}

	chirp, err := c.dbQueries.CreateChirp(r.Context(), createParams)
	if err != nil {
		respondWithError(w, "Failed to create chirp", http.StatusInternalServerError)
		log.Printf("Error creating chirp: %v", err)
		return
	}
	JSONChirp, err := createResponseStruct(chirp)
	if err != nil {
		respondWithError(w, "Failed to create chirp response", http.StatusInternalServerError)
		log.Printf("Error creating chirp response: %v", err)
		return
	}

	respondWithJSON(w, JSONChirp, http.StatusCreated)
}

func (c *apiConfig) getChirpsByUser(w http.ResponseWriter, r *http.Request, id string) {
	authorID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, "Invalid author ID", http.StatusBadRequest)
		log.Printf("Error parsing author ID: %v", err)
		return
	}
	sort := r.URL.Query().Get("sort")
	user, err := c.dbQueries.GetUserByID(r.Context(), authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, "User not found", http.StatusNotFound)
			return
		}
		respondWithError(w, "Failed to retrieve user", http.StatusInternalServerError)
		log.Printf("Error retrieving user: %v", err)
		return
	}
	var chirps []database.Chirp
	switch sort {
	case "":
		fallthrough
	case "asc":
		chirps, err = c.dbQueries.GetChirpsByUser(r.Context(), authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, fmt.Sprintf("No chirps found for user %s", user.Email), http.StatusNotFound)
				return
			}
			respondWithError(w, "Failed to retrieve chirps", http.StatusInternalServerError)
			log.Printf("Error retrieving chirps: %v", err)
			return
		}
	case "desc":
		chirps, err = c.dbQueries.GetChirpsByUserDesc(r.Context(), authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, fmt.Sprintf("No chirps found for user %s", user.Email), http.StatusNotFound)
				return
			}
			respondWithError(w, "Failed to retrieve chirps", http.StatusInternalServerError)
			log.Printf("Error retrieving chirps: %v", err)
			return
		}
	default:
		respondWithError(w, "Invalid sort query", http.StatusBadRequest)
		return
	}
	JSONChirps, err := createResponseStruct(chirps)
	if err != nil {
		respondWithError(w, "Failed to create chirp response", http.StatusInternalServerError)
		log.Printf("Error creating chirp response: %v", err)
		return
	}
	respondWithJSON(w, JSONChirps, http.StatusOK)
}

func (c *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	if id != "" {
		c.getChirpsByUser(w, r, id)
		return
	}

	var chirps []database.Chirp
	var err error
	switch r.URL.Query().Get("sort") {
	case "":
		fallthrough
	case "asc":

		chirps, err = c.dbQueries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, "Failed to retrieve chirps", http.StatusInternalServerError)
			log.Printf("Error retrieving chirps: %v", err)
			return
		}
	case "desc":
		chirps, err = c.dbQueries.GetChirpsDesc(r.Context())
		if err != nil {
			respondWithError(w, "Failed to retrieve chirps", http.StatusInternalServerError)
			log.Printf("Error retrieving chirps: %v", err)
			return
		}

	default:
		respondWithError(w, "Invalid sort query", http.StatusBadRequest)
	}

	chirpList, err := createResponseStruct(chirps)
	if err != nil {
		respondWithError(w, "Failed to create chirp response", http.StatusInternalServerError)
		log.Printf("Error creating chirp response: %v", err)
		return
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
		log.Printf("Error retrieving chirp: %v", err)
		return
	}
	JSONChirp, err := createResponseStruct(chirp)
	if err != nil {
		respondWithError(w, "Failed to create chirp response", http.StatusInternalServerError)
		log.Printf("Error creating chirp response: %v", err)
		return
	}
	respondWithJSON(w, JSONChirp, http.StatusOK)
}

func (c *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	userID, err := c.getLoggedInUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	idString := r.PathValue("id")
	chirpID, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}
	ChirpData, err := c.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Chirp not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve chirp", http.StatusInternalServerError)
		log.Printf("Error retrieving chirp: %v", err)
		return
	}

	if ChirpData.UserID != userID {
		http.Error(w, "Forbidden: You can only delete your own chirps", http.StatusForbidden)
		return
	}

	err = c.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		http.Error(w, "Failed to delete chirp", http.StatusInternalServerError)
		log.Printf("Error deleting chirp: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
