package usecase

import (
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
)

// FetchMarkets retrieves all trading markets from the exchange.
type FetchMarkets struct {
	Ex service.Exchange
}

// Execute returns a slice of Market models or an error.
func (u *FetchMarkets) Execute() ([]model.Market, error) {
	return u.Ex.GetMarkets()
}
