package users

import (
	"context"
	"log"
	"time"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
	"github.com/google/uuid"
)

func Create(ctx context.Context, username, email, password string) (types.User, error) {
	id := uuid.New().String()
	created_at := time.Now()

	_, err := db.Pool.Exec(ctx, "INSERT INTO users (id, username, email, password, created_at) VALUES ($1, $2, $3, $4, $5)", id, username, email, password, created_at)
	if err != nil {
		log.Fatal(err)
		return types.User{}, err
	}

	return types.User{ID: id, Username: username, Email: email, CreatedAt: created_at}, nil
}

func CheckUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return exists, nil
}

func CheckEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return exists, nil
}
