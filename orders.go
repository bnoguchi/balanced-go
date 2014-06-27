package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type OrderService struct {
	client *Client
}

// An Order is a construct that logically groups related transaction operations
// for a particular seller (Customer). An Order allows issuing payouts to only
// one Customer and the marketplace bank account. An Order is useful for
// reconciliation purposes, as each Order maintains its own individual escrow
// balance, which is separate from the total marketplace escrow. Attempts to
// credit an Order beyond the amount debited into the Order will fail.
type Order struct {
	Id              string                 `json:"id,omitempty"`
	Href            string                 `json:"href,omitempty"`
	Amount          int                    `json:"amount,omitempty"`
	AmountEscrowed  int                    `json:"amount_escrowed,omitempty"`
	Currency        string                 `json:"currency,omitempty"`
	DeliveryAddress *Address               `json:"delivery_address,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Links           *OrderLinks            `json:"links,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	CreatedAt       *time.Time             `json:"created_at,omitempty"`
	UpdatedAt       *time.Time             `json:"updated_at,omitempty"`
}

type OrderLinks struct {
	Merchant string `json:"merchant"`
}

type OrderResponse struct {
	Orders []Order                `json:"orders"`
	Meta   map[string]interface{} `json:"meta"`
	Links  *OrderResponseLinks    `json:"links"`
}

type OrderResponseLinks struct {
	Buyers    string `json:"orders.buyers"`
	Credits   string `json:"orders.credits"`
	Debits    string `json:"orders.debits"`
	Merchant  string `json:"orders.merchant"`
	Refunds   string `json:"orders.refunds"`
	Reversals string `json:"orders.reversals"`
}

type OrderPage struct {
	Orders []Order
	*PaginationParams
}

func (s *OrderService) Create(customerId string, order *Order) (*Order, *http.Response, error) {
	path := fmt.Sprintf("/customers/%v/orders", customerId)
	orderResponse := new(OrderResponse)
	httpResponse, err := s.client.POST(path, nil, order, orderResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &orderResponse.Orders[0], httpResponse, nil
}

func (s *OrderService) Fetch(orderId string) (*Order, *http.Response, error) {
	path := fmt.Sprintf("/orders/%v", orderId)
	orderResponse := new(OrderResponse)
	httpResponse, err := s.client.GET(path, nil, nil, orderResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &orderResponse.Orders[0], httpResponse, nil
}

func (s *OrderService) List(args ...interface{}) (*OrderPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	orderResponse := new(OrderResponse)
	httpResponse, err := s.client.GET("/orders", query, nil, orderResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &OrderPage{
		Orders:           orderResponse.Orders,
		PaginationParams: NewPaginationParams(orderResponse.Meta),
	}, httpResponse, nil
}

func (s *OrderService) Update(orderId string, params map[string]interface{}) (*Order, *http.Response, error) {
	path := fmt.Sprintf("/orders/%v", orderId)
	orderResponse := new(OrderResponse)
	httpResponse, err := s.client.PUT(path, nil, params, orderResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &orderResponse.Orders[0], httpResponse, nil
}
