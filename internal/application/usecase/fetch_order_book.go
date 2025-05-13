package usecase

import (
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
)

// FetchOrderBook retrieves the order book (bids/asks) for a market.
type FetchOrderBook struct {
	Ex service.Exchange
}

// Execute returns the OrderBook for the given market and depth.
func (u *FetchOrderBook) Execute(market string, depth int) (*model.OrderBook, error) {
	return u.Ex.GetOrderBook(market, depth)
}
