package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	file := http.Dir(".")

	fileServer := http.FileServer(file)

	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}
}
