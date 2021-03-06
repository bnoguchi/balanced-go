package balanced

import (
	"fmt"
	"net/http"
	"time"
)

// Bank account verifications work by sending microdeposits to the bank account
// which will be used to verify bank account ownership when supplied during a
// bank account verification confirmation.

// Bank accounts verifications verify the ownership of a bank account. Creating
// a verification triggers 2 micro-deposits of random amounts less than $1 into
// the target account. These amounts typically show up within 2 business days
// as "Balanced Verification" on the account owner's statement. After obtaining
// these amounts, the account owner should submit these to the Balanced API to
// confirm ownership. Balanced allows 3 attempts to enter the correct
// verification amounts. After 3 failed attempts, a new verification must be
// created. Only one verification may exist at a time for a given bank account.
// Verifications are *not* required for accounts to which *only* credits will
// be issued.
type VerificationService struct {
	client *Client
}

type Verification struct {
	Attempts           int                `json:"attempts"`           // e.g., 0
	AttemptsRemaining  int                `json:"attempts_remaining"` // e.g., 3
	DepositStatus      string             `json:"deposit_status"`     // e.g., "succeeded"
	Href               string             `json"href"`                // e.g., "/verifications/BZ25cVCn6wh6UZrfgFcF71RD"
	Id                 string             `json"id"`                  // e.g., "BZ25cVCn6wh6UZrfgFcF71RD"
	Links              *VerificationLinks `json:"links"`              // e.g., "bank_account" => "BA1RdDM12aF5N8WVA1kaewQZ"
	Meta               map[string]string  `json:"meta"`
	VerificationStatus string             `json:"verification_status"` // e.g., "pending"
	CreatedAt          *time.Time         `json:"created_at"`
	UpdatedAt          *time.Time         `json:"updated_at"`
}

type VerificationLinks struct {
	BankAccount string `json:"bank_account"`
}

type verificationResponse struct {
	Verifications []Verification             `json:"bank_account_verifications"`
	Meta          map[string]interface{}     `json:"meta"`
	Links         *verificationResponseLinks `json:"links"`
}

type verificationResponseLinks struct {
	BankAccount string `json:"bank_account_verifications.bank_account"`
}

func (s *VerificationService) Create(accountId string) (*Verification, *http.Response, error) {
	path := fmt.Sprintf("/bank_accounts/%v/verifications", accountId)
	verifResponse := new(verificationResponse)
	httpResponse, err := s.client.POST(path, nil, nil, verifResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &verifResponse.Verifications[0], httpResponse, nil
}

// Fetches the verification for a bank account
func (s *VerificationService) Fetch(verificationId string) (*Verification, *http.Response, error) {
	path := fmt.Sprintf("/verifications/%v", verificationId)
	verifResponse := new(verificationResponse)
	httpResponse, err := s.client.GET(path, nil, nil, verifResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &verifResponse.Verifications[0], httpResponse, nil
}

// Confirmation amounts are sent with an attempt to confirm a bank account
// verification
type ConfirmationAmounts struct {
	Amount1 int `json:"amount_1"`
	Amount2 int `json:"amount_2"`
}

func (s *VerificationService) Confirm(verificationId string, amount1 int, amount2 int) (*Verification, *http.Response, error) {
	path := fmt.Sprintf("/verifications/%v", verificationId)
	verifResponse := new(verificationResponse)
	httpResponse, err := s.client.PUT(path, nil, &ConfirmationAmounts{
		Amount1: amount1,
		Amount2: amount2,
	}, verifResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &verifResponse.Verifications[0], httpResponse, nil
}
