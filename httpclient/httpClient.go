// internal/httpclient/http_client.go
package httpclient

import (
	"context"
	"country-search/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client interface defines the method for fetching country data
type Client interface {
	FetchCountryData(name string) (*models.CountrySearchResponseModel, error)
}

// httpClient is the implementation of Client
type httpClient struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new instance of httpClient
func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchCountryData fetches country data from the REST Countries API
func (c *httpClient) FetchCountryData(name string) (*models.CountrySearchResponseModel, error) {
	url := fmt.Sprintf("%s/name/%s", c.baseURL, name)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch country data")
	}

	var countries []struct {
		Name       string `json:"name"`
		Capital    string `json:"capital"`
		Currencies []struct {
			Code string `json:"code"`
		} `json:"currencies"`
		Population int `json:"population"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, err
	}

	if len(countries) == 0 {
		return nil, errors.New("country not found")
	}

	country := &models.CountrySearchResponseModel{
		Name:       countries[0].Name,
		Capital:    countries[0].Capital,
		Currency:   countries[0].Currencies[0].Code,
		Population: countries[0].Population,
	}

	return country, nil
}
