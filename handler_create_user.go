package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters %s", err)
		respondWithError(w, 400, "Error decoding paramters")
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("error creating user: %s", err)
		respondWithError(w, 500, "Error creating user")
		return
	}

	u := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 201, u)
}
