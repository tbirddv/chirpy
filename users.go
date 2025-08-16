package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tbirddv/chirpy/internal/auth"
	"github.com/tbirddv/chirpy/internal/database"
)

func (c *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var p userParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		log.Printf("Error decoding user params: %v", err)
		return
	}
	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}
	createUserParams := database.CreateUserParams{
		Email:          p.Email,
		HashedPassword: hashedPassword,
	}
	user, err := c.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	JSONUser, err := createResponseStruct(user)
	if err != nil {
		http.Error(w, "Failed to create user response", http.StatusInternalServerError)
		log.Printf("Error creating user response: %v", err)
		return
	}

	respondWithJSON(w, JSONUser, http.StatusCreated)
}

func (c *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	UserID, err := c.getLoggedInUser(r)
	if err != nil {
		http.Error(w, "Failed to get logged in user", http.StatusUnauthorized)
		log.Printf("Error getting logged in user: %v", err)
		return
	}

	var p userParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to update user", http.StatusBadRequest)
		log.Printf("Error decoding user params: %v", err)
		return
	}

	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}

	updateUserParams := database.UpdateUserParams{
		ID:             UserID,
		Email:          p.Email,
		HashedPassword: hashedPassword,
	}

	user, err := c.dbQueries.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("Error updating user: %v", err)
		return
	}

	JSONUser, err := createResponseStruct(user)
	if err != nil {
		http.Error(w, "Failed to create user response", http.StatusInternalServerError)
		log.Printf("Error creating user response: %v", err)
		return
	}

	respondWithJSON(w, JSONUser, http.StatusOK)
}

func (c *apiConfig) GiveChirpyRed(w http.ResponseWriter, r *http.Request) {
	var eventData chirpyRedEvent
	if err := json.NewDecoder(r.Body).Decode(&eventData); err != nil {
		http.Error(w, "Failed to give Chirpy Red", http.StatusBadRequest)
		log.Printf("Error decoding Chirpy Red event: %v", err)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if apiKey != c.polkaKey {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	if eventData.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(eventData.Data.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Error parsing user ID: %v", err)
		return
	}

	_, err = c.dbQueries.GiveChirpyRed(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to give Chirpy Red", http.StatusInternalServerError)
		log.Printf("Error giving Chirpy Red: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
