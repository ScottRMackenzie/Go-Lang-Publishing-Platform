package users

import (
	"context"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
)

func GetAll(ctx context.Context) ([]types.User, error) {
	rows, err := db.Pool.Query(ctx, "SELECT id, username, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		// user.CreatedAtStr = user.CreatedAt.Format("2006-01-02 15:04:05")
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetByID(id string, ctx context.Context) (types.User, error) {
	var user types.User
	err := db.Pool.QueryRow(ctx, "SELECT id, username FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func GetByUsername(username string, ctx context.Context) (types.User, error) {
	var user types.User
	err := db.Pool.QueryRow(ctx, "SELECT id, username FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}
