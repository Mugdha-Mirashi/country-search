package handler

import (
	"country-search/business"
	"country-search/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CountrySearchController struct {
	service *business.CountrySearchService
}

// NewCountrySearchController creates a new CountrySearchController instance.
func NewCountrySearchController(service *business.CountrySearchService) *CountrySearchController {
	return &CountrySearchController{
		service: service,
	}
}

// HandleCountrySearch godoc
// @Summary Search for a country
// @Description Get country information by name
// @Tags countries
// @Accept json
// @Produce json
// @Param name query string true "Country name"
// @Success 200 {array} models.CountrySearchResponseModel
// @Failure 400 {object} gin.H{"error": "Name is required"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /api/countries/search [get]
func (controller *CountrySearchController) HandleCountrySearch(ctx *gin.Context) {

	name := ctx.Query("name")

	if name == "" {
		log.Println("Country search failed: Name is required")
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Name is required",
		})
		return
	}

	countries, err := controller.service.CountrySearch(name)
	if err != nil {
		log.Println("Country search failed: Name is required")
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	log.Printf("Country search successful for: %s", name)
	ctx.JSON(http.StatusOK, countries)
}
