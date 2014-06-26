package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type DisputeService struct {
	client *Client
}

// A dispute is a customer disputed charge, aka a chargeback.
type Dispute struct {
	Amount      int               `json:"amount"`
	Currency    string            `json:"currency,omitempty"`
	Id          string            `json:"id,omitemptry"`
	Href        string            `json:"href,omitempty"`
	Status      string            `json:"status,omitempty"`
	Links       *DisputeLinks     `json:"links,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
	Reason      string            `json:"reason,omitempty"`
	InitiatedAt *time.Time        `json:"initiated_at,omitempty"`
	RespondBy   *time.Time        `json:"respond_by,omitempty"`
	CreatedAt   *time.Time        `json:"created_at,omitempty"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty"`
}

type DisputeLinks struct {
	Transaction string `json:"transaction"`
}

type DisputeResponse struct {
	Disputes []Dispute              `json:"disputes"`
	Links    *DisputeResponseLinks  `json:"links"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

type DisputeResponseLinks struct {
	Events     string `json:"disputes.events"`
	Transation string `json:"disputes.transaction"`
}

type DisputePage struct {
	Disputes []Dispute
	*PaginationParams
}

func (s *DisputeService) Fetch(disputeId string) (*Dispute, *http.Response, error) {
	path := fmt.Sprintf("/disputes/%v", disputeId)
	disputeResponse := new(DisputeResponse)
	httpResponse, err := s.client.GET(path, nil, nil, disputeResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &disputeResponse.Disputes[0], httpResponse, nil
}

func (s *DisputeService) List(args ...interface{}) (*DisputePage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	disputeResponse := new(DisputeResponse)
	httpResponse, err := s.client.GET("/disputes", query, nil, disputeResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &DisputePage{
		Disputes:         disputeResponse.Disputes,
		PaginationParams: NewPaginationParams(disputeResponse.Meta),
	}, httpResponse, nil
}
