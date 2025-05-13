package service

import "trading-bot/internal/domain/model"

// Exchange defines the port that any trading exchange adapter must implement.
// This is the “driven” interface in Hexagonal/DDD architecture.
type Exchange interface {
	// Market data
	GetMarkets() ([]model.Market, error)
	GetOrderBook(market string, depth int) (*model.OrderBook, error)

	// Order management
	CreateOrder(o model.Order) (*model.Order, error)
	GetActiveOrders(market string) ([]model.Order, error)
	GetOrderByID(id string) (*model.Order, error)
	CancelOrder(id string) error
}
