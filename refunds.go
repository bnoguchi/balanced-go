package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type RefundService struct {
	client *Client
}

// A Refund is a refund of a Debit transaction. The amount of the refund may be
// any value up to the amount of the original Debit.
type Refund struct {
	Amount            int               `json:"amount,omitempty"`
	Description       string            `json:"description,omitempty"`
	Meta              map[string]string `json:"meta,omitempty"`
	Currency          string            `json:"currency,omitempty"`
	Href              string            `json:"href,omitempty"`
	Id                string            `json:"id,omitempty"`
	Links             *RefundLinks      `json:"links,omitempty"`
	Status            string            `json:"status,omitempty"`
	TransactionNumber string            `json:"transaction_number,omitempty"`
	CreatedAt         *time.Time        `json:"created_at,omitempty"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty"`
}

type RefundLinks struct {
	Debit   string `json:"debit"`
	Dispute string `json:"dispute"`
	Order   string `json:"order"`
}

type RefundResponse struct {
	Refunds []Refund               `json:"refunds"`
	Links   *RefundResponseLinks   `json:"links"`
	Meta    map[string]interface{} `json:"meta"`
}

type RefundResponseLinks struct {
	Debit   string `json:"refunds.debit"`
	Dispute string `json:"refunds.dispute"`
	Events  string `json:"refunds.events"`
	Order   string `json:"refunds.order"`
}

type RefundPage struct {
	Refunds []Refund
	*PaginationParams
}

func (s *RefundService) Fetch(refundId string) (*Refund, *http.Response, error) {
	path := fmt.Sprintf("/refunds/%v", refundId)
	refundResponse := new(RefundResponse)
	httpResponse, err := s.client.GET(path, nil, nil, refundResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &refundResponse.Refunds[0], httpResponse, nil
}

func (s *RefundService) List(args ...interface{}) (*RefundPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	refundResponse := new(RefundResponse)
	httpResponse, err := s.client.GET("/refunds", query, nil, refundResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &RefundPage{
		Refunds:          refundResponse.Refunds,
		PaginationParams: NewPaginationParams(refundResponse.Meta),
	}, httpResponse, nil
}

func (s *RefundService) Update(refundId string, params map[string]interface{}) (*Refund, *http.Response, error) {
	path := fmt.Sprintf("/refunds/%v", refundId)
	refundResponse := new(RefundResponse)
	httpResponse, err := s.client.PUT(path, nil, params, refundResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &refundResponse.Refunds[0], httpResponse, nil
}
