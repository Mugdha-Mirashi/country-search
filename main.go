package main

import (
	"country-search/business"
	"country-search/cache"
	"country-search/handler"
	"country-search/httpclient"
	"log"

	_ "country-search/docs" // Import the generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()

	// Initialize cache and HTTP client
	cacheInstance := cache.NewCache()
	httpClientInstance := httpclient.NewClient("https://restcountries.com/v3.1")

	// Create service and controller
	service := business.NewCountrySearchService(cacheInstance, httpClientInstance)
	controller := handler.NewCountrySearchController(service)

	// Register the handler
	router.GET("/api/countries/search", controller.HandleCountrySearch)

	// Add Swagger middleware
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The URL pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Start the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
