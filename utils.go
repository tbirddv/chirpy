package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tbirddv/chirpy/internal/auth"
	"github.com/tbirddv/chirpy/internal/database"
)

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(ValidationError{Error: message})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}

}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}

}

func cleanProfanity(body string, badWords []string) string {
	badWordsMap := make(map[string]struct{})

	for _, word := range badWords {
		badWordsMap[strings.ToLower(word)] = struct{}{}
	}

	splitBody := strings.Split(body, " ")
	cleanedBody := make([]string, 0, len(splitBody))
	for _, word := range splitBody {
		if _, found := badWordsMap[strings.ToLower(word)]; found {
			cleanedBody = append(cleanedBody, "****")
		} else {
			cleanedBody = append(cleanedBody, word)
		}
	}
	return strings.Join(cleanedBody, " ")
}

func (c *apiConfig) getLoggedInUser(r *http.Request) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, err
	}
	return auth.ValidateJWT(token, c.tokenSecret)
}

func createResponseStruct(input interface{}) (any, error) {
	switch v := input.(type) {
	case database.User:
		return User{
			ID:          v.ID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Email:       v.Email,
			IsChirpyRed: v.IsChirpyRed,
		}, nil
	case database.Chirp:
		return Chirp{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserID:    v.UserID,
		}, nil
	case []database.Chirp:
		var chirps []Chirp
		for _, c := range v {
			chirp, err := createResponseStruct(c)
			if err != nil {
				return nil, err
			}
			chirps = append(chirps, chirp.(Chirp))
		}
		return chirps, nil
	default:
		return nil, fmt.Errorf("unknown type: %T", input)
	}
}
