package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BlackDogJet/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) webhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve API Key", err)
		return
	}

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "API Key is invalid", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON request body", err)
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "No content", nil)
		return
	}

	_, err = cfg.db.UpdateUserIsChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}

		respondWithError(w, http.StatusNotFound, "Error updating user to Chirpy Red in database", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
