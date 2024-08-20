package users

import (
	"context"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
)

func UpdateBalance(username string, balance int, ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, "UPDATE users SET balance = $1 WHERE username = $2", balance, username)
	if err != nil {
		return err
	}

	return nil
}
