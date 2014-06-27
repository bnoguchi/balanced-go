package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type CardService struct {
	client *Client
}

type Card struct {
	ExpirationMonth int    `json:"expiration_month"` // Expiration month (e.g., 1 for January); required
	ExpirationYear  int    `json:"expiration_year"`  // Expiration year; required
	Number          string `json:"number"`           // The digits of the credit card number; required
	*Address        `json:"address,omitempty"`
	Customer        string            `json:"customer,omitempty"` // The customer this card is associated to
	Cvv             string            `json:"cvv,omitempty"`      // The 3-4 digit security code for the card
	Name            string            `json:"name,omitempty"`     // The customer's name on card
	AvsPostalMatch  interface{}       `json:"avs_postal_match,omitempty"`
	AvsResult       interface{}       `json:"avs_result,omitempty"`
	AvsStreetMatch  interface{}       `json:"avs_street_match,omitempty"`
	Brand           string            `json:"brand,omitempty"`      // e.g., "MasterCard"
	CvvMatch        string            `json:"cvv_match,omitempty"`  // e.g., "yes"
	CvvResult       string            `json:"cvv_result,omitempty"` // e.g., "Match"
	Fingerprint     string            `json:"fingerprint,omitempty"`
	Meta            map[string]string `json:"meta,omitempty"`
	Href            string            `json:"href,omitempty"` // e.g., "/cards/CC2t9628l4ecJics6T8RuLPf"
	Id              string            `json:"id,omitempty"`
	Links           *CardLinks        `json:"links,omitempty"`
	IsVerified      bool              `json:"is_verified,omitempty"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
}

type CardLinks struct {
	Customer string `json:"customer"`
}

type CardResponse struct {
	Cards []Card                 `json:"cards"`
	Links *CardResponseLinks     `json:"links"`
	Meta  map[string]interface{} `json:"meta"`
}

type CardResponseLinks struct {
	CardHolds string `json:"cards.card_holds"`
	Customers string `json:"cards.customers"`
	Debits    string `json:"cards.debits"`
}

type CardPage struct {
	Cards []Card
	*PaginationParams
}

// Creates a card on the server.
//
// Fraud and Card declinations can be reduced if the following information is
// supplied when tokenizing a card:
// Name (Name on card)
// Cvv
// PostalCode
// CountryCode (Country code, ISO 3166-1 alpha-3)
func (s *CardService) Create(card *Card) (*Card, *http.Response, error) {
	cardResponse := new(CardResponse)
	httpResponse, err := s.client.POST("/cards", nil, card, cardResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &cardResponse.Cards[0], httpResponse, nil
}

func (s *CardService) Delete(cardId string) (bool, *http.Response, error) {
	path := fmt.Sprintf("/cards/%v", cardId)
	httpResponse, err := s.client.DELETE(path, nil, nil, nil)
	if err != nil {
		return false, httpResponse, err
	}
	code := httpResponse.StatusCode
	didDelete := 200 <= code && code < 300
	return didDelete, httpResponse, nil
}

func (s *CardService) Fetch(cardId string) (*Card, *http.Response, error) {
	path := fmt.Sprintf("/cards/%v", cardId)
	cardResponse := new(CardResponse)
	httpResponse, err := s.client.GET(path, nil, nil, cardResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &cardResponse.Cards[0], httpResponse, nil
}

func (s *CardService) List(args ...interface{}) (*CardPage, *http.Response, error) {
	query := paginatedArgsToQuery(args)
	cardResponse := new(CardResponse)
	httpResponse, err := s.client.GET("/cards", query, nil, cardResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &CardPage{
		Cards:            cardResponse.Cards,
		PaginationParams: NewPaginationParams(cardResponse.Meta),
	}, httpResponse, nil
}

func (s *CardService) Update(cardId string, params map[string]interface{}) (*Card, *http.Response, error) {
	path := fmt.Sprintf("/cards/%v", cardId)
	cardResponse := new(CardResponse)
	httpResponse, err := s.client.PUT(path, nil, params, cardResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &cardResponse.Cards[0], httpResponse, nil
}

func (s *CardService) AssociateWithCustomer(cardId, customerId string) (*Card, *http.Response, error) {
	return s.Update(cardId, map[string]interface{}{
		"customer": fmt.Sprintf("/customers/%v", customerId),
	})
}

func (s *CardService) Charge(cardId string, debit *Debit) (*Debit, *http.Response, error) {
	path := fmt.Sprintf("/cards/%v/debits", cardId)
	debitResponse := new(DebitResponse)
	httpResponse, err := s.client.POST(path, nil, debit, debitResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &debitResponse.Debits[0], httpResponse, nil
}

func (s *CardService) Credit(cardId string, credit *Credit) (*Credit, *http.Response, error) {
	if credit.Amount > 250000 {
		return nil, nil, fmt.Errorf("Cannot credit more than $2,500 to a card, but tried crediting %.2f", float32(credit.Amount)/100)
	}
	path := fmt.Sprintf("/cards/%v/credits", cardId)
	creditResponse := new(CreditResponse)
	httpResponse, err := s.client.POST(path, nil, credit, creditResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &creditResponse.Credits[0], httpResponse, nil
}
