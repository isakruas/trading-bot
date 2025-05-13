package model

// OrderSide indicates whether the order is a buy or sell.
type OrderSide string

// OrderType indicates limit, market, or other specialized order types.
type OrderType string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"

	Limit       OrderType = "LIMIT"
	MarketOrder OrderType = "MARKET"
)

// Order is the domain entity for a trading order.
// Price and Quantity are strings to preserve exchange-specific formats.
type Order struct {
	ID           string    `json:"id,omitempty"`
	MarketSymbol string    `json:"market_symbol"`
	Side         OrderSide `json:"side"`
	Type         OrderType `json:"type"`
	Price        string    `json:"price,omitempty"`
	Quantity     string    `json:"quantity,omitempty"`
	State        string    `json:"state,omitempty"`
}
