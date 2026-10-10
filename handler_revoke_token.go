package main

import (
	"net/http"

	"github.com/tristan-renaud/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unauthorized", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke token", err)
		return
	}

	w.WriteHeader(204)
}
