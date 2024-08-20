package bj

import "sync"

type Card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
}

type Deck []Card

const (
	Hearts   = "H"
	Diamonds = "D"
	Clubs    = "C"
	Spades   = "S"
)

var suits = []string{Hearts, Diamonds, Clubs, Spades}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

var cardValues = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5,
	"6": 6, "7": 7, "8": 8, "9": 9,
	"10": 10, "J": 10, "Q": 10, "K": 10, "A": 11,
}

type Game struct {
	Deck               Deck
	PlayerHand         Deck
	DealerHand         Deck
	mu                 sync.Mutex
	PlayerInputAllowed bool
	BetAmount          int
	Username           string
}
