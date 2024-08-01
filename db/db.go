package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(connStr string) {
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func Ping() error {
	return Pool.Ping(context.Background())
}

func CreateTable(ctx context.Context) {
	_, err := Pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created successfully")
}

func InsertData(ctx context.Context, username, password string) {
	_, err := Pool.Exec(ctx, "INSERT INTO users (id, name, password) VALUES ($1, $2)", username, password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data inserted successfully")
}

func QueryData(ctx context.Context) {
	rows, err := Pool.Query(ctx, "SELECT id, username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, username, password string
		if err := rows.Scan(&id, &username, &password); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", id, username, password)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// func GetUsers(ctx context.Context) [][]string {
// 	rows, err := Pool.Query(ctx, "SELECT id, username FROM users")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	var users [][]string
// 	for rows.Next() {
// 		var id, username string
// 		if err := rows.Scan(&id, &username); err != nil {
// 			log.Fatal(err)
// 		}
// 		users = append(users, []string{id, username})
// 	}
// 	return users
// }
