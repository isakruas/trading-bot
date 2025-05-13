package usecase

import "trading-bot/internal/domain/service"

// CancelOrder instructs the exchange to cancel an existing order by ID.
type CancelOrder struct {
	Ex service.Exchange
}

// Execute cancels the order and returns an error if the operation fails.
func (u *CancelOrder) Execute(id string) error {
	return u.Ex.CancelOrder(id)
}
