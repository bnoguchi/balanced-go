package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type CustomerService struct {
	client *Client
}

type Customer struct {
	*Address     `json:"address,omitempty"`
	BusinessName string `json:"business_name,omitempty"`
	Customers    string `json:"customers,omitempty"`
	DobMonth     int    `json:"dob_month,omitempty"`
	DobYear      int    `json:"dob_year,omitempty"`
	Ein          string `json:"ein,omitempty"`
	Name         string `json:"ein,omitempty"`
	Phone        string `json:"phone,omitempty"`

	Email          string                 `json:"email,omitempty"`
	Href           string                 `json:"href,omitempty"`
	Id             string                 `json:"id,omitempty"`
	Links          *CustomerLinks         `json:"links,omitempty"`
	MerchantStatus string                 `json:"merchant_status,omitempty"`
	Meta           map[string]interface{} `json:"meta,omitempty"`
	SsnLast4       string                 `json:"ssn_last4,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`
}

type CustomerLinks struct {
	Destination string `json:"destination"`
	Source      string `json:"source"`
}

// CustomerPage holds a paginated set of customers
type CustomerPage struct {
	Customers []Customer
	*PaginationParams
}

type CustomerResponse struct {
	Customers []Customer             `json:"customers"`
	Links     CustomerResponseLinks  `json:"links"`
	Meta      map[string]interface{} `json:"meta"`
}

type CustomerResponseLinks struct {
	BankAccounts     string `json:"customers.bank_accounts"`
	CardHolds        string `json:"customers.card_holds"`
	Cards            string `json:"customers.cards"`
	Credits          string `json:"customers.credits"`
	Debits           string `json:"customers.debits"`
	Destination      string `json:"customers.destination"`
	ExternalAccounts string `json:"customers.external_accounts"`
	Orders           string `json:"customers.orders"`
	Refunds          string `json:"customers.refunds"`
	Reversals        string `json:"customers.reversals"`
	Source           string `json:"customers.source"`
	Transactions     string `json:"customers.transactions"`
}

func (s *CustomerService) Create(customer *Customer) (*Customer, *http.Response, error) {
	customerResponse := new(CustomerResponse)
	httpResponse, err := s.client.POST("/customers", nil, customer, customerResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &customerResponse.Customers[0], httpResponse, nil
}

func (s *CustomerService) Delete(customerId string) (bool, *http.Response, error) {
	path := fmt.Sprintf("/customers/%v", customerId)
	httpResponse, err := s.client.DELETE(path, nil, nil, nil)
	if err != nil {
		return false, httpResponse, err
	}
	code := httpResponse.StatusCode
	didDelete := 200 <= code && code < 300
	return didDelete, httpResponse, nil
}

func (s *CustomerService) List(args ...interface{}) (*CustomerPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	customerResponse := new(CustomerResponse)
	httpResponse, err := s.client.GET("/customers", query, nil, customerResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &CustomerPage{
		Customers:        customerResponse.Customers,
		PaginationParams: NewPaginationParams(customerResponse.Meta),
	}, httpResponse, nil
}

func (s *CustomerService) Fetch(customerId string) (*Customer, *http.Response, error) {
	path := fmt.Sprintf("/customers/%v", customerId)
	customerResponse := new(CustomerResponse)
	httpResponse, err := s.client.GET(path, nil, nil, customerResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &customerResponse.Customers[0], httpResponse, nil
}

func (s *CustomerService) Update(customerId string, params map[string]interface{}) (*Customer, *http.Response, error) {
	path := fmt.Sprintf("/customers/%v", customerId)
	customerResponse := new(CustomerResponse)
	httpResponse, err := s.client.PUT(path, nil, params, customerResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &customerResponse.Customers[0], httpResponse, nil
}

func (s *CustomerService) AssociateWithCard(customerId, cardId string) (*Card, *http.Response, error) {
	return s.client.Card.AssociateWithCustomer(cardId, customerId)
}

func (s *CustomerService) AssociateWithBankAccount(customerId, accountId string) (*BankAccount, *http.Response, error) {
	return s.client.BankAccount.AssociateWithCustomer(accountId, customerId)
}
