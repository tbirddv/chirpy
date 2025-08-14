package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Chirp struct {
	Body string `json:"body"`
}

type ValidationError struct {
	Error string `json:"error"`
}

type ValidResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateLength(w http.ResponseWriter, chirp Chirp) bool {
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

func validateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var chirp Chirp
	if err := json.NewDecoder(r.Body).Decode(&chirp); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !validateLength(w, chirp) {
		return
	}
	badWordsSlice := []string{"kerfuffle", "sharbert", "fornax"}
	chirp.Body = cleanProfanity(chirp.Body, badWordsSlice)
	respondWithJSON(w, ValidResponse{CleanedBody: chirp.Body}, http.StatusOK)
}
