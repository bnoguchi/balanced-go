package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type EventService struct {
	client *Client
}

type Event struct {
	Id                string            `json:"id,omitempty"`
	Href              string            `json:"href,omitempty"`
	Type              string            `json:"type"`
	OccurredAt        *time.Time        `json:"occurred_at"`
	Links             map[string]string `json:"links"`
	Entity            *EventEntity      `json:"entity"`
	*CallbackStatuses `json:"callback_statuses"`
}

type EventEntity struct {
	Customers     []Customer     `json:"customers,omitempty"`
	BankAccounts  []BankAccount  `json:"bank_accounts,omitempty"`
	Cards         []Card         `json:"cards,omitempty"`
	CardHolds     []CardHold     `json:"card_holds,omitempty"`
	Debits        []Debit        `json:"debits,omitempty"`
	Credits       []Credit       `json:"credits,omitempty"`       // Not found in docs
	Disputes      []Dispute      `json:"disputes,omitempty"`      // Not found in docs
	Orders        []Order        `json:"orders,omitempty"`        // Not found in docs
	Refunds       []Refund       `json:"refunds,omitempty"`       // Not found in docs
	Reversals     []Reversal     `json:"reversals,omitempty"`     // Not found in docs
	Verifications []Verification `json:"verifications,omitempty"` // Not found in docs
}

type CallbackStatuses struct {
	Failed    int `json:"failed"`
	Pending   int `json:"pending"`
	Retrying  int `json:"retrying"`
	Succeeded int `json:"succeeded"`
}

type EventResponse struct {
	Events []Event                `json:"events"`
	Links  *EventResponseLinks    `json:"links"`
	Meta   map[string]interface{} `json:"meta"`
}

type EventResponseLinks struct {
	Callbacks string `json:"events.callbacks"` // "/events/{events.self}/callbacks"
}

type EventPage struct {
	Events []Event
	*PaginationParams
}

func (s *EventService) Fetch(eventId string) (*Event, *http.Response, error) {
	path := fmt.Sprintf("/events/%v", eventId)
	eventResponse := new(EventResponse)
	httpResponse, err := s.client.GET(path, nil, nil, eventResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &eventResponse.Events[0], httpResponse, nil
}

func (s *EventService) List(args ...interface{}) (*EventPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	eventResponse := new(EventResponse)
	httpResponse, err := s.client.GET("/events", query, nil, eventResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &EventPage{
		Events:           eventResponse.Events,
		PaginationParams: NewPaginationParams(eventResponse.Meta),
	}, httpResponse, nil
}
