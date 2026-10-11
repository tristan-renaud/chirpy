package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/tristan-renaud/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	APIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "", err)
		return
	}

	if APIKey != cfg.polkaKey {
		respondWithError(w, 401, "", err)
	}

	type dataStruct struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string     `json:"event"`
		Data  dataStruct `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing userID", err)
	}

	err = cfg.db.UpgradeUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "user not found", err)
	}

	w.WriteHeader(204)
}
