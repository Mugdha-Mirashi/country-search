package business

import (
	"country-search/cache"
	"country-search/httpclient"
	"country-search/models"
	"log"
)

type CountrySearchService struct {
	cache      cache.Cache
	httpClient httpclient.Client
}

func NewCountrySearchService(cache cache.Cache, httpClient httpclient.Client) *CountrySearchService {
	return &CountrySearchService{
		cache:      cache,
		httpClient: httpClient,
	}
}

func (service *CountrySearchService) CountrySearch(name string) (*models.CountrySearchResponseModel, error) {
	if cachedData, found := service.cache.Get(name); found {
		if country, ok := cachedData.(*models.CountrySearchResponseModel); ok {
			return country, nil
		}
	}

	countryData, err := service.httpClient.FetchCountryData(name)
	if err != nil {
		log.Printf("Error fetching country data from API for %s: %v", name, err)
		return nil, err
	}

	service.cache.Set(name, countryData)

	log.Printf("Country data retrieval successful for: %s", name)
	return countryData, nil
}
