package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
