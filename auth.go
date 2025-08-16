package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tbirddv/chirpy/internal/auth"
	"github.com/tbirddv/chirpy/internal/database"
)

func (c *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var p userParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.dbQueries.GetUserByEmail(r.Context(), p.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = auth.CheckPasswordHash(p.Password, user.HashedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, c.tokenSecret, time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 days
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, loginResponse, http.StatusOK)
}

func (c *apiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenData, err := c.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tokenData.RevokedAt.Valid || tokenData.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Refresh token expired or revoked", http.StatusUnauthorized)
		return
	}

	user, err := c.dbQueries.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, c.tokenSecret, time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := refreshResponse{
		Token: newAccessToken,
	}

	respondWithJSON(w, response, http.StatusOK)
}

func (c *apiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = c.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
