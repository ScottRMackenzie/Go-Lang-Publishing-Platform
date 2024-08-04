package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/auth"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.VerifyJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to context if needed
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			// User is not authenticated
			ctx := context.WithValue(r.Context(), "authenticated", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, err := auth.VerifyJWT(cookie.Value)
		if err != nil {
			// Invalid token
			ctx := context.WithValue(r.Context(), "authenticated", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Valid token
		ctx := context.WithValue(r.Context(), "authenticated", true)
		ctx = context.WithValue(ctx, "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
