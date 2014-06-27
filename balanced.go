package balanced

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	Version    = "0.1"
	ApiVersion = "v1.1"
	baseUrlStr = "https://api.balancedpayments.com/"

	Pending   = "pending"
	Succeeded = "succeeded"
	Failed    = "failed"

	// Status values used for disputes
	Won  = "won"
	Lost = "lost"
)

var baseURL *url.URL = mustParseUrl(baseUrlStr)

func mustParseUrl(urlStr string) *url.URL {
	url, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return url
}

type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client
	secret string

	ApiKey       *ApiKeyService
	BankAccount  *BankAccountService
	Verification *VerificationService
	Callback     *CallbackService
	Card         *CardService
	CardHold     *CardHoldService
	Credit       *CreditService
	Customer     *CustomerService
	Debit        *DebitService
	Event        *EventService
	Order        *OrderService
	Refund       *RefundService
	Reversal     *ReversalService
	Marketplace  *MarketplaceService
}

func NewClient(httpClient *http.Client, secret string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, secret: secret}
	c.ApiKey = &ApiKeyService{client: c}
	c.BankAccount = &BankAccountService{client: c}
	c.Verification = &VerificationService{client: c}
	c.Callback = &CallbackService{client: c}
	c.Card = &CardService{client: c}
	c.CardHold = &CardHoldService{client: c}
	c.Credit = &CreditService{client: c}
	c.Customer = &CustomerService{client: c}
	c.Debit = &DebitService{client: c}
	c.Event = &EventService{client: c}
	c.Order = &OrderService{client: c}
	c.Refund = &RefundService{client: c}
	c.Reversal = &ReversalService{client: c}
	c.Marketplace = &MarketplaceService{client: c}

	return c
}

func (c *Client) GET(urlPath string, query map[string]interface{}, reqBody interface{}, resObj interface{}) (*http.Response, error) {
	return c.buildAndDoRequest("GET", urlPath, query, reqBody, resObj)
}

func (c *Client) POST(urlPath string, query map[string]interface{}, reqBody interface{}, resObj interface{}) (*http.Response, error) {
	return c.buildAndDoRequest("POST", urlPath, query, reqBody, resObj)
}

func (c *Client) PUT(urlPath string, query map[string]interface{}, reqBody interface{}, resObj interface{}) (*http.Response, error) {
	return c.buildAndDoRequest("PUT", urlPath, query, reqBody, resObj)
}

func (c *Client) DELETE(urlPath string, query map[string]interface{}, reqBody interface{}, resObj interface{}) (*http.Response, error) {
	return c.buildAndDoRequest("DELETE", urlPath, query, reqBody, resObj)
}

func (c *Client) buildAndDoRequest(method, urlPath string, query map[string]interface{}, reqBody interface{}, resObj interface{}) (*http.Response, error) {
	req, err := c.NewRequest(method, urlPath, query, reqBody)
	if err != nil {
		return nil, err
	}
	return c.Do(req, resObj)
}

func mapToQueryVals(params map[string]interface{}) url.Values {
	values := make(url.Values)
	for k, v := range params {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return values
}

func (c *Client) NewRequest(method, urlPath string, queryParams map[string]interface{}, body interface{}) (*http.Request, error) {
	url, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	u := baseURL.ResolveReference(url)

	if queryParams != nil {
		qs := mapToQueryVals(queryParams)
		if err != nil {
			return nil, err
		}
		u.RawQuery = qs.Encode()
	}

	buff := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buff).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buff)
	if err != nil {
		return nil, err
	}

	// req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.api+json;revision=1.1")

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("User-Agent", "balanced-go/1.1")

	if c.secret != "" {
		req.SetBasicAuth(c.secret, "")
	}

	return req, nil
}

// Do sends an API request and returns an API response
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	err = checkResponse(res)
	if err != nil {
		return res, err
	}

	if v != nil {
		// If v implements the io.Writer interface, the raw response body will be
		// written to v, without attempting to decode it first.
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, res.Body)
		} else {
			err = json.NewDecoder(res.Body).Decode(v)
		}
	}

	return res, nil
}

type ErrorResponse struct {
	// http response that caused this error
	*http.Response

	Errors []ErrorResponseError `json:"errors"`
}

type ErrorResponseError struct {
	Status       string `json:"status"`
	CategoryCode string `json:"category_code"`
	CategoryType string `json:"category_type"`
	Description  string `json:"description"`
	RequestId    string `json:"request_id"`
	StatusCode   int    `json:"status_code"`
}

func (r *ErrorResponse) Error() string {
	errors := r.Errors
	var errDescr string
	if len(errors) > 0 {
		errDescr = errors[0].Description
	} else {
		errDescr = ""
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, errDescr)
}

func checkResponse(res *http.Response) error {
	if c := res.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: res}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if data != nil {
		err = json.Unmarshal(data, errorResponse)
		// var dat map[string]interface{}
		// err = json.Unmarshal(data, &dat)
		if err != nil {
			return err
		}
	}
	return errorResponse
}
