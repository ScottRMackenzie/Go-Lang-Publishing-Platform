package api

import (
	"encoding/json"
	"net/http"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/books"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
)

func GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := books.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	book, err := books.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func GetSortedBooksHandler(w http.ResponseWriter, r *http.Request) {
	var req BookSortingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Count == -1 {
		req.Count = 1000000
	}

	books, err := books.GetRangeWithSortingAndOrder(req.Start, req.Count, req.Sort, req.IsAcceding)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

type BookSortingRequest struct {
	Start      int    `json:"start"`
	Count      int    `json:"count"`
	Sort       string `json:"sort"`
	IsAcceding bool   `json:"is_acceding"`
}

func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {
	var req BookSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Count == -1 {
		req.Count = 1000000
	}

	if _, ok := types.ValidSortColumns[req.Sort]; !ok {
		http.Error(w, "Invalid sort column", http.StatusBadRequest)
		return
	}

	books, err := books.SearchQueryWithRange_Sorting_Order(req.Query, req.Start, req.Count, req.Sort, req.IsAcceding)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

type BookSearchRequest struct {
	Query      string `json:"query"`
	Sort       string `json:"sort"`
	Start      int    `json:"start"`
	Count      int    `json:"count"`
	IsAcceding bool   `json:"is_acceding"`
}
