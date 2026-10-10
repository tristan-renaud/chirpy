package main

import (
	"net/http"
	"time"

	"github.com/tristan-renaud/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	const tokenExpiration = 3600 * time.Second

	type payload struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unauthorized", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "unauthorized", err)
		return
	}

	accessToken, err := auth.MakeJWT(userID, cfg.tokenSecret, tokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue generating token", err)
		return
	}

	respondWithJSON(w, 200, payload{
		Token: accessToken,
	})
}
