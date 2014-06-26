package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type ReversalService struct {
	client *Client
}

type Reversal struct {
	Amount      int               `json:"amount,omitempty"`
	Description string            `json:"description,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
	Links       *ReversalLinks    `json:"links,omitempty"`
	Id          string            `json:"id,omitempty"`
	Href        string            `json:"href,omitempty"`
	CreatedAt   *time.Time        `json:"created_at,omitempty"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty"`
}

type ReversalLinks struct {
	Credit string `json:"credit"`
	Order  string `json:"order,omitempty"`
}

type ReversalResponse struct {
	Reversals []Reversal             `json:"reversals"`
	Links     *ReversalResponseLinks `json:"links"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

type ReversalResponseLinks struct {
	Credit string `json:"reversals.credit"`
	Events string `json:"reversals.events"`
	Order  string `json:"reversal.order"`
}

type ReversalPage struct {
	Reversals []Reversal
	*PaginationParams
}

func (s *ReversalService) Create(creditId string, reversal *Reversal) (*Reversal, *http.Response, error) {
	path := fmt.Sprintf("/credits/%v/reversals", creditId)
	reversalResponse := new(ReversalResponse)
	httpResponse, err := s.client.POST(path, nil, reversal, reversalResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &reversalResponse.Reversals[0], httpResponse, nil
}

func (s *ReversalService) Fetch(reversalId string) (*Reversal, *http.Response, error) {
	path := fmt.Sprintf("/reversals/%v", reversalId)
	reversalResponse := new(ReversalResponse)
	httpResponse, err := s.client.GET(path, nil, nil, reversalResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &reversalResponse.Reversals[0], httpResponse, nil
}

func (s *ReversalService) List(args ...interface{}) (*ReversalPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	reversalResponse := new(ReversalResponse)
	httpResponse, err := s.client.GET("/reversals", query, nil, reversalResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &ReversalPage{
		Reversals:        reversalResponse.Reversals,
		PaginationParams: NewPaginationParams(reversalResponse.Meta),
	}, httpResponse, nil
}

func (s *ReversalService) Update(reversalId string, params map[string]interface{}) (*Reversal, *http.Response, error) {
	path := fmt.Sprintf("/reversals/%v", reversalId)
	reversalResponse := new(ReversalResponse)
	httpResponse, err := s.client.PUT(path, nil, params, reversalResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &reversalResponse.Reversals[0], httpResponse, nil
}
