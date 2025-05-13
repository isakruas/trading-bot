package model

// OrderBook holds the top-of-book bids and asks as price/quantity pairs.
type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}
