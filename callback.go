package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type CallbackService struct {
	client *Client
}

type Callback struct {
	Url string `json:"url"`

	// "post", "put", or "get"
	Method    string         `json:"method,omitempty"`
	Links     *CallbackLinks `json:"links,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	Id        string         `json:"id,omitempty"`
	Href      string         `json:"href,omitempty"`
	Revision  string         `json:"revision,omitempty"`
}

type CallbackLinks struct{}

type CallbackPage struct {
	Callbacks []Callback
	*PaginationParams
}

type CallbackResponse struct {
	Callbacks []Callback             `json:"callbacks"`
	Meta      map[string]interface{} `json:"meta"`
	Links     *CallbackResponseLinks `json:"links"`
}

type CallbackResponseLinks struct{}

func (s *CallbackService) Create(url, method string) (*Callback, *http.Response, error) {
	callbackResponse := new(CallbackResponse)
	callback := &Callback{
		Url:    url,
		Method: method,
	}
	httpResponse, err := s.client.POST("/callbacks", nil, callback, callbackResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &callbackResponse.Callbacks[0], httpResponse, nil
}

func (s *CallbackService) Fetch(callbackId string) (*Callback, *http.Response, error) {
	path := fmt.Sprintf("/callbacks/%v", callbackId)
	callbackResponse := new(CallbackResponse)
	httpResponse, err := s.client.GET(path, nil, nil, callbackResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &callbackResponse.Callbacks[0], httpResponse, nil
}

func (s *CallbackService) Delete(callbackId string) (bool, *http.Response, error) {
	path := fmt.Sprintf("/callbacks/%v", callbackId)
	httpResponse, err := s.client.DELETE(path, nil, nil, nil)
	if err != nil {
		return false, httpResponse, err
	}
	code := httpResponse.StatusCode
	didDelete := 200 <= code && code < 300
	return didDelete, httpResponse, nil
}

func (s *CallbackService) List(args ...interface{}) (*CallbackPage, *http.Response, error) {
	// Turns args into a map[string]int with "offset" and "limit" keys
	query := paginatedArgsToQuery(args)
	callbackResponse := new(CallbackResponse)
	httpResponse, err := s.client.GET("/callbacks", query, nil, callbackResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &CallbackPage{
		Callbacks:        callbackResponse.Callbacks,
		PaginationParams: NewPaginationParams(callbackResponse.Meta),
	}, httpResponse, nil
}
