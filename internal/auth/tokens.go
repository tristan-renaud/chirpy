package auth

import (
	"log"
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
	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		log.Printf("error signing string %s", err)
		return "", err
	}
	return tokenString, nil
}
