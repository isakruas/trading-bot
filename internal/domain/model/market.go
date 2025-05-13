package model

// Market represents a trading pair and its precision/increment rules.
type Market struct {
	Symbol            string `json:"symbol"`
	PriceMin          string `json:"price_min"`
	PriceIncrement    string `json:"price_increment"`
	PricePrecision    int    `json:"price_precision"`
	QuantityMin       string `json:"quantity_min"`
	QuantityIncrement string `json:"quantity_increment"`
	QuantityPrecision int    `json:"quantity_precision"`
}
