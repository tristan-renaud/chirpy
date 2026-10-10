package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/tristan-renaud/chirpy/internal/auth"
	"github.com/tristan-renaud/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	const jwtExpiration = 3600 * time.Second
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

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		log.Printf("Password did not match %s", err)
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret, jwtExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to created token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	err = cfg.db.StoreToken(r.Context(), database.StoreTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to store token in database", err)
		return
	}

	u := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, 200, u)
}
