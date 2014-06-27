package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type BankAccountService struct {
	client *Client
}

type BankAccount struct {
	AccountNumber string `json:"account_number,omitempty"`

	// "checking" or "savings"
	AccountType   string `json:"account_type,omitempty"`
	Name          string `json:"name,omitempty"`
	RoutingNumber string `json:"routing_number,omitempty"`
	*Address      `json:"address,omitempty"`

	BankName    string `json:"bank_name,omitempty"`
	CanCredit   bool   `json:"can_credit,omitempty"`
	CanDebit    bool   `json:"can_debit,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Href        string `json:"href,omitempty"`
	Id          string `json:"id,omitempty"`

	Links     *BankAccountLinks      `json:"links,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	CreatedAt *time.Time             `json:"created_at,omitempty"`
	UpdatedAt *time.Time             `json:"updated_at,omitempty"`
}

type BankAccountLinks struct {
	BankAccountVerification string `json:"bank_account_verification"`
	Customer                string `json:customer`
}

// BankAccountPage holds a paginated set of bank accounts
type BankAccountPage struct {
	BankAccounts []BankAccount
	*PaginationParams
}

type BankAccountResponse struct {
	BankAccounts []BankAccount             `json:"bank_accounts"`
	Meta         map[string]interface{}    `json:"meta"`
	Links        *BankAccountResponseLinks `json:"links"`
}

type BankAccountResponseLinks struct {
	Verification  string `json:"bank_accounts.bank_account_verification"`
	Verifications string `json:"bank_accounts.bank_account_verifications"`
	Credits       string `json:"bank_accounts.credits"`
	Customer      string `json:"bank_accounts.customer"`
	Debits        string `json:"bank_accounts.debits"`
}

func (s *BankAccountService) Create(account *BankAccount) (*BankAccount, *http.Response, error) {
	accountResponse := new(BankAccountResponse)
	httpResponse, err := s.client.POST("/bank_accounts", nil, account, accountResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &accountResponse.BankAccounts[0], httpResponse, nil
}

func (s *BankAccountService) Delete(accountId string) (bool, *http.Response, error) {
	path := fmt.Sprintf("/bank_accounts/%v", accountId)
	httpResponse, err := s.client.DELETE(path, nil, nil, nil)
	if err != nil {
		return false, httpResponse, err
	}
	code := httpResponse.StatusCode
	didDelete := 200 <= code && code < 300
	return didDelete, httpResponse, nil
}

func (s *BankAccountService) Fetch(accountId string) (*BankAccount, *http.Response, error) {
	path := fmt.Sprintf("/bank_accounts/%v", accountId)
	accountResponse := new(BankAccountResponse)
	httpResponse, err := s.client.GET(path, nil, nil, accountResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &accountResponse.BankAccounts[0], httpResponse, nil
}

func (s *BankAccountService) Update(accountId string, params map[string]interface{}) (*BankAccount, *http.Response, error) {
	path := fmt.Sprintf("/bank_accounts/%v", accountId)
	accountResponse := new(BankAccountResponse)
	httpResponse, err := s.client.PUT(path, nil, params, accountResponse)

	if err != nil {
		return nil, httpResponse, err
	}
	return &accountResponse.BankAccounts[0], httpResponse, nil
}

func (s *BankAccountService) UpdateMeta(accountId string, meta map[string]interface{}) (*BankAccount, *http.Response, error) {
	return s.Update(accountId, map[string]interface{}{"meta": meta})
}

func (s *BankAccountService) List(args ...interface{}) (*BankAccountPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	accountResponse := new(BankAccountResponse)
	httpResponse, err := s.client.GET("/bank_accounts", query, nil, accountResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &BankAccountPage{
		BankAccounts:     accountResponse.BankAccounts,
		PaginationParams: NewPaginationParams(accountResponse.Meta),
	}, httpResponse, nil
}

func (s *BankAccountService) AssociateWithCustomer(accountId string, customerId string) (*BankAccount, *http.Response, error) {
	return s.Update(accountId, map[string]interface{}{
		"customer": fmt.Sprintf("/customers/%v", customerId),
	})
}

func (s *BankAccountService) Debit(accountId string, debit *Debit) (*Debit, *http.Response, error) {
	path := fmt.Sprintf("/bank_accounts/%v/debits", accountId)
	debitResponse := new(DebitResponse)
	httpResponse, err := s.client.POST(path, nil, &DebitRequest{
		Debits: []Debit{*debit},
	}, debitResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &debitResponse.Debits[0], httpResponse, nil
}

func (s *BankAccountService) Credit(bankAccountId string, credit *Credit) (*Credit, *http.Response, error) {
	path := fmt.Sprintf("/bank_accounts/%v/credits", bankAccountId)
	creditResponse := new(CreditResponse)
	httpResponse, err := s.client.POST(path, nil, credit, creditResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &creditResponse.Credits[0], httpResponse, nil
}
