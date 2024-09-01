package roulette

type Bets struct {
	IsSkipped  bool           `json:"isSkipped"`
	StraightUp map[string]int `json:"straightUp"`
	Color      map[string]int `json:"color"`
	EvenOdd    map[string]int `json:"evenOdd"`
	LowHigh    map[string]int `json:"lowHigh"`
	Dozens     map[string]int `json:"dozens"`
	Columns    map[string]int `json:"columns"`
}

type RoundResult struct {
	Errors         []string `json:"errors"`
	WinningNumber  string   `json:"winningNumber"`
	TotalBets      int      `json:"totalBets"`
	AmountReturned int      `json:"amountReturned"`
	NewBalance     int      `json:"newBalance"`
}
