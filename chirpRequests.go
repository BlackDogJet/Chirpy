package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON request body")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp body exceeds maximum length of 140 characters")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanedBody := getCleanedBody(params.Body, badWords)

	responseWithJSON(w, http.StatusOK, returnVal{CleanedBody: cleanedBody})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	cleanedBody := body
	for badWord := range badWords {
		cleanedBody = replaceIgnoreCase(cleanedBody, badWord, "****")
	}

	return cleanedBody
}

func replaceIgnoreCase(cleanedBody, badWord, s string) string {
	lowerBody := strings.ToLower(cleanedBody)
	lowerBadWord := strings.ToLower(badWord)

	for {
		index := strings.Index(lowerBody, lowerBadWord)
		if index == -1 {
			break
		}

		cleanedBody = cleanedBody[:index] + s + cleanedBody[index+len(badWord):]
		lowerBody = strings.ToLower(cleanedBody)
	}

	return cleanedBody
}
