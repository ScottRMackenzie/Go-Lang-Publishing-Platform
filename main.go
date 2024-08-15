package main

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
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
	var wg sync.WaitGroup

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	base_url_frontend := os.Getenv("BASE_URL_FRONTEND")
	base_url_api := os.Getenv("BASE_URL_API")
	apiPort := os.Getenv("API_PORT")
	staticPort := os.Getenv("FRONTEND_PORT")

	if base_url_frontend == "" || base_url_api == "" || apiPort == "" || staticPort == "" {
		fmt.Println("BASE_URL_FRONTEND, BASE_URL_API, API_PORT, and FRONTEND_PORT must be set in the .env file")
		return
	}

	// Database
	if os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_NAME") == "" {
		fmt.Println("DB_USER, DB_PASSWORD, and DB_NAME must be set in the .env file")
		return
	}

	db.Connect(fmt.Sprintf("user=%s password=%s dbname=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
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
		mux.Handle("/api/users", middleware.CombinedAuthMiddleware(http.HandlerFunc(api.GetUsersHandler)))
		mux.Handle("/api/users/create", middleware.JWTMiddleware(http.HandlerFunc(api.CreateUserHandler)))
		mux.Handle("/api/users/verify-email/{token}", middleware.JWTMiddleware(http.HandlerFunc(email_verification.VerifyEmailHandler)))
		mux.Handle("/api/users/{id}", middleware.JWTMiddleware(http.HandlerFunc(api.GetUserByIDHandler)))

		mux.HandleFunc("/api/v1/public/books", api.GetAllBooksHandler)
		mux.HandleFunc("/api/v1/public/books/{id}", api.GetBookByIDHandler)
		mux.HandleFunc("/api/v1/public/books/sorted", api.GetSortedBooksHandler)
		mux.HandleFunc("/api/v1/public/books/search", api.FilteredSearchBooksHandler)

		middlewareHandler := corsMiddleware(LoggingMiddleware(mux))

		apiServer := &http.Server{
			Addr:    fmt.Sprintf(":%s", apiPort),
			Handler: middlewareHandler,
		}
		fmt.Printf("API server started on %s:%s\n", base_url_api, apiPort)
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
		mux.HandleFunc("GET /static/books/covers/{id}", controller.ServeBookCover)

		mux.Handle("GET /", middleware.AuthMiddleware(http.HandlerFunc(controller.HomeHandler)))
		mux.Handle("GET /login", middleware.AuthMiddleware(http.HandlerFunc(controller.LoginHandler)))
		mux.Handle("POST /login", middleware.AuthMiddleware(http.HandlerFunc(controller.LoginHandler)))
		mux.Handle("GET /logout", middleware.AuthMiddleware(http.HandlerFunc(controller.LogoutHandler)))

		mux.Handle("GET /browse", middleware.AuthMiddleware(http.HandlerFunc(controller.BrowseHandler)))

		middlewareHandler := corsMiddleware(LoggingMiddleware(FrontendMiddleware(mux)))

		frontendServer := &http.Server{
			Addr:    fmt.Sprintf(":%s", staticPort),
			Handler: middlewareHandler,
		}
		fmt.Printf("Template and static file server started on %s:%s\n", base_url_frontend, staticPort)
		if err := frontendServer.ListenAndServe(); err != nil {
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
		w.Header().Set("Access-Control-Allow-Origin", "http://tb-books.local")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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

func FrontendMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		parts := strings.Split(host, ".")

		// Check if the host has more than two parts (indicating a subdomain)
		if len(parts) > 2 {
			// Skip processing for subdomains
			http.Error(w, "Subdomain requests are not allowed for port 80", http.StatusForbidden)
			return
		}

		// Call the next handler if it's not a subdomain
		next.ServeHTTP(w, r)
	})
}
