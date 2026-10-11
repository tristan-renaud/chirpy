package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, r *http.Request) {
	var jsonChirps []Chirp
	userID := r.URL.Query().Get("author_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt parse userID", err)
	}

	if userID == "" {
		chirps, err := cfg.db.RetrieveChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database error: %s", err)
			return
		}
		for _, chirp := range chirps {
			c := Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
			jsonChirps = append(jsonChirps, c)
		}
	} else {
		chirps, err := cfg.db.RetrieveAuthorsChirps(r.Context(), userUUID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "", err)
		}
		for _, chirp := range chirps {
			c := Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
			jsonChirps = append(jsonChirps, c)
		}
	}

	respondWithJSON(w, 200, jsonChirps)
}
