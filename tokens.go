package main

import (
	"net/http"
	"time"

	"github.com/BlackDogJet/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Body != http.NoBody {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	type returnVal struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating new JWT", err)
		return
	}

	responseWithJSON(w, http.StatusOK, returnVal{
		Token: accessToken,
	})
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Body != http.NoBody {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
