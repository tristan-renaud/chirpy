package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokens(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "bigJuicyB00ties"
	tokenString, err := MakeJWT(userID, tokenSecret, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	validID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatal(err)
	}
	if validID != userID {
		t.Errorf("%v and %v != match", userID, validID)
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "wutAHooT!"
	tokenString, err := MakeJWT(userID, tokenSecret, -1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("Expected error here but did not receive one")
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "iluvmyDAWG"
	tokenString, err := MakeJWT(userID, tokenSecret, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ValidateJWT(tokenString, "dolaisbestCAT")
	if err == nil {
		t.Fatal("Ecpected error here but did not receive one")
	}
}
