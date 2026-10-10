package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	timeUTC := time.Now()
	timeDuration := time.Time.Add(timeUTC, expiresIn)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(timeUTC.UTC()),
		ExpiresAt: jwt.NewNumericDate(timeDuration.UTC()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("error signing string %s", err)
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("error: %s", err)
		return uuid.Nil, err
	}
	return uuid.Parse(claims.Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		log.Printf("Error, no Authorization header found")
		err := errors.New("authorization header empty")
		return "", err
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	return token, nil
}

func MakeRefreshToken() string {
	randomData := make([]byte, 32)
	rand.Read(randomData)
	return hex.EncodeToString(randomData)
}
