package users

import (
	"context"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(username, password string) (types.User, error) {
	hashedPassword, err := GetHashedPasswordByUsername(username, context.Background())
	if err != nil {
		return types.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return types.User{}, err
	}

	return GetByUsername(username, context.Background())
}
