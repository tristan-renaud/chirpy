package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tristan-renaud/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters %s", err)
		respondWithError(w, 400, "Error decoding paramters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("error creating user: %s", err)
		respondWithError(w, 500, "Error creating user", err)
		return
	}

	u := database.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 201, u)
}
