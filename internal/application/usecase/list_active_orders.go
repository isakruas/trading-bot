package usecase

import (
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
)

// ListActiveOrders returns all active orders for a market.
type ListActiveOrders struct {
	Ex service.Exchange
}

// Execute returns a slice of active Orders or an error.
func (u *ListActiveOrders) Execute(market string) ([]model.Order, error) {
	return u.Ex.GetActiveOrders(market)
}
