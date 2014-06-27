package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type CreditService struct {
	client *Client
}

type Credit struct {
	Amount               int                    `json:"amount"`
	AppearsOnStatementAs string                 `json:"appears_on_statement_as,omitempty"`
	Destination          string                 `json:"destination,omitempty"` // The funding destination for this credit
	Order                string                 `json:"order,omitempty"`
	Currency             string                 `json:"currency,omitempty"`
	Description          string                 `json:"description,omitempty"`
	FailureReason        string                 `json:"failure_reason,omitempty"`
	FailureReasonCode    string                 `json:"failure_reason_code,omitempty"`
	Href                 string                 `json:"href,omitempty"`
	Id                   string                 `json:"id,omitempty"`
	Links                *CreditLinks           `json:"links,omitempty"`
	Meta                 map[string]interface{} `json:"meta,omitempty"`
	Status               string                 `json:"status,omitempty"`
	TransactionNumber    string                 `json:"transaction_number,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
}

type CreditResponse struct {
	Credits []Credit               `json:"credits"`
	Links   *CreditResponseLinks   `json:"links"`
	Meta    map[string]interface{} `json:"meta"`
}

type CreditLinks struct {
	Customer    string `json:"customer"`
	Destination string `json:"destination"`
	Order       string `json:"order"`
}

type CreditResponseLinks struct {
	Customer    string `json:"credits.customer"`
	Destination string `json:"credits.destination"`
	Events      string `json:"credits.events"`
	Order       string `json:"credits.order"`
	Reversals   string `json:"credits.reversals"`
}

// CreditPage holds a paginated set of credits
type CreditPage struct {
	Credits []Credit
	*PaginationParams
}

func (s *CreditService) CreateToBankAccount(accountId string, credit *Credit) (*Credit, *http.Response, error) {
	return s.client.BankAccount.Credit(accountId, credit)
}

func (s *CreditService) CreateToCard(cardId string, credit *Credit) (*Credit, *http.Response, error) {
	return s.client.Card.Credit(cardId, credit)
}

// CreateForOrder credits money from the order to the seller's BankAccount
// represented by bankAccountId.
func (s *CreditService) CreateForOrder(bankAccountId, orderId string, credit *Credit) (*Credit, *http.Response, error) {
	credit.Order = fmt.Sprintf("/orders/%v", orderId)
	return s.CreateToBankAccount(bankAccountId, credit)
}

func (s *CreditService) Fetch(creditId string) (*Credit, *http.Response, error) {
	path := fmt.Sprintf("/credits/%v", creditId)
	creditResponse := new(CreditResponse)
	httpResponse, err := s.client.GET(path, nil, nil, creditResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &creditResponse.Credits[0], httpResponse, nil
}

func (s *CreditService) List(args ...interface{}) (*CreditPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	creditResponse := new(CreditResponse)
	httpResponse, err := s.client.GET("/credits", query, nil, creditResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &CreditPage{
		Credits:          creditResponse.Credits,
		PaginationParams: NewPaginationParams(creditResponse.Meta),
	}, httpResponse, nil
}

func (s *CreditService) ListForBankAccount(accountId string, args ...interface{}) (*CreditPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	path := fmt.Sprintf("/bank_accounts/%v/credits", accountId)
	creditResponse := new(CreditResponse)
	httpResponse, err := s.client.GET(path, query, nil, creditResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &CreditPage{
		Credits:          creditResponse.Credits,
		PaginationParams: NewPaginationParams(creditResponse.Meta),
	}, httpResponse, nil
}

func (s *CreditService) Update(creditId string, params map[string]interface{}) (*Credit, *http.Response, error) {
	path := fmt.Sprintf("/credits/%v", creditId)
	creditResponse := new(CreditResponse)
	httpResponse, err := s.client.PUT(path, nil, params, creditResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &creditResponse.Credits[0], httpResponse, nil
}
