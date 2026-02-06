package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	// Frist, we need to create a user ID and some tokens to test with
	user_id := uuid.New()
	token_secret := "the_right_secret"
	valid_token, _ := MakeJWT(user_id, token_secret, time.Hour*24)
	invalid_token, _ := MakeJWT(user_id, "the_wrong_secret", time.Hour*24)
	expired_token, _ := MakeJWT(user_id, token_secret, time.Nanosecond*1)
	time.Sleep(time.Millisecond * 10) // Sleep for a bit so expired_token expires
	tests := []struct {
		name     string
		token    string
		wantErr  bool
		expected uuid.UUID
	}{
		{
			name:     "Correct token",
			token:    valid_token,
			wantErr:  false,
			expected: user_id,
		},
		{
			name:     "Incorrect token",
			token:    invalid_token,
			wantErr:  true,
			expected: uuid.UUID{},
		},
		{
			name:     "Expired token",
			token:    expired_token,
			wantErr:  true,
			expected: uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ValidateJWT(tt.token, token_secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if actual != tt.expected {
				t.Errorf("ValidateJWT() expected: %v, got: %v", tt.expected, actual)
			}
		})
	}
}
