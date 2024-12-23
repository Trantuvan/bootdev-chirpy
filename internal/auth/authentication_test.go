package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	t.Run("Valid JWT", func(t *testing.T) {
		secretKey1 := "TestValidateJWT"
		expectUUID := uuid.UUID{}
		expectedToken, _ := MakeJWT(expectUUID, secretKey1, time.Second*2)
		actualUUID, err := ValidateJWT(expectedToken, secretKey1)

		if actualUUID.String() != expectUUID.String() {
			t.Errorf("ValidateJWT() error = %v", err)
		}
	})
	t.Run("InValid Secret key", func(t *testing.T) {
		secretKey1 := "TestValidateJWT"
		expectUUID := uuid.UUID{}
		expiredTime := time.Duration(1 * time.Second)
		expectedToken, _ := MakeJWT(expectUUID, "", expiredTime)
		actualUUID, err := ValidateJWT(expectedToken, secretKey1)

		if err == nil {
			t.Errorf("ValidateJWT() expected err Invalid SecretKey but given userID %v", actualUUID.String())
		}
	})
}
