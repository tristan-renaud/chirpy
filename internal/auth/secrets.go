package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		log.Printf("Error, no Authorization header found")
		err := errors.New("authorization header empty")
		return "", err
	}
	key := strings.TrimSpace(strings.TrimPrefix(auth, "ApiKey "))
	return key, nil
}
