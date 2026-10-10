package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tristan-renaud/chirpy/internal/auth"
	"github.com/tristan-renaud/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing path", err)
		return
	}

	chirp, err := cfg.db.RetrieveChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	if chirp.UserID == userID {
		err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
			UserID: userID,
			ID:     id,
		})
		if err != nil {
			respondWithError(w, 500, "could not delete chirp", err)
			return
		}
	} else {
		respondWithError(w, 403, "Unauthorized", err)
		return
	}

	w.WriteHeader(204)
}
