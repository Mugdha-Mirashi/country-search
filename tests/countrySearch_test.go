package tests

import (
	"country-search/business"
	"country-search/handler"
	"country-search/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCache is a mock implementation of the Cache interface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCache) Set(key string, value interface{}) {
	m.Called(key, value)
}

// MockHTTPClient is a mock implementation of the HTTP client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) FetchCountryData(name string) (*models.CountrySearchResponseModel, error) {
	args := m.Called(name)
	return args.Get(0).(*models.CountrySearchResponseModel), args.Error(1)
}

func TestGetCountrySearch200Success(t *testing.T) {
	mockCache := new(MockCache)
	mockHTTPClient := new(MockHTTPClient)
	service := business.NewCountrySearchService(mockCache, mockHTTPClient)

	expectedCountry := &models.CountrySearchResponseModel{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "INR",
		Population: 1380004385,
	}

	mockCache.On("Get", "India").Return(expectedCountry, true)

	country, err := service.CountrySearch("India")

	assert.NoError(t, err)
	assert.Equal(t, expectedCountry, country)
	mockCache.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestGetCountryInfo_CacheMiss(t *testing.T) {
	mockCache := new(MockCache)
	mockHTTPClient := new(MockHTTPClient)
	service := business.NewCountrySearchService(mockCache, mockHTTPClient)

	expectedCountry := &models.CountrySearchResponseModel{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "INR",
		Population: 1380004385,
	}

	mockCache.On("Get", "India").Return(nil, false)
	mockHTTPClient.On("FetchCountryData", "India").Return(expectedCountry, nil)
	mockCache.On("Set", "India", expectedCountry).Return()

	country, err := service.CountrySearch("India")

	assert.NoError(t, err)
	assert.Equal(t, expectedCountry, country)
	mockCache.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestGetCountryInfo_ApiError(t *testing.T) {
	mockCache := new(MockCache)
	mockHTTPClient := new(MockHTTPClient)
	service := business.NewCountrySearchService(mockCache, mockHTTPClient)

	mockCache.On("Get", "India").Return(nil, false)
	mockHTTPClient.On("FetchCountryData", "India").Return((*models.CountrySearchResponseModel)(nil), errors.New("API error"))

	country, err := service.CountrySearch("India")

	assert.Error(t, err)
	assert.Nil(t, country)
	mockCache.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestHandleCountrySearch_MissingName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.Default()

	mockCache := new(MockCache)
	mockHTTPClient := new(MockHTTPClient)
	service := business.NewCountrySearchService(mockCache, mockHTTPClient)
	controller := handler.NewCountrySearchController(service)

	// Register the handler
	router.GET("/api/countries/search", controller.HandleCountrySearch)

	// Create a new HTTP request without the 'name' query parameter
	req, _ := http.NewRequest(http.MethodGet, "/api/countries/search", nil)

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "Name is required"}`, w.Body.String())
}
