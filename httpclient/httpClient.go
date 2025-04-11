// internal/httpclient/http_client.go
package httpclient

import (
	"context"
	"country-search/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

type Country struct {
	Name struct {
		Common     string `json:"common"`
		Official   string `json:"official"`
		NativeName map[string]struct {
			Official string `json:"official"`
			Common   string `json:"common"`
		} `json:"nativeName"`
	} `json:"name"`
	TLD         []string `json:"tld"`
	CCA2        string   `json:"cca2"`
	CCN3        string   `json:"ccn3"`
	CCA3        string   `json:"cca3"`
	CIOC        string   `json:"cioc"`
	Independent bool     `json:"independent"`
	Status      string   `json:"status"`
	UNMember    bool     `json:"unMember"`
	Currencies  map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	IDD struct {
		Root     string   `json:"root"`
		Suffixes []string `json:"suffixes"`
	} `json:"idd"`
	Capital      []string          `json:"capital"`
	AltSpellings []string          `json:"altSpellings"`
	Region       string            `json:"region"`
	Subregion    string            `json:"subregion"`
	Languages    map[string]string `json:"languages"`
	Translations map[string]struct {
		Official string `json:"official"`
		Common   string `json:"common"`
	} `json:"translations"`
	Latlng     []float64 `json:"latlng"`
	Landlocked bool      `json:"landlocked"`
	Borders    []string  `json:"borders"`
	Area       float64   `json:"area"`
	Demonyms   map[string]struct {
		F string `json:"f"`
		M string `json:"m"`
	} `json:"demonyms"`
	Flag string `json:"flag"`
	Maps struct {
		GoogleMaps     string `json:"googleMaps"`
		OpenStreetMaps string `json:"openStreetMaps"`
	} `json:"maps"`
	Population int                `json:"population"`
	Gini       map[string]float64 `json:"gini"`
	Fifa       string             `json:"fifa"`
	Car        struct {
		Signs []string `json:"signs"`
		Side  string   `json:"side"`
	} `json:"car"`
	Timezones  []string `json:"timezones"`
	Continents []string `json:"continents"`
	Flags      struct {
		Png string `json:"png"`
		Svg string `json:"svg"`
		Alt string `json:"alt"`
	} `json:"flags"`
	CoatOfArms struct {
		Png string `json:"png"`
		Svg string `json:"svg"`
	} `json:"coatOfArms"`
	StartOfWeek string `json:"startOfWeek"`
	CapitalInfo struct {
		Latlng []float64 `json:"latlng"`
	} `json:"capitalInfo"`
	PostalCode struct {
		Format string `json:"format"`
		Regex  string `json:"regex"`
	} `json:"postalCode"`
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response Body: %s\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch country data")
	}

	var countries []Country

	if err := json.Unmarshal(bodyBytes, &countries); err != nil {
		return nil, err
	}

	if len(countries) == 0 {
		return nil, errors.New("country not found")
	}

	var currencySymbol string
	for _, currency := range countries[0].Currencies {
		// currencyCode = code
		// currencyName = currency.Name
		currencySymbol = currency.Symbol
		break
	}

	country := &models.CountrySearchResponseModel{
		Name:       countries[0].Name.Common,
		Capital:    countries[0].Capital[0],
		Currency:   currencySymbol,
		Population: countries[0].Population,
	}

	return country, nil
}
