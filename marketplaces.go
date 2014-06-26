package balanced

import (
	"net/http"
	"time"
)

type MarketplaceService struct {
	client *Client
}

type Marketplace struct {
	Id                  string                 `json:"id"`
	Name                string                 `json:"name"`
	SupportEmailAddress string                 `json:"support_email_address"`
	SupportPhoneNumber  string                 `json:"support_phone_number"`
	InEscrow            int                    `json:"in_escrow"`
	DomainUrl           string                 `json:"domain_url"`
	Links               map[string]string      `json:"links"`
	Href                string                 `json:"href"`
	CreatedAt           *time.Time             `json:"created_at"`
	UpdatedAt           *time.Time             `json:"updated_at"`
	Production          bool                   `json:"production"`
	Meta                map[string]interface{} `json:"meta"`
	UnsettledFees       int                    `json:"unsettled_fees"`
}

type MarketplaceResponse struct {
	Marketplaces []Marketplace             `json:"marketplaces"`
	Links        *MarketplaceResponseLinks `json:"links"`
}

type MarketplaceResponseLinks struct {
	Reversals     string `json:"marketplaces.reversals"`
	Cards         string `json:"marketplaces.cards"`
	Refunds       string `json:"marketplaces.refunds"`
	BankAccounts  string `json:"marketplaces.bank_accounts"`
	Debits        string `json:"marketplaces.debits"`
	Customers     string `json:"marketplaces.customers"`
	Credits       string `json:"marketplaces.credits"`
	CardHolds     string `json:"marketplaces.card_holds"`
	OwnerCustomer string `json:"marketplaces.owner_customer"`
	Transactions  string `json:"marketplaces.transactions"`
	Callbacks     string `json:"marketplaces.callbacks"`
	Events        string `json:"marketplaces.events"`
}

func (s *MarketplaceService) Create() (*Marketplace, *http.Response, error) {
	marketplaceResponse := new(MarketplaceResponse)
	httpResponse, err := s.client.POST("/marketplaces", nil, nil, marketplaceResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &marketplaceResponse.Marketplaces[0], httpResponse, err
}
