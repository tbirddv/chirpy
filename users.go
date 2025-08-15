package main

import (
	"encoding/json"
	"net/http"

	"github.com/tbirddv/chirpy/internal/auth"
	"github.com/tbirddv/chirpy/internal/database"
)

func (c *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var p userParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createUserParams := database.CreateUserParams{
		Email:          p.Email,
		HashedPassword: hashedPassword,
	}
	user, err := c.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSONUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, JSONUser, http.StatusCreated)
}
