package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/tristan-renaud/chirpy/internal/database"

	_ "github.com/lib/pq"
)

func cleanChirps(payloadBody string) any {
	type cleanChirp struct {
		Body string `json:"cleaned_body"`
	}

	words := strings.Split(payloadBody, " ")
	var cleanWords []string

	for _, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			cleanWords = append(cleanWords, "****")
		default:
			cleanWords = append(cleanWords, word)
		}
	}

	clean := cleanChirp{
		Body: strings.Join(cleanWords, " "),
	}
	return clean
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Printf("Error code %d: %s", code, msg)
	type errorResponse struct {
		Error string `json:"error"`
	}
	err := errorResponse{
		Error: msg,
	}

	respondWithJSON(w, code, err)
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	PLATFORM       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) requestCounter(w http.ResponseWriter, req *http.Request) {
	content := fmt.Sprintf("<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

func (cfg *apiConfig) resetCounter(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	// loading .env file for postgres URL
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	// opening connection to that database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("%s", err)
	}

	const filepathRoot = "."
	const port = "8080"

	var cfg apiConfig
	cfg.db = database.New(db)

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", cfg.requestCounter)

	mux.HandleFunc("POST /admin/reset", cfg.resetCounter)

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, req *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters %s", err)
			respondWithError(w, 400, "Error decoding paramters")
			return
		}

		if len(params.Body) > 140 {
			log.Printf("Length of body exceeded 140 characters")
			respondWithError(w, 400, "Length of body exceeded 140 characters")
			return
		}

		cleanResp := cleanChirps(params.Body)

		respondWithJSON(w, 200, cleanResp)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
