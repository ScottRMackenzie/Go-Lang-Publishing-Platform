package main

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"sync"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/api"
	email_verification "github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/api/verification"
	controller "github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/controllers"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/middleware"
	"github.com/joho/godotenv"
)

func init() {
	// Register the JavaScript MIME type
	mime.AddExtensionType(".js", "application/javascript")
}

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

		mux.HandleFunc("/api/login", api.LoginHandler)
		mux.Handle("/api/users", middleware.JWTMiddleware(http.HandlerFunc(api.GetUsersHandler)))
		mux.Handle("/api/users/create", middleware.JWTMiddleware(http.HandlerFunc(api.CreateUserHandler)))
		mux.Handle("/api/users/verify-email/{token}", middleware.JWTMiddleware(http.HandlerFunc(email_verification.VerifyEmailHandler)))
		mux.Handle("/api/users/{id}", middleware.JWTMiddleware(http.HandlerFunc(api.GetUserByIDHandler)))

		mux.HandleFunc("/api/v1/public/books", api.GetAllBooksHandler)
		mux.HandleFunc("/api/v1/public/books/{id}", api.GetBookByIDHandler)
		mux.HandleFunc("/api/v1/public/books/sorted", api.GetSortedBooksHandler)
		mux.HandleFunc("/api/v1/public/books/search", api.FilteredSearchBooksHandler)

		middlewareHandler := corsMiddleware(LoggingMiddleware(mux))

		apiServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", apiPort),
			Handler: middlewareHandler,
		}
		fmt.Printf("API server started on :%d\n", apiPort)
		if err := apiServer.ListenAndServe(); err != nil {
			fmt.Printf("API server error: %v\n", err)
		}
	}()

	// Template server
	wg.Add(1)
	go func() {
		mux := http.NewServeMux()
		defer wg.Done()

		mux.Handle("GET /static/", staticHandle(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
		mux.Handle("GET /favicon.ico", staticHandle(http.FileServer(http.Dir("static"))))
		// mux.HandleFunc("GET /static/js/", serveJavaScript)
		mux.HandleFunc("GET /static/books/covers/{id}", controller.ServeBookCover)

		mux.Handle("GET /", middleware.AuthMiddleware(http.HandlerFunc(controller.HomeHandler)))
		mux.Handle("GET /login", middleware.AuthMiddleware(http.HandlerFunc(controller.LoginHandler)))
		mux.Handle("POST /login", middleware.AuthMiddleware(http.HandlerFunc(controller.LoginHandler)))
		mux.Handle("GET /logout", middleware.AuthMiddleware(http.HandlerFunc(controller.LogoutHandler)))

		mux.Handle("GET /browse", middleware.AuthMiddleware(http.HandlerFunc(controller.BrowseHandler)))

		mux.Handle("GET /client-api/{path}", middleware.AuthMiddleware(http.HandlerFunc(controller.ClientAPIHandler)))
		mux.Handle("GET /check-authenticated", middleware.AuthMiddleware(http.HandlerFunc(controller.CheckIfAuthenticatedHandler)))

		middlewareHandler := corsMiddleware(LoggingMiddleware(mux))

		staticServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", staticPort),
			Handler: middlewareHandler,
		}
		fmt.Printf("Template and static file server started on :%d\n", staticPort)
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// print each request
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
