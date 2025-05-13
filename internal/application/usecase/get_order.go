package usecase

import (
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
)

// GetOrder retrieves a single order by its ID.
type GetOrder struct {
	Ex service.Exchange
}

// Execute returns the Order or an error if not found.
func (u *GetOrder) Execute(id string) (*model.Order, error) {
	return u.Ex.GetOrderByID(id)
}
