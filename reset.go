package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.PLATFORM == "dev" {
		err := cfg.db.DeleteUsers(r.Context())
		if err != nil {
			log.Printf("error deleting users database: %s", err)
			respondWithError(w, 500, "error deleting users database", err)
			return
		}
		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0, users database deleted"))
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}
}
