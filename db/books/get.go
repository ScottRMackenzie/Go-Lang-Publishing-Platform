package books

import (
	"context"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
)

func GetAll() ([]types.Book, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, author, published_date, isbn, genre, language_code, publisher, created_at, updated_at, summary, word_count FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []types.Book
	for rows.Next() {
		var book types.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Isbn, &book.Genre, &book.LanguageCode, &book.Publisher, &book.CreatedAt, &book.UpdatedAt, &book.Summary, &book.WordCount); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func GetByID(id string) (types.Book, error) {
	var book types.Book
	err := db.Pool.QueryRow(context.Background(), "SELECT id, title, author, published_date, isbn, genre, language_code, publisher, created_at, updated_at, summary, word_count FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Isbn, &book.Genre, &book.LanguageCode, &book.Publisher, &book.CreatedAt, &book.UpdatedAt, &book.Summary, &book.WordCount)
	if err != nil {
		return types.Book{}, err
	}

	return book, nil
}

func GetRangeWithSortingAndOrder(start, count int, sort string, isAcceding bool) ([]types.Book, error) {
	order := "DESC"
	if isAcceding {
		order = "ASC"
	}

	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, author, published_date, isbn, genre, language_code, publisher, created_at, updated_at, summary, word_count FROM books ORDER BY "+sort+" "+order+" LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []types.Book
	for rows.Next() {
		var book types.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Isbn, &book.Genre, &book.LanguageCode, &book.Publisher, &book.CreatedAt, &book.UpdatedAt, &book.Summary, &book.WordCount); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func SearchQueryWithRange_Sorting_Order(searchQuery string, start, count int, sort string, isAcceding bool) ([]types.Book, error) {
	order := "DESC"
	if isAcceding {
		order = "ASC"
	}

	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, author, published_date, isbn, genre, language_code, publisher, created_at, updated_at, summary, word_count FROM books WHERE title ILIKE $1 OR author ILIKE $1 OR isbn ILIKE $1 OR genre ILIKE $1 OR publisher ILIKE $1 OR summary ILIKE $1 ORDER BY "+sort+" "+order+" LIMIT $2 OFFSET $3", "%"+searchQuery+"%", count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []types.Book
	for rows.Next() {
		var book types.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Isbn, &book.Genre, &book.LanguageCode, &book.Publisher, &book.CreatedAt, &book.UpdatedAt, &book.Summary, &book.WordCount); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
