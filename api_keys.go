package balanced

import (
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
}

func (s *ApiKeyService) Create() (*ApiKey, *http.Response, error) {
	apiKeyResponse := new(ApiKeyResponse)
	httpResponse, err := s.client.POST("/api_keys", nil, nil, apiKeyResponse)
	if err != nil {
		return nil, httpResponse, err
	}
	return &apiKeyResponse.ApiKeys[0], httpResponse, nil
}
