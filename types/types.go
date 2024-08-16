package types

import (
	"time"
)

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	IsVerified bool      `json:"is_verified"`
}

type Book struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	PublishedDate   time.Time `json:"published_date"`
	Isbn            string    `json:"isbn"`
	Genre           string    `json:"genre"`
	LanguageCode    string    `json:"language_code"`
	Publisher       string    `json:"publisher"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Summary         string    `json:"summary"`
	WordCount       int       `json:"word_count"`
	CoverImgUrl     string    `json:"cover_img_url"`
	Price           Price     `json:"price"`
	DiscountedPrice float64   `json:"discounted_price"`
}

type Discount struct {
	ID         string    `json:"id"`
	BookID     string    `json:"book_id"`
	Percentage int       `json:"percentage"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Price struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Discount  Discount  `json:"discount"`
}

var ValidSortColumns = map[string]bool{
	"id": true, "title": true, "author": true, "published_date": true,
	"isbn": true, "genre": true, "language_code": true, "publisher": true,
	"created_at": true, "updated_at": true, "summary": true, "word_count": true,
}

var ValidOrder = map[string]bool{
	"ASC": true, "DESC": true,
}

type Filters struct {
	CaseSensitive map[string]bool   `json:"case_sensitive"`
	Values        map[string]string `json:"values"`
}

func InitializeNewFilters() Filters {
	return Filters{
		CaseSensitive: make(map[string]bool),
		Values:        make(map[string]string),
	}
}
