package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, r *http.Request) {
	var jsonChirps []Chirp
	sortQ := r.URL.Query().Get("sort")
	userID := r.URL.Query().Get("author_id")

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
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldnt parse userID", err)
			return
		}
		chirps, err := cfg.db.RetrieveAuthorsChirps(r.Context(), userUUID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "", err)
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
	}
	if sortQ == "desc" {
		sort.Slice(jsonChirps, func(i, j int) bool {
			return jsonChirps[i].CreatedAt.After(jsonChirps[j].CreatedAt)
		})
	}
	respondWithJSON(w, 200, jsonChirps)
}
