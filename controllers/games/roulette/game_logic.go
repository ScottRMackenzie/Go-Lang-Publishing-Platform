package roulette

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
)

func RouletteGameHandler(w http.ResponseWriter, r *http.Request) {
	errors := []string{}

	winningNumber := generateWinningNumber()

	var bets Bets
	if err := json.NewDecoder(r.Body).Decode(&bets); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	user, err := users.GetByUsername(username, context.Background())
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	charge, err := chargeForBets(bets)
	if err != nil {
		http.Error(w, "Error calculating charge", http.StatusInternalServerError)
		return
	}

	if user.Balance < charge {
		errors = append(errors, "insufficient funds")
		return
	}
	user.Balance -= charge

	amountReturned, err := calculateAmountReturned(bets, winningNumber)
	if err != nil {
		errors = append(errors, "error calculating amount returned")
		amountReturned = charge
	}

	user.Balance += amountReturned

	if err := users.UpdateBalance(username, user.Balance, context.Background()); err != nil {
		http.Error(w, "Error updating user balance", http.StatusInternalServerError)
		return
	}

	result := RoundResult{
		Errors:         []string{},
		WinningNumber:  winningNumber,
		TotalBets:      charge,
		AmountReturned: amountReturned,
		NewBalance:     user.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func generateWinningNumber() string {
	return strconv.Itoa(rand.Intn(37))
}

func chargeForBets(bets Bets) (int, error) {
	totalCharge := 0
	for _, bet := range bets.StraightUp {
		totalCharge += bet
	}

	for _, bet := range bets.Color {
		totalCharge += bet
	}

	for _, bet := range bets.EvenOdd {
		totalCharge += bet
	}

	for _, bet := range bets.LowHigh {
		totalCharge += bet
	}

	for _, bet := range bets.Dozens {
		totalCharge += bet
	}

	for _, bet := range bets.Columns {
		totalCharge += bet
	}

	return totalCharge, nil
}

func calculateAmountReturned(bets Bets, winningNumber string) (int, error) {
	totalReturns := 0
	if _, ok := bets.StraightUp[winningNumber]; ok {
		totalReturns += bets.StraightUp[winningNumber]*35 + bets.StraightUp[winningNumber]
	}

	betColor, err := getColor(winningNumber)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if _, ok := bets.Color[betColor]; ok {
		totalReturns += bets.Color[betColor]*2 + bets.Color[betColor]
	}

	if _, ok := bets.EvenOdd["even"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num%2 == 0 {
			totalReturns += bets.EvenOdd["even"]*2 + bets.EvenOdd["even"]
		}
	}

	if _, ok := bets.EvenOdd["odd"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num%2 != 0 {
			totalReturns += bets.EvenOdd["odd"]*2 + bets.EvenOdd["odd"]
		}
	}

	if _, ok := bets.LowHigh["low"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num >= 1 && num <= 18 {
			totalReturns += bets.LowHigh["low"]*2 + bets.LowHigh["low"]
		}
	}

	if _, ok := bets.LowHigh["high"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num >= 19 && num <= 36 {
			totalReturns += bets.LowHigh["high"]*2 + bets.LowHigh["high"]
		}
	}

	if _, ok := bets.Dozens["1st"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num >= 1 && num <= 12 {
			totalReturns += bets.Dozens["1st"]*3 + bets.Dozens["1st"]
		}
	}

	if _, ok := bets.Dozens["2nd"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num >= 13 && num <= 24 {
			totalReturns += bets.Dozens["2nd"]*3 + bets.Dozens["2nd"]
		}
	}

	if _, ok := bets.Dozens["3rd"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num >= 25 && num <= 36 {
			totalReturns += bets.Dozens["3rd"]*3 + bets.Dozens["3rd"]
		}
	}

	if _, ok := bets.Columns["1st"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num%3 == 1 {
			totalReturns += bets.Columns["1st"]*3 + bets.Columns["1st"]
		}
	}

	if _, ok := bets.Columns["2nd"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num%3 == 2 {
			totalReturns += bets.Columns["2nd"]*3 + bets.Columns["2nd"]
		}
	}

	if _, ok := bets.Columns["3rd"]; ok {
		if num, err := strconv.Atoi(winningNumber); err == nil && num%3 == 0 {
			totalReturns += bets.Columns["3rd"]*3 + bets.Columns["3rd"]
		}
	}

	return totalReturns, nil
}

func getColor(number string) (string, error) {
	num, err := strconv.Atoi(number)
	if err != nil {
		return "0", err
	}

	if num < 0 || num > 36 {
		return "", errors.New("invalid number")
	}

	if number == "0" {
		return "green", nil
	}

	blackNumbers := map[string]bool{
		"2":  true,
		"4":  true,
		"6":  true,
		"8":  true,
		"10": true,
		"11": true,
		"13": true,
		"15": true,
		"17": true,
		"20": true,
		"22": true,
		"24": true,
		"26": true,
		"28": true,
		"29": true,
		"31": true,
		"33": true,
		"35": true,
	}

	if _, ok := blackNumbers[number]; ok {
		return "black", nil
	}

	return "red", nil
}
