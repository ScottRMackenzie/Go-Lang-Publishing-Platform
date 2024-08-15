package controller

import (
	"net/http"
	"os"
	"path/filepath"
)

func ServeBookCover(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("id") + ".jpg"

	if path == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join("static", "books", "covers", path)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		filepath = "static/books/covers/default.jpg"
	}

	http.ServeFile(w, r, filepath)

}
