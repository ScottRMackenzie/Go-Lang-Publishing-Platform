package bj

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func NewDeck() Deck {
	var deck Deck
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

func (g *Game) ShuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
}

func (g *Game) DealCard() Card {
	if g.Deck == nil || len(g.Deck) == 0 {
		g.Deck = NewDeck()
		g.ShuffleDeck()
	}

	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	return card
}

func (g *Game) DealCardForPlayer() {
	g.PlayerHand = append(g.PlayerHand, g.DealCard())
}

func (g *Game) DealCardForDealer() {
	g.DealerHand = append(g.DealerHand, g.DealCard())
}

func CalculateHandValue(hand []Card) string {
	total := 0
	aceCount := 0

	for _, card := range hand {
		card.Value = strings.ToUpper(card.Value)
		if value, exists := cardValues[card.Value]; exists {
			total += value
			if card.Value == "A" {
				aceCount++
			}
		} else {
			fmt.Printf("Invalid card: %s\n", card)
		}
	}

	for aceCount > 0 && total > 21 {
		total -= 10
		aceCount--
	}

	if total > 21 {
		return "Bust"
	}

	if total == 21 {
		return "Blackjack"
	}

	stringScore := fmt.Sprintf("%d", total)
	for aceCount > 0 {
		stringScore += fmt.Sprintf(" / %d", total-10)
		aceCount--
	}
	return stringScore
}

// func GetCardValue(card Card) int {
// 	card.Value = strings.ToUpper(card.Value)
// 	if value, exists := cardValues[card.Value]; exists {
// 		return value
// 	} else {
// 		fmt.Printf("Invalid card: %s\n", card)
// 	}

// 	return 0
// }

func (g *Game) CalculatePlayerHandValue() string {
	return CalculateHandValue(g.PlayerHand)
}

func (g *Game) CalculateDealerHandValue() string {
	return CalculateHandValue(g.DealerHand)
}

func (g *Game) DealerFirstCard() Deck {
	return g.DealerHand[:1]
}

func softHandValue(handValue string) string {
	if !strings.Contains(handValue, "/") {
		return handValue
	}

	values := strings.Split(handValue, " / ")
	for _, value := range values {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error converting string to int")
			return ""
		}

		if valueInt < 21 {
			return value
		}
	}

	return values[len(values)-1]
}

func (g *Game) DealerShouldHit() bool {
	// if g.CalculateDealerHandValue() == "Blackjack" || g.CalculateDealerHandValue() == "Bust" {
	// 	return false
	// }

	dealerValueStr := g.CalculateDealerHandValue()

	// if dealerValue contains a "/", then it's a soft hand
	if strings.Contains(dealerValueStr, "/") {
		dealerValueStr = softHandValue(dealerValueStr)
	}

	dealerValue, err := strconv.Atoi(dealerValueStr)
	if err != nil {
		return false
	}

	return dealerValue < 17
}

// func (g *Game) CalculateWinner() string {
// 	if g.CalculateDealerHandValue() > 21 || g.CalculateDealerHandValue() < g.CalculatePlayerHandValue() {
// 		return "player"
// 	} else if g.CalculateDealerHandValue() == g.CalculatePlayerHandValue() {
// 		return "push"
// 	} else {
// 		return "dealer"
// 	}
// }

func (g *Game) GenerateWinnerAndEarnings() (string, float32) {
	playerHandValue := softHandValue(g.CalculatePlayerHandValue())
	dealerHandValue := softHandValue(g.CalculateDealerHandValue())

	if playerHandValue == "Blackjack" {
		return "Player wins with Blackjack!", 2.5
	} else if playerHandValue == "Bust" {
		return "Player busts! Dealer wins!", 0
	} else if dealerHandValue == "Blackjack" {
		return "Dealer wins with Blackjack!", 0
	} else if dealerHandValue == "Bust" {
		return "Dealer busts! Player wins!", 2
	} else if playerHandValue > dealerHandValue {
		return "Player wins!", 2
	} else if playerHandValue < dealerHandValue {
		return "Dealer wins!", 0
	} else {
		return "Push!", 1
	}
}
