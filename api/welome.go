package api

import "net/http"

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from API server"))
}
