package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) requestCounter(w http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) resetCounter(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	var cfg apiConfig

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /metrics", cfg.requestCounter)
	mux.HandleFunc("POST /reset", cfg.resetCounter)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
