package balanced

import (
	"fmt"
	"net/http"
	"time"
)

// CardHoldService provides methods to interact with Balanced API for card
// holds. Balanced does not recommend holds and considers their usage to be an
// advanced feature because of the confusion they cause and the difficulty in
// managing them.
type CardHoldService struct {
	client *Client
}

type CardHold struct {
	Amount            int                    `json:"amount"`
	Currency          string                 `json:"currency,omitempty"`
	Description       string                 `json:"description,omitempty"`
	ExpiresAt         *time.Time             `json:"expires_at,omitempty"`
	FailureReason     string                 `json:"failure_reason,omitempty"`
	FailureReasonCode string                 `json:"failure_reason_code,omitempty"`
	Href              string                 `json:"href,omitempty"`
	Id                string                 `json:"id,omitempty"`
	Links             *CardHoldLinks         `json:"links,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Status            string                 `json:"status,omitempty"`
	TransactionNumber string                 `json:"transaction_number,omitempty"`
	VoidedAt          *time.Time             `json:"voided_at,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
}

type CardHoldLinks struct {
	Card  string `json:"card"`
	Debit string `json:"debit"`
}

// CardHoldPage holds a paginated set of card holds
type CardHoldPage struct {
	CardHolds []CardHold
	*PaginationParams
}

type cardHoldResponse struct {
	CardHolds []CardHold             `json:"card_holds"`
	Links     *cardHoldResponseLinks `json:"links"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

type cardHoldResponseLinks struct {
	Card   string `json:"card_holds.card"`
	Debit  string `json:"card_holds.debit"`
	Debits string `json:"card_holds.debits"`
	Events string `json:"card_holds.events"`
}

func (s *CardHoldService) Create(cardId string, hold *CardHold) (*CardHold, *http.Response, error) {
	path := fmt.Sprintf("/cards/%v/card_holds", cardId)
	holdResponse := new(cardHoldResponse)
	httpResponse, err := s.client.POST(path, nil, hold, holdResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &holdResponse.CardHolds[0], httpResponse, nil
}

func (s *CardHoldService) Fetch(holdId string) (*CardHold, *http.Response, error) {
	path := fmt.Sprintf("/card_holds/%v", holdId)
	holdResponse := new(cardHoldResponse)
	httpResponse, err := s.client.GET(path, nil, nil, holdResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &holdResponse.CardHolds[0], httpResponse, nil
}

// The holds are returned sorted from recent to oldest
func (s *CardHoldService) List(args ...interface{}) (*CardHoldPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	holdResponse := new(cardHoldResponse)
	httpResponse, err := s.client.GET("/card_holds", query, nil, holdResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &CardHoldPage{
		CardHolds:        holdResponse.CardHolds,
		PaginationParams: NewPaginationParams(holdResponse.Meta),
	}, httpResponse, nil
}

func (s *CardHoldService) Update(holdId string, params map[string]interface{}) (*CardHold, *http.Response, error) {
	path := fmt.Sprintf("/card_holds/%v", holdId)
	holdResponse := new(cardHoldResponse)
	httpResponse, err := s.client.PUT(path, nil, params, holdResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &holdResponse.CardHolds[0], httpResponse, nil
}

// Captures a previously created card hold. This creates a Debit. Any amount up
// to the hold amount may be captured.
func (s *CardHoldService) Capture(holdId string, debit *Debit) (*Debit, *http.Response, error) {
	path := fmt.Sprintf("/card_holds/%v/debits", holdId)
	debitResponse := new(debitResponse)
	httpResponse, err := s.client.POST(path, nil, debit, debitResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &debitResponse.Debits[0], httpResponse, nil
}

// Cancels the hold. Once voided, the hold can no longer be captured.
func (s *CardHoldService) Void(holdId string) (*CardHold, *http.Response, error) {
	path := fmt.Sprintf("/card_holds/%v", holdId)
	reqBody := map[string]interface{}{
		"is_void": true,
	}
	holdResponse := new(cardHoldResponse)
	httpResponse, err := s.client.PUT(path, nil, reqBody, holdResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &holdResponse.CardHolds[0], httpResponse, nil
}
