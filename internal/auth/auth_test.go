package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("Password hash mismatch: %v", err)
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	id, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Failed to validate JWT: %v", err)
	}

	if id != userID {
		t.Errorf("Expected user ID %v, got %v", userID, id)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	id, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Failed to validate JWT: %v", err)
	}

	if id != userID {
		t.Errorf("Expected user ID %v, got %v", userID, id)
	}
}

func TestTimeOutJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := -time.Second * 5

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	time.Sleep(6 * time.Second) // Wait for the token to expire

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("Expected token to be expired")
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Errorf("Expected token to be invalid")
	}
}
