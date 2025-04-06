package models

type CountrySearchResponseModel struct {
	Name       string `json:"name" example:"India"`
	Capital    string `json:"capital" example:"New Delhi"`
	Currency   string `json:"currency" example:"₹"`
	Population int    `json:"population" example:"1380000000"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
