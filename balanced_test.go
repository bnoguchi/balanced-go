package balanced

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

var sharedClient *Client

type CardFixture struct {
	Brand  string
	Number string
	Cvv    string

	// "Success", "Processor Failure", "Tokenization Error", "Cvv Match Fail",
	// "Cvv Unsupported", "Disputed Charge"
	Result string
}

func newCardFixture(number string, cvv string) *Card {
	return &Card{
		Name:            "Brian Noguchi",
		ExpirationMonth: 12,
		ExpirationYear:  2016,
		Number:          number,
		Cvv:             cvv,
	}
}

var cardFixtures map[string]*Card = map[string]*Card{
	"VisaSuccess": &Card{
		Number: "4111111111111111",
		Cvv:    "123",
	},
	"MasterCardSuccess": &Card{
		Number: "5105105105105100",
		Cvv:    "123",
	},
	"AmexSuccess": &Card{
		Number: "341111111111111",
		Cvv:    "1234",
	},
	"VisaCreditable": newCardFixture(
		"4342561111111118",
		"123",
	),
	"VisaProcessorFailure": &Card{
		Number: "4444444444444448",
		Cvv:    "123",
	},
	"VisaTokenizationError": &Card{
		Number: "4222222222222220",
		Cvv:    "123",
	},
	"MasterCardCvvFail": &Card{
		Number: "5112000200000002",
		Cvv:    "200",
	},
	"VisaCvvUnsupported": &Card{
		Number: "4457000300000007",
		Cvv:    "901",
	},
	"DiscoverDisputedCharge": &Card{
		Number: "6500000000000002",
		Cvv:    "123",
	},
}

var bankAccountFixtures map[string]*BankAccount = map[string]*BankAccount{
	"invalid_routing_a": &BankAccount{
		RoutingNumber: "100000007",
		AccountNumber: "8887776665555",
	},
	"invalid_routing_b": &BankAccount{
		RoutingNumber: "111111118",
		AccountNumber: "8887776665555",
	},
	"pending_a": &BankAccount{
		RoutingNumber: "021000021",
		AccountNumber: "9900000000",
	},
	"pending_b": &BankAccount{
		RoutingNumber: "321174851",
		AccountNumber: "9900000001",
	},
	"succeeded_a": &BankAccount{
		RoutingNumber: "021000021",
		AccountNumber: "9900000002",
	},
	"succeeded_b": &BankAccount{
		RoutingNumber: "321174851",
		AccountNumber: "9900000003",
	},
	"failed_a": &BankAccount{
		RoutingNumber: "021000021",
		AccountNumber: "9900000004",
	},
	"failed_b": &BankAccount{
		RoutingNumber: "321174851",
		AccountNumber: "9900000005",
	},
}

func init() {
	apiKey, err := createApiKey()
	if err != nil {
		panic(err)
	}
	if apiKey == nil {
		panic("Expected apiKey to not be nil")
	}
	secret := apiKey.Secret
	if secret == "" {
		panic("Expected non-empty secret")
	}
	sharedClient = NewClient(nil, secret)

	// Setup a test marketplace
	marketplace, _, err := sharedClient.Marketplace.Create()
	if err != nil {
		panic(err)
	}
	if marketplace.Production != false {
		panic("Tests need to be run on a test marketplace. This marketplace is a production marketplace")
	}
}

// Hook up gocheck into the "go test" runner
func Test(t *testing.T) { TestingT(t) }

var testSecret = ""

func createApiKey() (*ApiKey, error) {
	c := NewClient(nil, testSecret)
	apiKey, _, err := c.ApiKey.Create()
	return apiKey, err
}

func TestNewMarketplace(t *testing.T) {
	apiKey, err := createApiKey()
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if apiKey == nil {
		t.Error("Expected apiKey to not be nil")
	}
	secret := apiKey.Secret
	if secret == "" {
		t.Error("Expected non-empty. Got an empty secret.")
	}
}

type CardSuite struct{}

var _ = Suite(&CardSuite{})

func createCardFixture(client *Client, fixtureAlias string) (*Card, error) {
	card, _, err := client.Card.Create(cardFixtures[fixtureAlias])
	return card, err
}

func mustCreateCardFixture(client *Client, fixtureAlias string) *Card {
	card, err := createCardFixture(client, fixtureAlias)
	if err != nil {
		panic(err)
	}
	return card
}

func createCard(client *Client) (*Card, error) {
	card, _, err := client.Card.Create(&Card{
		ExpirationMonth: 12,
		ExpirationYear:  2016,
		Number:          "4111111111111111",
	})
	return card, err
}

func mustCreateCard(client *Client) *Card {
	card, err := createCard(client)
	if err != nil {
		panic(err)
	}
	return card
}

func deleteCard(client *Client, card *Card, c *C) {
	didDelete, _, err := client.Card.Delete(card.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)
}

func (s *CardSuite) TestCreate(c *C) {
	card, err := createCard(sharedClient)
	c.Assert(err, IsNil)
	defer deleteCard(sharedClient, card, c)

	c.Assert(card.Number, Equals, "xxxxxxxxxxxx1111")
	c.Assert(card.ExpirationMonth, Equals, 12)
	c.Assert(card.ExpirationYear, Equals, 2016)
}

func (s *CardSuite) TestDelete(c *C) {
	startingCardPage, _, err := sharedClient.Card.List()
	c.Assert(err, IsNil)

	card := mustCreateCard(sharedClient)
	cardPage, _, err := sharedClient.Card.List()
	c.Assert(err, IsNil)
	c.Assert(len(cardPage.Cards), Equals, 1+startingCardPage.Total)

	didDelete, _, err := sharedClient.Card.Delete(card.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)

	cardPage, _, err = sharedClient.Card.List()
	c.Assert(err, IsNil)
	c.Assert(len(cardPage.Cards), Equals, startingCardPage.Total)

	// Deleted cards still resolve their hrefs, even though they do not show up
	// in a list of cards.
	fetchedCard, _, err := sharedClient.Card.Fetch(card.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedCard, Not(IsNil))
	c.Assert(fetchedCard.Id, Equals, card.Id)
}

func (s *CardSuite) TestFetch(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	fetchedCard, _, err := sharedClient.Card.Fetch(card.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedCard.Number, Equals, "xxxxxxxxxxxx1111")
	c.Assert(fetchedCard.ExpirationMonth, Equals, 12)
	c.Assert(fetchedCard.ExpirationYear, Equals, 2016)
}

func (s *CardSuite) TestList(c *C) {
	cardPage, _, err := sharedClient.Card.List()
	c.Assert(err, IsNil)
	originalTotal := cardPage.Total

	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	cardPage, _, err = sharedClient.Card.List()
	c.Assert(err, IsNil)
	c.Assert(cardPage.Total, Equals, originalTotal+1)
	c.Assert(cardPage.Limit, Equals, 10)
	c.Assert(cardPage.Offset, Equals, 0)
	c.Assert(cardPage.Cards[0].Id, Equals, card.Id)
}

func (s *CardSuite) TestUpdate(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	// Should be able to update meta
	updatedCard, _, err := sharedClient.Card.Update(card.Id, map[string]interface{}{
		"meta": map[string]string{
			"twitter.id": "1234987650",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(updatedCard.Meta["twitter.id"], Equals, "1234987650")

	// TODO Should be able to add new primary card info (e.g., name)
}

func (s *CardSuite) TestAssociateWithCustomer(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	c.Assert(card.Links["customer"], Equals, "")

	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	updatedCard, _, err := sharedClient.Card.AssociateWithCustomer(card.Id, customer.Id)

	c.Assert(err, IsNil)
	c.Assert(updatedCard.Links["customer"], Equals, customer.Id)
}

func (s *CardSuite) TestCharge(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount:               50,
		AppearsOnStatementAs: "Starbucks Gift Card",
		Description:          "Test Charge",
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 50)
}

func (s *CardSuite) TestCredit(c *C) {
	card := mustCreateCardFixture(sharedClient, "VisaCreditable")
	defer deleteCard(sharedClient, card, c)

	// Debit funds from a card into Escrow, so we can have money to pay out
	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 70,
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 70)

	credit, _, err := sharedClient.Card.Credit(card.Id, &Credit{
		Amount:               70,
		AppearsOnStatementAs: "Credit back",
		Description:          "Test Credit",
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 70)
}

func createCustomer(client *Client) (*Customer, error) {
	customer, _, err := client.Customer.Create(&Customer{})
	return customer, err
}

func mustCreateCustomer(client *Client) *Customer {
	customer, err := createCustomer(client)
	if err != nil {
		panic(err)
	}
	return customer
}

func deleteCustomer(client *Client, customerId string, c *C) {
	didDelete, _, err := client.Customer.Delete(customerId)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)
}

type CustomerSuite struct{}

var _ = Suite(&CustomerSuite{})

func (s *CustomerSuite) TestCreate(c *C) {
	customer, err := createCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	c.Assert(err, IsNil)
	c.Assert(customer, Not(IsNil))
	c.Assert(customer.Id, Not(Equals), "")
}

func (s *CustomerSuite) TestDelete(c *C) {
	startingCustomerPage, _, err := sharedClient.Customer.List()

	customer := mustCreateCustomer(sharedClient)

	didDelete, _, err := sharedClient.Customer.Delete(customer.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)

	finalCustomerPage, _, err := sharedClient.Customer.List()
	c.Assert(err, IsNil)
	c.Assert(finalCustomerPage.Total, Equals, startingCustomerPage.Total)

	// TODO Show that fetch still brings up customer
}

func (s *CustomerSuite) TestFetch(c *C) {
	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	fetchedCustomer, _, err := sharedClient.Customer.Fetch(customer.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedCustomer, Not(IsNil))
	c.Assert(fetchedCustomer.Id, Equals, customer.Id)
}

func (s *CustomerSuite) TestList(c *C) {
	customerPage, _, err := sharedClient.Customer.List()
	c.Assert(err, IsNil)
	originalTotal := customerPage.Total

	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	customerPage, _, err = sharedClient.Customer.List()
	c.Assert(err, IsNil)
	c.Assert(customerPage.Total, Equals, originalTotal+1)
	c.Assert(customerPage.Limit, Equals, 10)
	c.Assert(customerPage.Offset, Equals, 0)
	c.Assert(customerPage.Customers[0].Id, Equals, customer.Id)
}

func (s *CustomerSuite) TestUpdate(c *C) {
	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	// Should be able to update meta
	updatedCustomer, _, err := sharedClient.Customer.Update(customer.Id, map[string]interface{}{
		"meta": map[string]string{
			"shipping-preference": "ground",
		},
		"email": "email@newdomain.com",
	})
	c.Assert(err, IsNil)
	c.Assert(updatedCustomer.Meta["shipping-preference"], Equals, "ground")
	c.Assert(updatedCustomer.Email, Equals, "email@newdomain.com")
}

func (s *CustomerSuite) TestAssociateWithCard(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	c.Assert(card.Links["customer"], Equals, "")

	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	updatedCard, _, err := sharedClient.Customer.AssociateWithCard(customer.Id, card.Id)

	c.Assert(err, IsNil)
	c.Assert(updatedCard.Links["customer"], Equals, customer.Id)
}

func (s *CustomerSuite) TestAssociateWithBankAccount(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	c.Assert(account.Links.Customer, Equals, "")

	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	updatedAccount, _, err := sharedClient.Customer.AssociateWithBankAccount(customer.Id, account.Id)

	c.Assert(err, IsNil)
	c.Assert(updatedAccount.Links.Customer, Equals, customer.Id)
}

type BankAccountSuite struct{}

var _ = Suite(&BankAccountSuite{})

func createBankAccount(client *Client, fixture *BankAccount) (*BankAccount, error) {
	if fixture == nil {
		fixture = bankAccountFixtures["succeeded_a"]
	}
	fixture.Name = "Test Name"
	fixture.AccountType = "checking"
	account, _, err := client.BankAccount.Create(fixture)
	return account, err
}

func mustCreateBankAccount(client *Client, fixture *BankAccount) *BankAccount {
	account, err := createBankAccount(client, fixture)
	if err != nil {
		panic(err)
	}
	return account
}

func deleteBankAccount(client *Client, account *BankAccount, c *C) {
	didDelete, _, err := client.BankAccount.Delete(account.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)
}

func (s *BankAccountSuite) TestCreate(c *C) {
	account, err := createBankAccount(sharedClient, nil)
	c.Assert(err, IsNil)
	defer deleteBankAccount(sharedClient, account, c)

	c.Assert(account.AccountNumber, Equals, "xxxxxx0002")
	c.Assert(account.RoutingNumber, Equals, "021000021")
	c.Assert(account.Name, Equals, "Test Name")
	c.Assert(account.AccountType, Equals, "checking")
}

func (s *BankAccountSuite) TestDelete(c *C) {
	startingAccountPage, _, err := sharedClient.BankAccount.List()

	account := mustCreateBankAccount(sharedClient, nil)
	didDelete, _, err := sharedClient.BankAccount.Delete(account.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)

	accountPage, _, err := sharedClient.BankAccount.List()
	c.Assert(err, IsNil)
	c.Assert(len(accountPage.BankAccounts), Equals, len(startingAccountPage.BankAccounts))
	c.Assert(accountPage.Total, Equals, startingAccountPage.Total)

	// Deleted accounts still resolve their hrefs, even though they do not show
	// up in a list of bank accounts.
	fetchedAccount, _, err := sharedClient.BankAccount.Fetch(account.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedAccount, Not(IsNil))
	c.Assert(fetchedAccount.Id, Equals, account.Id)
}

func (s *BankAccountSuite) TestFetch(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	fetchedAccount, _, err := sharedClient.BankAccount.Fetch(account.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedAccount.Id, Equals, account.Id)
}

func (s *BankAccountSuite) TestList(c *C) {
	startingAccountPage, _, err := sharedClient.BankAccount.List()
	c.Assert(err, IsNil)

	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	accountPage, _, err := sharedClient.BankAccount.List()
	c.Assert(err, IsNil)
	c.Assert(accountPage.Total, Equals, 1+startingAccountPage.Total)
	c.Assert(accountPage.Offset, Equals, 0)
	c.Assert(accountPage.Limit, Equals, 10)
	currAccounts := accountPage.BankAccounts
	c.Assert(currAccounts[0].Id, Equals, account.Id)
}

func (s *BankAccountSuite) TestUpdateMeta(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	// Should be able to update meta
	updatedAccount, _, err := sharedClient.BankAccount.UpdateMeta(account.Id, map[string]interface{}{
		"twitter.id": "1234987650",
	})
	c.Assert(err, IsNil)
	c.Assert(updatedAccount.Meta["twitter.id"], Equals, "1234987650")
}

func (s *BankAccountSuite) TestAssociateWithCustomer(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	c.Assert(account.Links.Customer, Equals, "")

	customer := mustCreateCustomer(sharedClient)
	defer deleteCustomer(sharedClient, customer.Id, c)

	updatedAccount, _, err := sharedClient.BankAccount.AssociateWithCustomer(account.Id, customer.Id)

	c.Assert(err, IsNil)
	c.Assert(updatedAccount.Links.Customer, Equals, customer.Id)
}

func (s *BankAccountSuite) TestCredit(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	// Debit funds from a card into Escrow, so we can have money to pay out
	card := mustCreateCard(sharedClient)
	sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 50,
	})

	credit, _, err := sharedClient.BankAccount.Credit(account.Id, &Credit{
		Amount: 50,
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 50)
	c.Assert(credit.Status, Equals, Succeeded)
}

func mustVerifyAccount(accountId string) {
	verif, _, err := sharedClient.Verification.Create(accountId)
	if err != nil {
		panic(err)
	}
	_, _, err = sharedClient.Verification.Confirm(verif.Id, 1, 1)
	if err != nil {
		panic(err)
	}
}

func (s *BankAccountSuite) TestDebit(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	// Debit funds from a card into Escrow, so we can have money to pay out
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)
	sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 10000,
	})

	// Credit funds to the bank account that we'll later debit
	_, _, err := sharedClient.BankAccount.Credit(account.Id, &Credit{
		Amount: 10000,
	})
	c.Assert(err, IsNil)

	// An account must be verified before it can be successfully debited
	mustVerifyAccount(account.Id)

	debit, _, err := sharedClient.BankAccount.Debit(account.Id, &Debit{
		Amount:               10000,
		AppearsOnStatementAs: "Starbucks Gift Card",
		Description:          "Test Charge",
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 10000)
}

// TODO Test debit error due to insufficient funds in account

func (s *BankAccountSuite) TestDebitUnverifiedAccount(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	c.Assert(account.CanDebit, Equals, false)

	debit, res, err := sharedClient.BankAccount.Debit(account.Id, &Debit{
		Amount:               50,
		AppearsOnStatementAs: "Starbucks Gift Card",
		Description:          "Test Charge",
	})
	c.Assert(debit, IsNil)
	c.Assert(err, Not(IsNil))
	// TODO Assert expected error details
	c.Assert(debit, IsNil)
	c.Assert(res.StatusCode, Equals, 409)
	c.Assert(err.(*ErrorResponse).Errors[0].CategoryCode, Equals, "funding-source-not-debitable")
}

type VerificationSuite struct{}

var _ = Suite(&VerificationSuite{})

func (s *VerificationSuite) TestCreate(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	verif, _, err := sharedClient.Verification.Create(account.Id)
	c.Assert(err, IsNil)
	c.Assert(verif.Attempts, Equals, 0)
	c.Assert(verif.AttemptsRemaining, Equals, 3)
	c.Assert(verif.DepositStatus, Equals, Succeeded)
	c.Assert(verif.Href, Equals, "/verifications/"+verif.Id)
	c.Assert(verif.Links["bank_account"], Equals, account.Id)
	c.Assert(verif.VerificationStatus, Equals, Pending)
}

func (s *VerificationSuite) TestFetch(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	verif, _, err := sharedClient.Verification.Create(account.Id)
	c.Assert(err, IsNil)

	fetchedVerif, _, err := sharedClient.Verification.Fetch(verif.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedVerif.Id, Equals, verif.Id)
}

func (s *VerificationSuite) TestConfirm(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	verif, _, err := sharedClient.Verification.Create(account.Id)
	c.Assert(verif.Attempts, Equals, 0)
	c.Assert(verif.AttemptsRemaining, Equals, 3)
	c.Assert(verif.VerificationStatus, Equals, Pending)

	confirmedVerif, _, err := sharedClient.Verification.Confirm(verif.Id, 1, 1)
	c.Assert(err, IsNil)
	c.Assert(confirmedVerif.Attempts, Equals, 1)
	c.Assert(confirmedVerif.AttemptsRemaining, Equals, 2)
	c.Assert(confirmedVerif.VerificationStatus, Equals, Succeeded)
}

func (s *VerificationSuite) TestConfirmFail(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	verif, _, err := sharedClient.Verification.Create(account.Id)
	c.Assert(verif.Attempts, Equals, 0)
	c.Assert(verif.AttemptsRemaining, Equals, 3)
	c.Assert(verif.VerificationStatus, Equals, Pending)

	confirmedVerif, _, err := sharedClient.Verification.Confirm(verif.Id, 2, 2)
	c.Assert(confirmedVerif, IsNil)
	expectedErr := fmt.Sprintf("PUT https://api.balancedpayments.com/verifications/%v: 409 Authentication amounts do not match. Your request id is .*", verif.Id)
	c.Assert(err, ErrorMatches, expectedErr)
	fetchedVerif, _, err := sharedClient.Verification.Fetch(verif.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedVerif.Attempts, Equals, 1)
	c.Assert(fetchedVerif.AttemptsRemaining, Equals, 2)
	c.Assert(fetchedVerif.VerificationStatus, Equals, Pending)
}

type CallbackSuite struct{}

var _ = Suite(&CallbackSuite{})

func createCallback(client *Client) (*Callback, error) {
	callback, _, err := client.Callback.Create("http://requestb.in/122e0yu1", "get")
	return callback, err
}

func mustCreateCallback(client *Client) *Callback {
	callback, err := createCallback(client)
	if err != nil {
		panic(err)
	}
	return callback
}

func deleteCallback(sharedClient *Client, callback *Callback, c *C) {
	didDelete, _, err := sharedClient.Callback.Delete(callback.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)
}

func (s *CallbackSuite) TestCreate(c *C) {
	callback := mustCreateCallback(sharedClient)
	defer deleteCallback(sharedClient, callback, c)

	c.Assert(callback.Url, Equals, "http://requestb.in/122e0yu1")
	// TODO Test callback
}

func (s *CallbackSuite) TestDelete(c *C) {
	startingCallbackPage, _, err := sharedClient.Callback.List()
	c.Assert(err, IsNil)

	callback := mustCreateCallback(sharedClient)
	callbackPage, _, err := sharedClient.Callback.List()
	c.Assert(err, IsNil)
	c.Assert(len(callbackPage.Callbacks), Equals, 1+startingCallbackPage.Total)

	didDelete, _, err := sharedClient.Callback.Delete(callback.Id)
	c.Assert(err, IsNil)
	c.Assert(didDelete, Equals, true)

	callbackPage, _, err = sharedClient.Callback.List()
	c.Assert(err, IsNil)
	c.Assert(callbackPage.Total, Equals, startingCallbackPage.Total)

	fetchedCallback, _, err := sharedClient.Callback.Fetch(callback.Id)
	c.Assert(fetchedCallback, IsNil)
	c.Assert(err, ErrorMatches, ".*404.*The requested URL was not found on the server.*")
}

func (s *CallbackSuite) TestFetch(c *C) {
	callback := mustCreateCallback(sharedClient)
	defer deleteCallback(sharedClient, callback, c)

	fetchedCallback, _, err := sharedClient.Callback.Fetch(callback.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedCallback.Url, Equals, callback.Url)
	c.Assert(fetchedCallback.Id, Equals, callback.Id)
}

func (s *CallbackSuite) TestList(c *C) {
	startingCallbackPage, _, err := sharedClient.Callback.List()
	c.Assert(err, IsNil)

	callback := mustCreateCallback(sharedClient)
	defer deleteCallback(sharedClient, callback, c)

	callbackPage, _, err := sharedClient.Callback.List()
	c.Assert(err, IsNil)
	c.Assert(callbackPage.Total, Equals, 1+startingCallbackPage.Total)
	c.Assert(callbackPage.Limit, Equals, 10)
	c.Assert(callbackPage.Offset, Equals, 0)
	c.Assert(callbackPage.Callbacks[0].Id, Equals, callback.Id)
}

type CardHoldSuite struct{}

var _ = Suite(&CardHoldSuite{})

func createHold(client *Client, card *Card) (*CardHold, error) {
	hold, _, err := client.CardHold.Create(card.Id, &CardHold{
		Amount: 100,
	})
	return hold, err
}

func mustCreateHold(client *Client, card *Card) *CardHold {
	hold, err := createHold(client, card)
	if err != nil {
		panic(err)
	}
	return hold
}

func voidHold(sharedClient *Client, hold *CardHold, c *C) {
	voidedHold, _, err := sharedClient.CardHold.Void(hold.Id)
	c.Assert(err, IsNil)
	c.Assert(voidedHold.VoidedAt, Not(IsNil))
}

func (s *CardHoldSuite) TestCreate(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	hold := mustCreateHold(sharedClient, card)
	defer voidHold(sharedClient, hold, c)

	c.Assert(hold.Amount, Equals, 100)
}

func (s *CardHoldSuite) TestFetch(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	hold := mustCreateHold(sharedClient, card)
	defer voidHold(sharedClient, hold, c)

	fetchedHold, _, err := sharedClient.CardHold.Fetch(hold.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedHold.Id, Equals, hold.Id)
}

func (s *CardHoldSuite) TestList(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	startingHoldPage, _, err := sharedClient.CardHold.List()
	c.Assert(err, IsNil)

	hold := mustCreateHold(sharedClient, card)
	defer voidHold(sharedClient, hold, c)

	holdPage, _, err := sharedClient.CardHold.List()
	c.Assert(err, IsNil)
	c.Assert(holdPage.Total, Equals, 1+startingHoldPage.Total)
}

func (s *CardHoldSuite) TestUpdate(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	hold := mustCreateHold(sharedClient, card)
	defer voidHold(sharedClient, hold, c)

	updatedHold, _, err := sharedClient.CardHold.Update(hold.Id, map[string]interface{}{
		"description": "Sample description",
	})
	c.Assert(err, IsNil)
	c.Assert(updatedHold.Description, Equals, "Sample description")

	updatedHold, _, err = sharedClient.CardHold.Update(hold.Id, map[string]interface{}{
		"meta": map[string]interface{}{
			"xxx": "yyy",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(updatedHold.Meta["xxx"], Equals, "yyy")
}

func (s *CardHoldSuite) TestCapture(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	hold := mustCreateHold(sharedClient, card)

	debit, _, err := sharedClient.CardHold.Capture(hold.Id, &Debit{
		Amount:               100,
		AppearsOnStatementAs: "Test Capture",
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 100)
	c.Assert(debit.AppearsOnStatementAs, Equals, "BAL*Test Capture")
}

func (s *CardHoldSuite) TestVoid(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	hold := mustCreateHold(sharedClient, card)
	defer voidHold(sharedClient, hold, c)

	c.Assert(hold.VoidedAt, IsNil)

	voidedHold, _, err := sharedClient.CardHold.Void(hold.Id)
	c.Assert(err, IsNil)
	c.Assert(voidedHold.VoidedAt, Not(IsNil))

	debit, _, err := sharedClient.CardHold.Capture(hold.Id, &Debit{
		Amount: 100,
	})
	c.Assert(debit, IsNil)
	c.Assert(err, ErrorMatches, ".* 409 This hold .* has already been voided.*")
}

func (s *CardHoldSuite) TestVoidCaptured(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	hold := mustCreateHold(sharedClient, card)
	c.Assert(hold.VoidedAt, IsNil)

	debit, _, err := sharedClient.CardHold.Capture(hold.Id, &Debit{
		Amount: 100,
	})
	c.Assert(debit.Amount, Equals, 100)
	c.Assert(err, IsNil)

	failedVoidedHold, _, err := sharedClient.CardHold.Void(hold.Id)
	c.Assert(failedVoidedHold, IsNil)
	c.Assert(err, ErrorMatches, ".* 409 This hold .* has already been captured.*")
}

type CreditSuite struct{}

var _ = Suite(&CreditSuite{})

func (s *CreditSuite) TestCreateToBankAccount(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	// Debit funds from a card into Escrow, so we can have money to pay out
	card := mustCreateCard(sharedClient)
	sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 50,
	})

	credit, _, err := sharedClient.Credit.CreateToBankAccount(account.Id, &Credit{
		Amount: 50,
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 50)
	c.Assert(credit.Status, Equals, Succeeded)
}

func (s *CreditSuite) TestCreateToCard(c *C) {
	card := mustCreateCardFixture(sharedClient, "VisaCreditable")
	defer deleteCard(sharedClient, card, c)

	// Debit funds from a card into Escrow, so we can have money to pay out
	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 70,
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 70)

	credit, _, err := sharedClient.Credit.CreateToCard(card.Id, &Credit{
		Amount:               70,
		AppearsOnStatementAs: "Credit back",
		Description:          "Test Credit",
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 70)
}

func (s *CreditSuite) TestCreateForOrder(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	customer := mustCreateCustomer(sharedClient)
	// Can't delete a customer associated with an order
	// defer deleteCustomer(sharedClient, customer.Id, c)

	sharedClient.BankAccount.AssociateWithCustomer(account.Id, customer.Id)

	order, _, err := sharedClient.Order.Create(customer.Id, &Order{
		Description: "Order #TestCreateForOrder",
	})
	c.Assert(err, IsNil)

	// Debit funds from a card into the order escrow, so we can have money to pay
	// out. Aka create a debit for an order.
	card := mustCreateCard(sharedClient)
	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 50,
		Order:  order.Href,
	})
	c.Assert(debit.Amount, Equals, 50)
	c.Assert(err, IsNil)

	credit, _, err := sharedClient.Credit.CreateForOrder(account.Id, order.Id, &Credit{
		Amount: 50,
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 50)
	c.Assert(credit.Status, Equals, Succeeded)
}

func (s *CreditSuite) TestList(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	startingCreditPage, _, err := sharedClient.Credit.List()
	c.Assert(err, IsNil)

	credit, _, err := sharedClient.BankAccount.Credit(account.Id, &Credit{
		Amount: 50,
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 50)

	creditPage, _, err := sharedClient.Credit.List()
	c.Assert(err, IsNil)
	c.Assert(creditPage.Total, Equals, 1+startingCreditPage.Total)
}

func (s *CreditSuite) TestListForBankAccount(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	startingCreditPage, _, err := sharedClient.Credit.ListForBankAccount(account.Id)
	c.Assert(err, IsNil)

	credit, _, err := sharedClient.BankAccount.Credit(account.Id, &Credit{
		Amount: 50,
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 50)

	// We should list credits associated with the account

	creditPage, _, err := sharedClient.Credit.ListForBankAccount(account.Id)
	c.Assert(err, IsNil)
	c.Assert(creditPage.Total, Equals, 1+startingCreditPage.Total)
	c.Assert(creditPage.Total, Equals, 1)

	// We should not list credits not associated with the account

	ignoreAccount := mustCreateBankAccount(sharedClient, bankAccountFixtures["succeeded_b"])
	defer deleteBankAccount(sharedClient, ignoreAccount, c)

	ignoreCredit, _, err := sharedClient.BankAccount.Credit(ignoreAccount.Id, &Credit{
		Amount: 50,
	})
	c.Assert(err, IsNil)
	c.Assert(ignoreCredit.Amount, Equals, 50)

	creditPage, _, err = sharedClient.Credit.ListForBankAccount(account.Id)
	c.Assert(err, IsNil)
	c.Assert(creditPage.Total, Equals, 1)
}

func (s *CreditSuite) TestUpdate(c *C) {
	card := mustCreateCardFixture(sharedClient, "VisaCreditable")
	defer deleteCard(sharedClient, card, c)

	// Debit funds from a card into Escrow, so we can have money to pay out
	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 70,
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 70)

	credit, _, err := sharedClient.Credit.CreateToCard(card.Id, &Credit{
		Amount:               70,
		AppearsOnStatementAs: "Credit back",
		Description:          "Test Credit",
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 70)

	updatedCredit, _, err := sharedClient.Credit.Update(credit.Id, map[string]interface{}{
		"description": "updated description",
		"meta": map[string]interface{}{
			"xxx": "yyy",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(updatedCredit.Description, Equals, "updated description")
	c.Assert(updatedCredit.Meta["xxx"], Equals, "yyy")
}

func (s *CreditSuite) TestFilter(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	startingCreditPage, _, err := sharedClient.Credit.List(map[string]interface{}{
		"amount[>]": 1000,
		"amount[<]": 2000,
	})
	c.Assert(err, IsNil)

	amounts := [4]int{50, 1100, 1900, 2050}

	for _, amount := range amounts {
		credit, _, err := sharedClient.BankAccount.Credit(account.Id, &Credit{
			Amount: amount,
		})
		c.Assert(err, IsNil)
		c.Assert(credit, Not(IsNil))
	}

	creditPage, _, err := sharedClient.Credit.List(map[string]interface{}{
		"amount[>]": 1000,
		"amount[<]": 2000,
	})

	c.Assert(err, IsNil)
	c.Assert(creditPage.Total, Equals, 2+startingCreditPage.Total)
}

type DebitSuite struct{}

var _ = Suite(&DebitSuite{})

func (s *DebitSuite) TestFetch(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount:               50,
		AppearsOnStatementAs: "Starbucks Gift Card",
		Description:          "Test Charge",
	})
	c.Assert(err, IsNil)
	c.Assert(debit, Not(IsNil))

	fetchedDebit, _, err := sharedClient.Debit.Fetch(debit.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedDebit.Id, Equals, debit.Id)
	c.Assert(fetchedDebit.Amount, Equals, 50)
}

func (s *DebitSuite) TestList(c *C) {
	startingDebitList, _, err := sharedClient.Debit.List()
	c.Assert(err, IsNil)

	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount:               50,
		AppearsOnStatementAs: "Starbucks Gift Card",
		Description:          "Test Charge",
	})
	c.Assert(err, IsNil)
	c.Assert(debit, Not(IsNil))

	debitList, _, err := sharedClient.Debit.List()
	c.Assert(err, IsNil)
	c.Assert(debitList.Total, Equals, 1+startingDebitList.Total)
}

func (s *DebitSuite) TestUpdate(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount:               50,
		AppearsOnStatementAs: "Starbucks Gift Card",
		Description:          "Test Charge",
	})
	c.Assert(err, IsNil)
	c.Assert(debit, Not(IsNil))

	updatedDebit, _, err := sharedClient.Debit.Update(debit.Id, map[string]interface{}{
		"description": "updated description",
		"meta": map[string]interface{}{
			"xxx": "yyy",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(updatedDebit.Description, Equals, "updated description")
	c.Assert(updatedDebit.Meta["xxx"], Equals, "yyy")
}

func (s *DebitSuite) TestRefund(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 5000,
	})
	c.Assert(err, IsNil)

	refund, _, err := sharedClient.Debit.Refund(debit.Id, &Refund{
		Amount: 5000,
	})
	c.Assert(err, IsNil)
	c.Assert(refund.Status, Equals, Succeeded)
}

type RefundSuite struct{}

var _ = Suite(&RefundSuite{})

func (s *RefundSuite) TestFetch(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 5000,
	})
	c.Assert(err, IsNil)

	refund, _, err := sharedClient.Debit.Refund(debit.Id, &Refund{
		Amount: 5000,
	})
	c.Assert(err, IsNil)

	fetchedRefund, _, err := sharedClient.Refund.Fetch(refund.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedRefund.Id, Equals, refund.Id)
	c.Assert(fetchedRefund.Amount, Equals, refund.Amount)
}

func (s *RefundSuite) TestList(c *C) {
	startingRefundPage, _, err := sharedClient.Refund.List()
	c.Assert(err, IsNil)

	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 5000,
	})
	c.Assert(err, IsNil)

	refund, _, err := sharedClient.Debit.Refund(debit.Id, &Refund{
		Amount: 5000,
	})
	c.Assert(err, IsNil)
	c.Assert(refund, Not(IsNil))

	refundPage, _, err := sharedClient.Refund.List()
	c.Assert(err, IsNil)
	c.Assert(refundPage.Total, Equals, 1+startingRefundPage.Total)
}

func (s *RefundSuite) TestUpdate(c *C) {
	card := mustCreateCard(sharedClient)
	defer deleteCard(sharedClient, card, c)

	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 5000,
	})
	c.Assert(err, IsNil)

	refund, _, err := sharedClient.Debit.Refund(debit.Id, &Refund{
		Amount: 5000,
	})
	c.Assert(err, IsNil)

	updatedRefund, _, err := sharedClient.Refund.Update(refund.Id, map[string]interface{}{
		"description": "sample description",
		"meta": map[string]interface{}{
			"xxx": "yyy",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(updatedRefund.Description, Equals, "sample description")
	c.Assert(updatedRefund.Meta["xxx"], Equals, "yyy")
}

type DisputeSuite struct{}

var _ = Suite(&DisputeSuite{})

func (s *DisputeSuite) TestFetch(c *C) {
	c.Skip("Unimplemented")
}

func (s *DisputeSuite) TestList(c *C) {
	c.Skip("Unimplemented")
}

type OrderSuite struct{}

var _ = Suite(&OrderSuite{})

func (s *OrderSuite) TestCreate(c *C) {
	customer := mustCreateCustomer(sharedClient)

	// Can't delete a customer associated with an order
	// defer deleteCustomer(sharedClient, customer.Id, c)

	order, _, err := sharedClient.Order.Create(customer.Id, &Order{
		Description: "TestCreate",
	})

	c.Assert(err, IsNil)
	c.Assert(order.Description, Equals, "TestCreate")
	c.Assert(order.AmountEscrowed, Equals, 0)
}

func (s *OrderSuite) TestFetch(c *C) {
	account := mustCreateBankAccount(sharedClient, nil)
	defer deleteBankAccount(sharedClient, account, c)

	customer := mustCreateCustomer(sharedClient)
	// Can't delete a customer associated with an order
	// defer deleteCustomer(sharedClient, customer.Id, c)

	sharedClient.BankAccount.AssociateWithCustomer(account.Id, customer.Id)

	order, _, err := sharedClient.Order.Create(customer.Id, &Order{
		Description: "TestFetch",
	})
	c.Assert(err, IsNil)

	fetchedOrder, _, err := sharedClient.Order.Fetch(order.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedOrder.Id, Equals, order.Id)
	c.Assert(fetchedOrder.Description, Equals, order.Description)

	card := mustCreateCard(sharedClient)
	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 100,
		Order:  order.Href,
	})
	c.Assert(debit.Amount, Equals, 100)
	c.Assert(err, IsNil)

	fetchedOrder, _, err = sharedClient.Order.Fetch(order.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedOrder.Amount, Equals, 100)
	c.Assert(fetchedOrder.AmountEscrowed, Equals, 100)

	credit, _, err := sharedClient.Credit.CreateForOrder(account.Id, order.Id, &Credit{
		Amount: 25,
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 25)
	c.Assert(credit.Status, Equals, Succeeded)

	fetchedOrder, _, err = sharedClient.Order.Fetch(order.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedOrder.Amount, Equals, 100)
	c.Assert(fetchedOrder.AmountEscrowed, Equals, 75)
}

func (s *OrderSuite) TestList(c *C) {
	startingOrderPage, _, err := sharedClient.Order.List()
	c.Assert(err, IsNil)

	customer := mustCreateCustomer(sharedClient)
	// Can't delete a customer associated with an order
	// defer deleteCustomer(sharedClient, customer.Id, c)

	order, _, err := sharedClient.Order.Create(customer.Id, &Order{
		Description: "TestList",
	})
	c.Assert(err, IsNil)
	c.Assert(order, Not(IsNil))

	orderPage, _, err := sharedClient.Order.List()
	c.Assert(err, IsNil)
	c.Assert(orderPage.Total, Equals, 1+startingOrderPage.Total)
}

func (s *OrderSuite) TestUpdate(c *C) {
	customer := mustCreateCustomer(sharedClient)
	// Can't delete a customer associated with an order
	// defer deleteCustomer(sharedClient, customer.Id, c)

	order, _, err := sharedClient.Order.Create(customer.Id, &Order{
		Description: "TestUpdate",
	})
	c.Assert(err, IsNil)

	updatedOrder, _, err := sharedClient.Order.Update(order.Id, map[string]interface{}{
		"description": "sample description",
		"meta": map[string]interface{}{
			"xxx": "yyy",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(updatedOrder.Description, Equals, "sample description")
	c.Assert(updatedOrder.Meta["xxx"], Equals, "yyy")
}

type ReversalSuite struct{}

var _ = Suite(&ReversalSuite{})

func setupReversalPrereqs(c *C) (*Credit, *Card, *BankAccount) {
	card := mustCreateCardFixture(sharedClient, "VisaCreditable")
	account := mustCreateBankAccount(sharedClient, nil)

	// Debit funds from a card into Escrow, so we can have money to pay out
	debit, _, err := sharedClient.Card.Charge(card.Id, &Debit{
		Amount: 70,
	})
	c.Assert(err, IsNil)
	c.Assert(debit.Amount, Equals, 70)

	credit, _, err := sharedClient.BankAccount.Credit(account.Id, &Credit{
		Amount:               70,
		AppearsOnStatementAs: "Credit",
		Description:          "Credit for Test Reversal",
	})
	c.Assert(err, IsNil)
	c.Assert(credit.Amount, Equals, 70)

	return credit, card, account
}

func teardownReversalPrereqs(card *Card, account *BankAccount, c *C) {
	defer deleteCard(sharedClient, card, c)
	defer deleteBankAccount(sharedClient, account, c)
}

func (s *ReversalSuite) TestCreate(c *C) {
	credit, card, account := setupReversalPrereqs(c)
	defer teardownReversalPrereqs(card, account, c)

	reversal, _, err := sharedClient.Reversal.Create(credit.Id, &Reversal{
		Amount:      70,
		Description: "Test reversal",
		Meta: map[string]string{
			"merchant.feedback":          "positive",
			"user.refund_reason":         "not happy with product",
			"fulfillment.item_condition": "OK",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(reversal.Amount, Equals, 70)
	c.Assert(reversal.Meta["merchant.feedback"], Equals, "positive")
	c.Assert(reversal.Meta["user.refund_reason"], Equals, "not happy with product")
	c.Assert(reversal.Meta["fulfillment.item_condition"], Equals, "OK")
}

func (s *ReversalSuite) TestFetch(c *C) {
	credit, card, account := setupReversalPrereqs(c)
	defer teardownReversalPrereqs(card, account, c)

	reversal, _, err := sharedClient.Reversal.Create(credit.Id, &Reversal{
		Amount:      70,
		Description: "Test reversal",
		Meta: map[string]string{
			"merchant.feedback":          "positive",
			"user.refund_reason":         "not happy with product",
			"fulfillment.item_condition": "OK",
		},
	})
	c.Assert(err, IsNil)

	fetchedReversal, _, err := sharedClient.Reversal.Fetch(reversal.Id)
	c.Assert(err, IsNil)
	c.Assert(fetchedReversal.Id, Equals, reversal.Id)
	c.Assert(fetchedReversal.Amount, Equals, reversal.Amount)
	c.Assert(fetchedReversal.Description, Equals, reversal.Description)
	c.Assert(fetchedReversal.Meta["merchant.feedback"], Equals, reversal.Meta["merchant.feedback"])
	c.Assert(fetchedReversal.Meta["user.refund_reason"], Equals, reversal.Meta["user.refund_reason"])
	c.Assert(fetchedReversal.Meta["fulfillment.item_condition"], Equals, reversal.Meta["fulfillment.item_condition"])
}

func (s *ReversalSuite) TestList(c *C) {
	credit, card, account := setupReversalPrereqs(c)
	defer teardownReversalPrereqs(card, account, c)

	startingReversalPage, _, err := sharedClient.Reversal.List()
	c.Assert(err, IsNil)

	reversal, _, err := sharedClient.Reversal.Create(credit.Id, &Reversal{Amount: 70})
	c.Assert(err, IsNil)
	c.Assert(reversal, Not(IsNil))

	reversalPage, _, err := sharedClient.Reversal.List()
	c.Assert(err, IsNil)
	c.Assert(reversalPage.Total, Equals, 1+startingReversalPage.Total)
}

func (s *ReversalSuite) TestUpdate(c *C) {
	credit, card, account := setupReversalPrereqs(c)
	defer teardownReversalPrereqs(card, account, c)

	reversal, _, err := sharedClient.Reversal.Create(credit.Id, &Reversal{
		Amount:      70,
		Description: "Starting reversal description",
	})
	c.Assert(err, IsNil)

	updatedReversal, _, err := sharedClient.Reversal.Update(reversal.Id, map[string]interface{}{
		"description": "Ending reversal description",
	})
	c.Assert(err, IsNil)
	c.Assert(updatedReversal.Description, Equals, "Ending reversal description")
}
