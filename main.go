package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/api"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/joho/godotenv"
)

func main() {
	apiPort := 2337
	staticPort := 80
	var wg sync.WaitGroup

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	db.Connect(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	defer db.Pool.Close()
	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to ping the database: %v", err)
		return
	}

	// API server
	wg.Add(1)
	go func() {
		defer wg.Done()

		mux := http.NewServeMux()

		mux.HandleFunc("POST /api", api.WelcomeHandler)
		mux.HandleFunc("POST /api/users/create", api.CreateUserHandler)

		// !! Dangerous code
		mux.HandleFunc("POST /api/users", api.GetUsersHandler)
		mux.HandleFunc("POST /api/users/{id}", api.GetUserByIDHandler)

		apiServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", apiPort),
			Handler: mux,
		}
		fmt.Printf("API server started on :%d\n", apiPort)
		if err := apiServer.ListenAndServe(); err != nil {
			fmt.Printf("API server error: %v\n", err)
		}
	}()

	// Static file server
	wg.Add(1)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", staticHandle(http.FileServer(http.Dir("./static"))))
		defer wg.Done()
		staticServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", staticPort),
			Handler: mux,
		}
		fmt.Printf("Static file server started on :%d\n", staticPort)
		if err := staticServer.ListenAndServe(); err != nil {
			fmt.Printf("Static file server error: %v\n", err)
		}
	}()

	wg.Wait()
}

func staticHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
