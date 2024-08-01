package users

import (
	"context"
	"fmt"
	"log"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/google/uuid"
)

func Create(ctx context.Context, username, password string) error {
	id := uuid.New().String()

	_, err := db.Pool.Exec(ctx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", id, username, password)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Data inserted successfully")
	return nil
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
