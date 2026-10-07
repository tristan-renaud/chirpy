package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tristan-renaud/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters %s", err)
		respondWithError(w, 400, "Error decoding paramters", err)
		return
	}

	user, err := cfg.db.RetrieveUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("email not found in database %s", err)
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	u := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		log.Printf("Password did not match %s", err)
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	respondWithJSON(w, 200, u)
}
