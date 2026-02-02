package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BlackDogJet/Chirpy/internal/auth"
	"github.com/BlackDogJet/Chirpy/internal/databases"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVal struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	params.Password = hashedPassword

	user, err := cfg.db.CreateUser(r.Context(), databases.CreateUserParams{
		Email:          params.Email,
		HashedPassword: params.Password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user in database", err)
		return
	}

	responseWithJSON(w, http.StatusCreated, returnVal{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
