package bj_controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
	bj "github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/games/blackjack"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func writeMessage(conn *websocket.Conn, msg Message) error {
	fmt.Println("Writing message:", msg)
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}
	return nil
}

func setupGame(ws *websocket.Conn, deck bj.Deck, username string, betAmount int) bj.Game {
	err := UpdateBalance(ws, username, -betAmount)
	if err != nil {
		writeMessage(ws, Message{"error", "Insufficient funds", ""})
		return bj.Game{}
	}

	game := bj.Game{Deck: deck}
	game.Username = username
	game.BetAmount = betAmount

	// if the deck is empty, create a new one
	if len(deck) == 0 {
		deck = bj.NewDeck()
		game.Deck = deck
		game.ShuffleDeck()
	}

	game.PlayerInputAllowed = true

	game.PlayerHand = []bj.Card{game.DealCard(), game.DealCard()}
	game.DealerHand = []bj.Card{game.DealCard(), game.DealCard()}

	log.Println("Client Connected")
	writeMessage(ws, Message{"dealer_hand", game.DealerFirstCard(), bj.CalculateHandValue(game.DealerFirstCard()) + "+?"})
	writeMessage(ws, Message{"player_hand", game.PlayerHand, game.CalculatePlayerHandValue()})

	return game
}

func BlackjackGameHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upgrading connection")
	fmt.Println("Authenticated: ", r.Context().Value("authenticated"))
	fmt.Println("Username: ", r.Context().Value("username"))

	username := r.Context().Value("username").(string)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	game := bj.Game{}

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		fmt.Println("Received message:", string(message))

		// Handle different message types
		switch string(message) {
		case "hit":
			if game.PlayerInputAllowed {
				onPlayerHits(ws, &game)
			}
		case "stand":
			if game.PlayerInputAllowed {
				dealersTurn(ws, &game)
			}
		default:
			var receivedMsg receivedMsg
			err := json.Unmarshal(message, &receivedMsg)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
			}

			if receivedMsg.Type == "new_game" {
				if receivedMsg.Value == "" {
					writeMessage(ws, Message{"error", "Invalid bet amount", ""})
					break
				}
				betAmount, err := strconv.Atoi(receivedMsg.Value)
				if err != nil {
					log.Println("Error converting bet amount:", err)
					break
				}
				if betAmount < 1 {
					writeMessage(ws, Message{"error", "Invalid bet amount", ""})
					break
				}
				game.BetAmount = betAmount
				fmt.Println("New game with bet amount:", game.BetAmount)
				game = setupGame(ws, game.Deck, username, game.BetAmount)
				// if game.Deck == nil {
				// 	break
				// }
			}
		}
	}
}

func onPlayerHits(ws *websocket.Conn, game *bj.Game) {
	game.DealCardForPlayer()

	if err := writeMessage(ws, Message{"player_hand", game.PlayerHand, game.CalculatePlayerHandValue()}); err != nil {
		log.Println("Write error:", err)
		// break
	}

	if game.CalculatePlayerHandValue() == "Bust" || game.CalculatePlayerHandValue() == "Blackjack" {
		dealersTurn(ws, game)
	}
}

func dealersTurn(ws *websocket.Conn, game *bj.Game) {
	game.PlayerInputAllowed = false
	if err := writeMessage(ws, Message{"dealers_turn", nil, ""}); err != nil {
		log.Println("Write error:", err)
		// break
	}

	// show the full hand of the dealer
	if err := writeMessage(ws, Message{"dealer_hand", game.DealerHand, game.CalculateDealerHandValue()}); err != nil {
		log.Println("Write error:", err)
		// break
	}

	for game.DealerShouldHit() {
		time.Sleep(1000 * time.Millisecond)
		game.DealCardForDealer()
		if err := writeMessage(ws, Message{"dealer_hand", game.DealerHand, game.CalculateDealerHandValue()}); err != nil {
			log.Println("Write error:", err)
			// break
		}
	}

	resultMsg, multiplier := game.GenerateWinnerAndEarnings()
	earnings := int(float32(game.BetAmount) * multiplier)

	// Update the user's balance
	if earnings > 0 {
		UpdateBalance(ws, game.Username, earnings)
	}

	if err := writeMessage(ws, Message{"result", resultMsg, fmt.Sprintf("%d", earnings)}); err != nil {
		log.Println("Write error:", err)
		// break
	}
}

func UpdateBalance(ws *websocket.Conn, username string, amount int) error {
	user, err := users.GetByUsername(username, context.Background())
	if err != nil {
		return err
	}
	user.Balance += amount

	if user.Balance < 0 {
		return errors.New("Insufficient funds")
	}

	err = users.UpdateBalance(username, user.Balance, context.Background())
	if err != nil {
		return err
	}

	writeMessage(ws, Message{"balance", "", fmt.Sprintf("%d", user.Balance)})

	return nil
}
