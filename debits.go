package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type DebitService struct {
	client *Client
}

// A Debit represents a transaction that charges/takes money from a funding
// instrument (i.e., from a credit card or bank account).
type Debit struct {
	Amount               int               `json:"amount,omitempty"` // in cents
	AppearsOnStatementAs string            `json:"appears_on_statement_as,omitempty"`
	Description          string            `json:"description,omitempty"` // dashboard description
	Meta                 map[string]string `json:"meta,omitempty"`
	Order                string            `json:"order,omitempty"`
	Currency             string            `json:"currency,omitempty"`
	FailureReason        string            `json:"failure_reason,omitempty"`
	FailureReasonCode    string            `json:"failure_reason_code,omitempty"`
	Href                 string            `json:"href,omitempty"`
	Id                   string            `json:"id,omitempty"`
	Links                *DebitLinks       `json:"links,omitempty"`
	Status               string            `json:"status,omitempty"` // "succeeded", "failed", "pending"
	TransactionNumber    string            `json:"transaction_number,omitempty"`
	CreatedAt            *time.Time        `json:"created_at,omitempty"`
	UpdatedAt            *time.Time        `json:"updated_at,omitempty"`
}

type DebitLinks struct {
	Customer string `json:"customer"`
	Dispute  string `json:"dispute"`
	Order    string `json:"order"`
	Source   string `json:"source"`
}

type debitResponse struct {
	Debits []Debit                `json:"debits"`
	Links  *debitResponseLinks    `json:"links,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

type DebitRequest debitResponse

type debitResponseLinks struct {
	Customer string `json:"debits.customer"`
	Dispute  string `json:"debits.dispute"`
	Events   string `json:"debits.events"`
	Order    string `json:"debits.order"`
	Refunds  string `json:"debits.refunds"`
	Source   string `json:"debits.source"`
}

type DebitPage struct {
	Debits []Debit
	*PaginationParams
}

func (s *DebitService) Fetch(debitId string) (*Debit, *http.Response, error) {
	path := fmt.Sprintf("/debits/%v", debitId)
	debitResponse := new(debitResponse)
	httpResponse, err := s.client.GET(path, nil, nil, debitResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &debitResponse.Debits[0], httpResponse, nil
}

func (s *DebitService) List(args ...interface{}) (*DebitPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	debitResponse := new(debitResponse)
	httpResponse, err := s.client.GET("/debits", query, nil, debitResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &DebitPage{
		Debits:           debitResponse.Debits,
		PaginationParams: NewPaginationParams(debitResponse.Meta),
	}, httpResponse, nil
}

func (s *DebitService) Update(debitId string, params map[string]interface{}) (*Debit, *http.Response, error) {
	path := fmt.Sprintf("/debits/%v", debitId)
	debitResponse := new(debitResponse)
	httpResponse, err := s.client.PUT(path, nil, params, debitResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &debitResponse.Debits[0], httpResponse, nil
}

func (s *DebitService) Refund(debitId string, refund *Refund) (*Refund, *http.Response, error) {
	path := fmt.Sprintf("/debits/%v/refunds", debitId)
	refundResponse := new(refundResponse)
	httpResponse, err := s.client.POST(path, nil, refund, refundResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &refundResponse.Refunds[0], httpResponse, err
}
