package balanced

import (
	"fmt"
	"net/http"
	"time"
)

type ApiKeyService struct {
	client *Client
}

type ApiKey struct {
	Id        string            `json:"id"`
	Href      string            `json:"href"`
	Links     map[string]string `json:"links"`
	Meta      map[string]string `json:"meta"`
	Secret    string            `json:"secret"`
	CreatedAt *time.Time        `json:"created_at"`
}

type ApiKeyResponse struct {
	ApiKeys []ApiKey          `json:"api_keys"`
	Links   map[string]string `json:"links"`
	Meta    map[string]interface{}
}

type ApiKeyPage struct {
	ApiKeys []ApiKey
	*PaginationParams
}

func (s *ApiKeyService) Create() (*ApiKey, *http.Response, error) {
	apiKeyResponse := new(ApiKeyResponse)
	httpResponse, err := s.client.POST("/api_keys", nil, nil, apiKeyResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &apiKeyResponse.ApiKeys[0], httpResponse, nil
}

func (s *ApiKeyService) Fetch(id string) (*ApiKey, *http.Response, error) {
	path := fmt.Sprintf("/api_keys/%v", id)
	apiKeyResponse := new(ApiKeyResponse)
	httpResponse, err := s.client.GET(path, nil, nil, apiKeyResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &apiKeyResponse.ApiKeys[0], httpResponse, nil
}

func (s *ApiKeyService) List(args ...interface{}) (*ApiKeyPage, *http.Response, error) {
	query := paginatedArgsToQuery(args)
	apiKeyResponse := new(ApiKeyResponse)
	httpResponse, err := s.client.GET("/api_keys", query, nil, apiKeyResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &ApiKeyPage{
		ApiKeys:          apiKeyResponse.ApiKeys,
		PaginationParams: NewPaginationParams(apiKeyResponse.Meta),
	}, httpResponse, nil
}

func (s *ApiKeyService) Delete(id string) (bool, *http.Response, error) {
	path := fmt.Sprintf("/api_keys/%v", id)
	httpResponse, err := s.client.DELETE(path, nil, nil, nil)
	if err != nil {
		return false, httpResponse, err
	}
	code := httpResponse.StatusCode
	didDelete := 200 <= code && code < 300
	return didDelete, httpResponse, nil
}
