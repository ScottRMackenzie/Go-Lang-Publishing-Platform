package main

import (
	"fmt"
	"mime"
	"net/http"
	"os"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/api"
	email_verification "github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/api/verification"
	controller "github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/controllers"
	bj_controller "github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/controllers/games/blackjack"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/middleware"
	"github.com/joho/godotenv"
)

func init() {
	// Register the JavaScript MIME type
	mime.AddExtensionType(".js", "application/javascript")
}

func main() {
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

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", api.LoginHandler)
	mux.Handle("POST /api/users", middleware.CombinedAuthMiddleware(http.HandlerFunc(api.GetUsersHandler)))
	mux.Handle("POST /api/users/create", middleware.JWTMiddleware(http.HandlerFunc(api.CreateUserHandler)))
	mux.Handle("GET /api/users/verify-email/{token}", middleware.JWTMiddleware(http.HandlerFunc(email_verification.VerifyEmailHandler)))
	mux.Handle("GET /api/users/{id}", middleware.JWTMiddleware(http.HandlerFunc(api.GetUserByIDHandler)))

	mux.HandleFunc("POST /api/v1/public/books", api.GetAllBooksHandler)
	mux.HandleFunc("GET /api/v1/public/books/{id}", api.GetBookByIDHandler)
	mux.HandleFunc("POST /api/v1/public/books/sorted", api.GetSortedBooksHandler)
	mux.HandleFunc("POST /api/v1/public/books/search", api.FilteredSearchBooksHandler)

	mux.Handle("GET /static/", staticHandle(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	mux.Handle("GET /favicon.ico", staticHandle(http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /static/books/covers/{id}", controller.ServeBookCover)

	mux.Handle("/", middleware.AuthMiddleware(http.HandlerFunc(controller.HomeHandler)))
	mux.Handle("GET /login", middleware.AuthMiddleware(http.HandlerFunc(controller.LoginHandler)))
	mux.Handle("POST /login", middleware.AuthMiddleware(http.HandlerFunc(controller.LoginHandler)))
	mux.Handle("GET /logout", middleware.AuthMiddleware(http.HandlerFunc(controller.LogoutHandler)))

	mux.Handle("GET /browse", middleware.AuthMiddleware(http.HandlerFunc(controller.BrowseHandler)))
	mux.Handle("GET /book/{id}", middleware.AuthMiddleware(http.HandlerFunc(controller.BookHandler)))

	mux.Handle("GET /games", middleware.AuthMiddleware(http.HandlerFunc(controller.GamesHandler)))
	mux.Handle("GET /games/bj", middleware.AuthMiddleware(http.HandlerFunc(bj_controller.BlackjackFrontendHandler)))

	mux.Handle("/ws/games/bj", middleware.AuthMiddleware(http.HandlerFunc(bj_controller.BlackjackGameHandler)))

	middlewareHandler := corsMiddleware(LoggingMiddleware(mux))

	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", apiPort),
		Handler: middlewareHandler,
	}
	fmt.Printf("Web server started on %s:%s\n", base_url_api, apiPort)
	if err := webServer.ListenAndServe(); err != nil {
		fmt.Printf("Web server error: %v\n", err)
	}
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

// func FrontendMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		host := r.Host
// 		parts := strings.Split(host, ".")

// 		// Check if the host has more than two parts (indicating a subdomain)
// 		if len(parts) > 2 {
// 			// Skip processing for subdomains
// 			http.Error(w, "Subdomain requests are not allowed for port 80", http.StatusForbidden)
// 			return
// 		}

// 		// Call the next handler if it's not a subdomain
// 		next.ServeHTTP(w, r)
// 	})
// }
