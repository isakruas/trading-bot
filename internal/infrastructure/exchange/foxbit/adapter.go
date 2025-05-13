package foxbit

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
	"trading-bot/internal/infrastructure/httputil"
)

// FoxbitAdapter implements service.Exchange using Foxbit REST v3.
type FoxbitAdapter struct {
	apiKey     string
	secret     string
	baseURL    string
	httpClient *http.Client
}

// New returns an initialized FoxbitAdapter.
func New(apiKey, secret string) service.Exchange {
	return &FoxbitAdapter{
		apiKey:     apiKey,
		secret:     secret,
		baseURL:    "https://api.foxbit.com.br",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetMarkets implements Exchange.GetMarkets.
func (f *FoxbitAdapter) GetMarkets() ([]model.Market, error) {
	var reply struct {
		Data []model.Market `json:"data"`
	}
	err := httputil.DoRequest(f.httpClient, httputil.RequestParams{
		Method:     http.MethodGet,
		BaseURL:    f.baseURL,
		Path:       "/rest/v3/markets",
		Query:      nil,
		Body:       nil,
		APIKey:     f.apiKey,
		Secret:     f.secret,
		ResultDest: &reply,
	})
	if err != nil {
		return nil, err
	}
	return reply.Data, nil
}

// GetOrderBook implements Exchange.GetOrderBook.
func (f *FoxbitAdapter) GetOrderBook(market string, depth int) (*model.OrderBook, error) {
	params := map[string]string{"depth": strconv.Itoa(depth)}
	var ob model.OrderBook
	path := "/rest/v3/markets/" + url.PathEscape(market) + "/orderbook"
	err := httputil.DoRequest(f.httpClient, httputil.RequestParams{
		Method:     http.MethodGet,
		BaseURL:    f.baseURL,
		Path:       path,
		Query:      params,
		Body:       nil,
		APIKey:     f.apiKey,
		Secret:     f.secret,
		ResultDest: &ob,
	})
	if err != nil {
		return nil, err
	}
	return &ob, nil
}

// CreateOrder implements Exchange.CreateOrder.
func (f *FoxbitAdapter) CreateOrder(o model.Order) (*model.Order, error) {
	// parse price/quantity

	payload := map[string]interface{}{
		"side":          o.Side,
		"type":          o.Type,
		"market_symbol": o.MarketSymbol,
		"quantity":      o.Quantity,
		"price":         o.Price,
		"post_only":     true,
		"time_in_force": "GTC",
	}
	var resp struct {
		ID string `json:"id"`
	}
	err := httputil.DoRequest(f.httpClient, httputil.RequestParams{
		Method:     http.MethodPost,
		BaseURL:    f.baseURL,
		Path:       "/rest/v3/orders",
		Query:      nil,
		Body:       payload,
		APIKey:     f.apiKey,
		Secret:     f.secret,
		ResultDest: &resp,
	})
	if err != nil {
		return nil, err
	}
	o.ID = resp.ID
	return &o, nil
}

// GetActiveOrders implements Exchange.GetActiveOrders.
func (f *FoxbitAdapter) GetActiveOrders(market string) ([]model.Order, error) {
	params := map[string]string{"state": "ACTIVE"}
	if market != "" {
		params["market_symbol"] = market
	}
	var reply struct {
		Data []model.Order `json:"data"`
	}
	err := httputil.DoRequest(f.httpClient, httputil.RequestParams{
		Method:     http.MethodGet,
		BaseURL:    f.baseURL,
		Path:       "/rest/v3/orders",
		Query:      params,
		Body:       nil,
		APIKey:     f.apiKey,
		Secret:     f.secret,
		ResultDest: &reply,
	})
	if err != nil {
		return nil, err
	}
	return reply.Data, nil
}

// GetOrderByID implements Exchange.GetOrderByID.
func (f *FoxbitAdapter) GetOrderByID(id string) (*model.Order, error) {
	var o model.Order
	path := "/rest/v3/orders/by-order-id/" + url.PathEscape(id)
	err := httputil.DoRequest(f.httpClient, httputil.RequestParams{
		Method:     http.MethodGet,
		BaseURL:    f.baseURL,
		Path:       path,
		Query:      nil,
		Body:       nil,
		APIKey:     f.apiKey,
		Secret:     f.secret,
		ResultDest: &o,
	})
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// CancelOrder implements Exchange.CancelOrder.
func (f *FoxbitAdapter) CancelOrder(id string) error {
	payload := map[string]interface{}{
		"type": "ID",
		"id":   id,
	}
	return httputil.DoRequest(f.httpClient, httputil.RequestParams{
		Method:     http.MethodPut,
		BaseURL:    f.baseURL,
		Path:       "/rest/v3/orders/cancel",
		Query:      nil,
		Body:       payload,
		APIKey:     f.apiKey,
		Secret:     f.secret,
		ResultDest: nil,
	})
}
