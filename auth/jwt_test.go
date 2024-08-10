package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	// Set up the environment variable for testing
	os.Setenv("JWT_SECRET", "testsecret")

	// Test data
	userID := "123"
	username := "testuser"

	// Generate JWT
	token, err := GenerateJWT(userID, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token
	claims, err := VerifyJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
}

func TestVerifyJWT_InvalidToken(t *testing.T) {
	// Set up the environment variable for testing
	os.Setenv("JWT_SECRET", "testsecret")

	// Test with an invalid token
	invalidToken := "invalid.token.string"

	claims, err := VerifyJWT(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateJWT_WithoutEnv(t *testing.T) {
	// Clear the environment variable
	os.Unsetenv("JWT_SECRET")

	// Test data
	userID := "123"
	username := "testuser"

	// Generate JWT
	token, err := GenerateJWT(userID, username)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestVerifyJWT_WithoutEnv(t *testing.T) {
	// Clear the environment variable
	os.Unsetenv("JWT_SECRET")

	// Test with a valid token (assuming you have a valid token)
	//validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjZhYTVkZjctZmQ0NS00ZDcxLTk4MzEtMzNkODliNjgxYmU0IiwidXNlcm5hbWUiOiJpYW10YWxsIiwiZXhwIjoxNzIzMjMwNjMyfQ.A6fWj5sLQCT6Q_1KNCfLc9sDv-2l56OElIr8mkAws7g" // Replace with a valid token
	validToken := "x"

	claims, err := VerifyJWT(validToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
