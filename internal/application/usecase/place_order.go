package usecase

import (
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
)

// PlaceOrder is the application service to create a new order.
// It orchestrates domain calls and returns the created Order.
type PlaceOrder struct {
	Ex service.Exchange
}

// Execute sends the order request to the exchange and returns the filled Order.
func (u *PlaceOrder) Execute(req model.Order) (*model.Order, error) {
	return u.Ex.CreateOrder(req)
}
